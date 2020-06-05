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
	"strings"
)

const (
	HTMLPathLogin = "https://cronometer.com/login/"
	APIPathLogin  = "https://cronometer.com/login"
)

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

	// Executing the request.
	req, err := http.NewRequest("POST", APIPathLogin, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed while building http request for login: %s", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed while executing http request for login: %s", err)
	}
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
	req, err := http.NewRequest("GET", HTMLPathLogin, nil)
	if err != nil {
		return "", fmt.Errorf("failed to build request to retreive anticsrf value: %s", err)
	}
	req = req.WithContext(ctx)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed issuing HTTP request: %s", err)
	}
	defer resp.Body.Close()

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
// a "token" that is provided a a nonce to the export API calls.
func (c *Client) GenerateAuthToken(ctx context.Context) (string, error) {

}

func (c *Client) ExportDailyNutrition() {

}

func (c *Client) ExportServings() {

}

func (c *Client) ExportExercises() {

}

func (c *Client) ExportBiometrics() {

}

func (c *Client) ExportNotes() {

}
