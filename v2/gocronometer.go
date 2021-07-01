package v2

import (
	"context"
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"net/http/cookiejar"
)

const (
	// HTMLLoginURL is the full URL to the Cronometer login page.
	HTMLLoginURL = "https://cronometer.com/login/"

	// APILoginURL is the full URL for login requests.
	APILoginURL = "https://cronometer.com/login"

	// GWTBaseURL is the full URL for accessing the GWT API.
	GWTBaseURL = "https://cronometer.com/app"

	// APIExportURL is the full URL for requesting data exports.
	APIExportURL = "https://cronometer.com/export"
)

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

	if opts.GWTContentType == "" {
		c.GWTContentType = opts.GWTContentType
	}
	if opts.GWTModuleBase == "" {
		c.GWTModuleBase = opts.GWTModuleBase
	}
	if opts.GWTPermutation == "" {
		c.GWTPermutation = opts.GWTPermutation
	}
	if opts.GWTHeader == "" {
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
	}

	client.updateOpts(opts)

	return client
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
