package gocronometer

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"
)

const (
	// HTMLLoginURL is the full URL to the Cronometer login page.
	HTMLLoginURL = "https://cronometer.com/login/"

	// APILoginURL is the full URL for login requests.
	APILoginURL = "https://cronometer.com/login"

	// GWTBaseURL is the full URL for accessing the GWT API.
	GWTBaseURL = "https://cronometer.com/cronometer/app"

	// APIExportURL is the full URL for requesting data exports.
	APIExportURL = "https://cronometer.com/export"
)

var GWTTokenRegex = regexp.MustCompile("\"(?P<token>.*)\"")

const GWTAuthRegex = `OK\[(?P<userid>\d*),.*`

var GWTAuthenticationRegexp = regexp.MustCompile(GWTAuthRegex)

// Client represents a client to the Cronometer API. The zero value is not a valid configuration. A new client should
// be generated with the NewClient function.
type Client struct {
	HTTPClient *http.Client
	Nonce      string
	UserID     string

	GWTContentType string
	GWTModuleBase  string
	GWTPermutation string
	GWTHeader      string
}

// ClientOptions represents the options that can be provided to the client. Zero values revert to the library defaults.
type ClientOptions struct {
	GWTContentType string
	GWTModuleBase  string
	GWTPermutation string
	GWTHeader      string
}

// updateOpts updates the client with the opts provided
func (c *Client) updateOpts(opts *ClientOptions) {
	// A nil opt is the same as a zero value opt.
	if opts == nil {
		opts = &ClientOptions{}
	}

	if opts.GWTContentType != "" {
		c.GWTContentType = opts.GWTContentType
	}
	if opts.GWTModuleBase != "" {
		c.GWTModuleBase = opts.GWTModuleBase
	}
	if opts.GWTPermutation != "" {
		c.GWTPermutation = opts.GWTPermutation
	}
	if opts.GWTHeader != "" {
		c.GWTHeader = opts.GWTHeader
	}
}

// NewClient generates a new client for the Cronometer API. If opts is nil the default values are utilized.
func NewClient(opts *ClientOptions) *Client {
	jar, _ := cookiejar.New(nil)
	client := &Client{
		HTTPClient: &http.Client{
			Jar: jar,
		},
		GWTContentType: GWTContentType,
		GWTModuleBase:  GWTModuleBase,
		GWTPermutation: GWTPermutation,
	}

	client.updateOpts(opts)

	return client
}

// NewGWTRequestWithContext creates a new http request with the proper headers for a GWT request.
func (c *Client) NewGWTRequestWithContext(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("content-type", c.GWTContentType)
	req.Header.Add("x-gwt-module-base", c.GWTModuleBase)
	req.Header.Add("x-gwt-permutation", c.GWTPermutation)

	return req, nil
}

// ObtainAntiCSRF connects to the login page of Cronometer and parses out the anticsrf value from the HTML form.
func (c *Client) ObtainAntiCSRF(ctx context.Context) (string, error) {

	// Building and executing request to obtain the login page HTML.
	req, err := http.NewRequestWithContext(ctx, "GET", HTMLLoginURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build request to retreive anticsrf value: %s", err)
	}
	req = req.WithContext(ctx)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed issuing HTTP request: %s", err)
	}
	defer closeAndExhaustReader(resp.Body)

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

type LoginResponse struct {
	Redirect string `json:"redirect"`
	Success  bool   `json:"success"`
	Error    string `json:"error"`
}

// Login logs into the Cronometer and the GWT API. Nil is returned on login success.
func (c *Client) Login(ctx context.Context, username string, password string) error {
	// Obtaining a new anticsrf from the login page.
	antiCSRF, err := c.ObtainAntiCSRF(ctx)
	if err != nil {
		return fmt.Errorf("failed to retrieve anit CSRF: %s", err)
	}

	// Building login request.
	formData := url.Values{}
	formData.Set("anticsrf", antiCSRF)
	formData.Set("password", password)
	formData.Set("username", username)

	req, err := http.NewRequestWithContext(ctx, "POST", APILoginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed while building http request for login: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed while executing http request for login: %s", err)
	}
	defer closeAndExhaustReader(resp.Body)

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

	// Storing the nonce from provided cookies.
	c.updateSesnonce(resp)

	// Authenticating with GWT.
	err = c.GWTAuthenticate(ctx)
	if err != nil {
		return fmt.Errorf("failed to authenticate with GWT: %s", err)
	}

	return nil
}

func (c *Client) updateSesnonce(resp *http.Response) {
	if resp == nil {
		return
	}

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "sesnonce" {
			c.Nonce = cookie.Value
		}
	}
}

// Logout logs out from the API.
func (c *Client) Logout(ctx context.Context) error {
	// Building the request.
	reqBody := fmt.Sprintf(GWTLogout, c.Nonce)

	req, err := c.NewGWTRequestWithContext(ctx, "POST", GWTBaseURL, strings.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed while building http request for gwt authentication: %s", err)
	}

	// Executing the request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed while executing http request for gwt logout: %s", err)
	}
	//noinspection GoUnhandledErrorResult
	defer closeAndExhaustReader(resp.Body)

	// Handling the response.
	if resp.StatusCode != 200 {
		return fmt.Errorf("received non 200 response of %d for gwt logout", resp.StatusCode)
	}

	c.UserID = ""
	c.Nonce = ""

	return nil
}

// GWTAuthenticate will authenticate with the GWT API using the sesnonce of the client. Login() calls this by default so
// in most cases this should never be called directly.
func (c *Client) GWTAuthenticate(ctx context.Context) error {
	// Building and sending the request.
	//reqBody := fmt.Sprintf(GWTAuthenticate, c.Nonce)

	req, err := c.NewGWTRequestWithContext(ctx, "POST", GWTBaseURL, strings.NewReader(GWTAuthenticate))
	if err != nil {
		return fmt.Errorf("failed while building http request for gwt authentication: %s", err)
	}

	// Executing the request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed while executing http request for gwt authentication: %s", err)
	}
	defer closeAndExhaustReader(resp.Body)

	// Handling the response.
	if resp.StatusCode != 200 {
		return fmt.Errorf("received non 200 response of %d for gwt token generation", resp.StatusCode)
	}

	c.updateSesnonce(resp)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body of gwt token authentication: %s", err)
	}

	match := GWTAuthenticationRegexp.FindStringSubmatch(string(body))

	if len(match) != 2 {
		return fmt.Errorf("failed to find GWT Authentication token in response data, expected 2 matches but received %d", len(match))
	}

	c.UserID = match[1]

	return nil
}

// GenerateAuthToken requests an authentication token from the API. This token is used to request the generation of
// a "token" that is provided as a nonce to the export API calls.
func (c *Client) GenerateAuthToken(ctx context.Context) (string, error) {

	// Building the request.
	reqBody := fmt.Sprintf(GWTGenerateAuthToken, c.Nonce, c.UserID)

	req, err := c.NewGWTRequestWithContext(ctx, "POST", GWTBaseURL, strings.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("failed while building http request for gwt token generation: %s", err)
	}

	// Executing the request.
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed while executing http request for gwt token generation: %s", err)
	}
	//noinspection GoUnhandledErrorResult
	defer closeAndExhaustReader(resp.Body)

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
	req, err := c.NewExportRequest(ctx, "GET", APIExportURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed while building http request for daily nutrition export: %s", err)
	}

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
	defer closeAndExhaustReader(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of daily nutrition export response: %s", err)
	}

	// Handling the response.
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non 200 response of %d for daily nutrition export: body %s", resp.StatusCode, string(body))
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
	req, err := c.NewExportRequest(ctx, "GET", APIExportURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed while building http request for servings export: %s", err)
	}

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
	defer closeAndExhaustReader(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of servings export response: %s", err)
	}

	// Handling the response.
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non 200 response of %d for servings export: body [%s]", resp.StatusCode, string(body))
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
	req, err := c.NewExportRequest(ctx, "GET", APIExportURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed while building http request for exercises export: %s", err)
	}

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
	defer closeAndExhaustReader(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body of exercises export response: %s", err)
	}

	// Handling the response.
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("received non 200 response of %d for exercises export: body %s", resp.StatusCode, string(body))
	}

	return string(body), nil
}

// NewExportRequest creates a new http request for exports.
func (c *Client) NewExportRequest(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Add("sec-fetch-mode", "navigate")
	req.Header.Add("sec-fetch-site", "same-origin")

	return req, nil
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
	req, err := c.NewExportRequest(ctx, "GET", APIExportURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed while building http request for biometrics export: %s", err)
	}

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
	defer closeAndExhaustReader(resp.Body)

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
	req, err := c.NewExportRequest(ctx, "GET", APIExportURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed while building http request for notes export: %s", err)
	}

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
	defer closeAndExhaustReader(resp.Body)

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

// ExportServingsParsed exports the servings within the date range and parses them into a go struct. Only the YYYY-mm-dd is utilized of startDate and
// endDate. The export is the raw string data.
func (c *Client) ExportServingsParsed(ctx context.Context, startDate time.Time, endDate time.Time) (ServingRecords, error) {
	raw, err := c.ExportServings(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("retreiving raw data: %s", err)
	}

	servings, err := ParseServingsExport(strings.NewReader(raw), time.UTC)
	if err != nil {
		return nil, fmt.Errorf("parsing raw data: %s", err)
	}

	return servings, nil
}

// ExportServingsParsedWithLocation is the same as ExportServingsParsed but sets the location of every recorded time
// to the location provided.
func (c *Client) ExportServingsParsedWithLocation(ctx context.Context, startDate time.Time, endDate time.Time, location *time.Location) (ServingRecords, error) {
	raw, err := c.ExportServings(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("retreiving raw data: %s", err)
	}

	servings, err := ParseServingsExport(strings.NewReader(raw), location)
	if err != nil {
		return nil, fmt.Errorf("parsing raw data: %s", err)
	}

	return servings, nil
}

// ExportExercisesParsedWithLocation exports the exercises within the date range and parses them into a go struct. Only the YYYY-mm-dd is utilized of startDate and
// endDate. The export is parsed and dates set to the location provided.
func (c *Client) ExportExercisesParsedWithLocation(ctx context.Context, startDate time.Time, endDate time.Time, location *time.Location) (ExerciseRecords, error) {
	raw, err := c.ExportExercises(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("retreiving raw data: %s", err)
	}

	exercises, err := ParseExerciseExport(strings.NewReader(raw), location)
	if err != nil {
		return nil, fmt.Errorf("parsing raw data: %s", err)
	}

	return exercises, nil
}

// closeAndExhaustReader will first try and exhaust r and then call close. Errors are intentionally ignored
// as this is only to be called in with deferred and where the error would have no action to be taken.
func closeAndExhaustReader(r io.ReadCloser) {
	if _, err := io.Copy(io.Discard, r); err != nil {
		// Do nothing.
	}
	if err := r.Close(); err != nil {
		// Do nothing.
	}
	return
}

// ExportBiometricRecordsParsedWithLocation exports the biometric records within the date range and parses them into a go struct. Only the YYYY-mm-dd is utilized of startDate and
// endDate. The export is parsed and dates set to the location provided.
func (c *Client) ExportBiometricRecordsParsedWithLocation(ctx context.Context, startDate time.Time, endDate time.Time, location *time.Location) (BiometricRecords, error) {
	raw, err := c.ExportBiometrics(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("retreiving raw data: %s", err)
	}

	exercises, err := ParseBiometricRecordsExport(strings.NewReader(raw), location)
	if err != nil {
		return nil, fmt.Errorf("parsing raw data: %s", err)
	}

	return exercises, nil
}
