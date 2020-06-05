package gocronometer

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	// HTMLPathLogin is the full path for the login page.
	HTMLPathLogin = "https://cronometer.com/login/"

	// APIPathLogin is the full path to the login API resource.
	APIPathLogin = "https://cronometer.com/login"

	// APIGWTPath is the full path for GWT RPC calls.
	APIGWTPath = "https://cronometer.com/cronometer/app"

	// APIPathExport is the full path for requesting exports.
	APIPathExport = "https://cronometer.com/export"
)

// These constants are considered "magic" values that allow GWT requests to process. The minimum amount of effort was
// put forth to get the GWT calls to work and it's expected some of these values will change with time.
const (
	GWTContentType = "text/x-gwt-rpc; charset=UTF-8"
	GWTModuleBase  = "https://cronometer.com/cronometer/"
	GWTPermutation = "9D62616AE775E1F90E83CD7804DC7AFE"
)

// These constants are the RPC body for the various GWT calls. It's expected that these may change out form under us.
const (
	// GWTRPCGenerateAuthorizationToken generates an authorization token. The only use known at this point is to provide
	// it as a nonce to the export calls.
	// Need to provide fmt.Sprintf() a nonce string to insert into the call. The nonce appears to really be used as a
	// session for GWT calls...
	GWTRPCGenerateAuthorizationToken = "7|0|8|https://cronometer.com/cronometer/|5BCB62A9B6F57CF6161F9EE3C6B77CD2|com.cronometer.client.CronometerService|generateAuthorizationToken|java.lang.String/2004016611|I|com.cronometer.client.data.AuthScope/3692935123|%s|1|2|3|4|4|5|6|6|7|8|70646|3600|7|2|"
)

var GWTTokenRegex = regexp.MustCompile("\"(?P<token>.*)\"")

type Client struct {
	HTTPClient *http.Client
	Nonce      string
}

func NewClient() *Client {
	cookieJar, _ := cookiejar.New(nil)
	return &Client{
		HTTPClient: &http.Client{
			Jar: cookieJar,
		},
	}
}

// Login retrieves a new Anit CSRF value and logs into the API. Upon login success error is nil.
func (c *Client) Login(ctx context.Context, username string, password string) error {
	// Retrieving AntiCSRF from the login form.
	antiCSRF, err := c.RetrieveAntiCSRF(ctx)
	if err != nil {
		return fmt.Errorf("failed while retreiving Anti CSRF for login: %s", err)
	}

	// Building url encoded values for login request.
	formData := url.Values{}
	formData.Set("anticsrf", antiCSRF)
	formData.Set("password", password)
	formData.Set("username", username)

	// Building the request.
	req, err := http.NewRequestWithContext(ctx, "POST", APIPathLogin, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed while building http request for login: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Executing the request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed while executing http request for login: %s", err)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Handling the response.
	if resp.StatusCode != 200 {
		return fmt.Errorf("received non 200 response of %d for login", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body of login response: %s", err)
	}

	var loginResponse LoginResponse
	if err = json.Unmarshal(body, &loginResponse); err != nil {
		return fmt.Errorf("failed to unmarshal login response json: %s", err)
	}

	if loginResponse.Error != "" {
		return fmt.Errorf("failed to login: %s", loginResponse.Error)
	}

	// Storing the nonse from provided cookies.
	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "sesnonce" {
			c.Nonce = cookie.Value
		}
	}

	return nil
}

// RetrieveAntiCSRF retrieves an Anti CSRF value needed for login. The only method to obtain this value is via the login
// form at https://cronometer.com/login/. An HTTP request is performed and the HTML page parsed to obtain the value.
func (c *Client) RetrieveAntiCSRF(ctx context.Context) (string, error) {

	// Building the request.
	req, err := http.NewRequestWithContext(ctx, "GET", HTMLPathLogin, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build request to retreive anticsrf value: %s", err)
	}
	req = req.WithContext(ctx)

	// Executing request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed issuing HTTP request: %s", err)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Handling the response.
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("recevied non OK error code %d", resp.StatusCode)
	}

	z, err := html.Parse(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to parse HTML response: %s", err)
	}

	var csrf string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, a := range n.Attr {
				if a.Key == "name" {
					if a.Val == "anticsrf" {
						for _, c := range n.Attr {
							if c.Key == "value" {
								csrf = c.Val
								break
							}
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(z)

	if csrf == "" {
		return "", fmt.Errorf("failed to find csrf value in HTML document")
	}

	return csrf, nil
}

// GenerateAuthToken requests an authentication token from the API. This token is used to request the generation of
// a "token" that is provided as a nonce to the export API calls.
func (c *Client) GenerateAuthToken(ctx context.Context) (string, error) {

	// Building the request.
	reqBody := fmt.Sprintf(GWTRPCGenerateAuthorizationToken, c.Nonce)

	req, err := http.NewRequestWithContext(ctx, "POST", APIGWTPath, strings.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed while building http request for gwt token generation: %s", err)
	}
	req.Header.Set("content-type", GWTContentType)
	req.Header.Add("x-gwt-module-base", GWTModuleBase)
	req.Header.Add("x-gwt-permutation", GWTPermutation)

	// Executing the request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed while executing http request for gwt token generation: %s", err)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Handling the response.
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non 200 response of %d for gwt token generation", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of gwt token generation response: %s", err)
	}

	match := GWTTokenRegex.FindStringSubmatch(string(body))

	if len(match) != 2 {
		return "", fmt.Errorf("failed to find token in response data, expected 2 matches but received %d", len(match))
	}

	return match[1], nil
}

// ExportDailyNutrition exports the daily nutrition values within the date range. Only the YYYY-mm-dd is utilized of startDate and
// endDate. The export is the raw string data.
func (c *Client) ExportDailyNutrition(ctx context.Context, startDate time.Time, endDate time.Time) (string, error) {
	// Generating the required token.
	token, err := c.GenerateAuthToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token to make request: %s", err)
	}

	// Building the request.
	req, err := http.NewRequestWithContext(ctx, "GET", APIPathExport, nil)
	if err != nil {
		return "", fmt.Errorf("failed while building http request for daily nutrition export: %s", err)
	}
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")

	q := req.URL.Query()
	q.Add("nonce", token)
	q.Add("generate", "dailySummary")
	q.Add("start", startDate.Format("2006-01-02"))
	q.Add("end", endDate.Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	// Executing the request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed while executing http request for daily nutrition export: %s", err)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Handling the response.
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non 200 response of %d for daily nutrition export", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of daily nutrition export response: %s", err)
	}

	return string(body), nil
}

// ExportServings exports all the services within the date range. Only the YYYY-mm-dd is utilized of startDate and
// endDate. The export is the raw string data.
func (c *Client) ExportServings(ctx context.Context, startDate time.Time, endDate time.Time) (string, error) {

	// Generating the required token.
	token, err := c.GenerateAuthToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token to make request: %s", err)
	}

	// Building the request.
	req, err := http.NewRequestWithContext(ctx, "GET", APIPathExport, nil)
	if err != nil {
		return "", fmt.Errorf("failed while building http request for servings export: %s", err)
	}
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")

	q := req.URL.Query()
	q.Add("nonce", token)
	q.Add("generate", "servings")
	q.Add("start", startDate.Format("2006-01-02"))
	q.Add("end", endDate.Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	// Executing the request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed while executing http request for servings export: %s", err)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Handling the response.
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non 200 response of %d for servings export", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of servings export response: %s", err)
	}

	return string(body), nil

}

// ExportExercises exports the exercises within the date range. Only the YYYY-mm-dd is utilized of startDate and
// endDate. The export is the raw string data.
func (c *Client) ExportExercises(ctx context.Context, startDate time.Time, endDate time.Time) (string, error) {
	// Generating the required token.
	token, err := c.GenerateAuthToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token to make request: %s", err)
	}

	// Building the request.
	req, err := http.NewRequestWithContext(ctx, "GET", APIPathExport, nil)
	if err != nil {
		return "", fmt.Errorf("failed while building http request for exercises export: %s", err)
	}
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")

	q := req.URL.Query()
	q.Add("nonce", token)
	q.Add("generate", "exercises")
	q.Add("start", startDate.Format("2006-01-02"))
	q.Add("end", endDate.Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	// Executing the request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed while executing http request for exercises export: %s", err)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Handling the response.
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non 200 response of %d for exercises export", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of exercises export response: %s", err)
	}

	return string(body), nil
}

// ExportBiometrics exports the biometrics within the date range. Only the YYYY-mm-dd is utilized of startDate and
// endDate. The export is the raw string data.
func (c *Client) ExportBiometrics(ctx context.Context, startDate time.Time, endDate time.Time) (string, error) {
	// Generating the required token.
	token, err := c.GenerateAuthToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token to make request: %s", err)
	}

	// Building the request.
	req, err := http.NewRequestWithContext(ctx, "GET", APIPathExport, nil)
	if err != nil {
		return "", fmt.Errorf("failed while building http request for biometrics export: %s", err)
	}
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")

	q := req.URL.Query()
	q.Add("nonce", token)
	q.Add("generate", "biometrics")
	q.Add("start", startDate.Format("2006-01-02"))
	q.Add("end", endDate.Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	// Executing the request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed while executing http request for biometrics export: %s", err)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Handling the response.
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non 200 response of %d for biometrics export", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of biometrics export response: %s", err)
	}

	return string(body), nil
}

// ExportNotes exports the notes within the date range. Only the YYYY-mm-dd is utilized of startDate and
// endDate. The export is the raw string data.
func (c *Client) ExportNotes(ctx context.Context, startDate time.Time, endDate time.Time) (string, error) {
	// Generating the required token.
	token, err := c.GenerateAuthToken(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get token to make request: %s", err)
	}

	// Building the request.
	req, err := http.NewRequestWithContext(ctx, "GET", APIPathExport, nil)
	if err != nil {
		return "", fmt.Errorf("failed while building http request for notes export: %s", err)
	}
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")

	q := req.URL.Query()
	q.Add("nonce", token)
	q.Add("generate", "notes")
	q.Add("start", startDate.Format("2006-01-02"))
	q.Add("end", endDate.Format("2006-01-02"))
	req.URL.RawQuery = q.Encode()

	// Executing the request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed while executing http request for notes export: %s", err)
	}
	//noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Handling the response.
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non 200 response of %d for notes export", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of notes export response: %s", err)
	}

	return string(body), nil
}
