package jiosaavn

import "net/http"

// Client.
type Client struct {
	httpClient *http.Client
}

// NewClient returns a new JioSaavn client
func NewClient(c *http.Client) *Client {
	if c == nil {
		c = &http.Client{}
	}

	return &Client{c}
}
