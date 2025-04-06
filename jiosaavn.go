package jiosaavn

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
)

// constants
const (
	baseURL      = "https://www.jiosaavn.com/api.php"
	callEndpoint = "__call"
)

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

// SearchSongs returns
func (c *Client) SearchSongs(ctx context.Context, q string, opts ...SearchOption) (SearchSongsResults, error) {
	searchOpts := defaultSearchOpts()
	for _, opt := range opts {
		opt(searchOpts)
	}

	return c.searchSongs(ctx, q, searchOpts)
}

// SearchArtists
func (c *Client) SearchArtists(ctx context.Context, q string, opts ...SearchOption) (SearchArtistsResults, error) {
	searchOpts := defaultSearchOpts()
	for _, opt := range opts {
		opt(searchOpts)
	}

	return c.searchArtists(ctx, q, searchOpts)
}

func (c *Client) searchSongs(ctx context.Context, q string, opts *searchOptions) (SearchSongsResults, error) {
	opts.query = strings.TrimSpace(q)

	params, err := buildSearchParams(opts)
	if err != nil {
		return SearchSongsResults{}, err
	}
	params[callEndpoint] = searchSongsEndpoint

	req, err := makeRequest(ctx, params)
	if err != nil {
		return SearchSongsResults{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SearchSongsResults{}, err
	}
	defer resp.Body.Close()

	apiResp := new(searchSongsAPIResponse)
	err = json.NewDecoder(resp.Body).Decode(apiResp)
	if err != nil {
		return SearchSongsResults{}, err
	}

	return apiResp.toResults(c, opts)
}

func (c *Client) searchArtists(ctx context.Context, q string, opts *searchOptions) (SearchArtistsResults, error) {
	opts.query = strings.TrimSpace(q)

	params, err := buildSearchParams(opts)
	if err != nil {
		return SearchArtistsResults{}, err
	}
	params[callEndpoint] = searchArtistsEndpoint

	req, err := makeRequest(ctx, params)
	if err != nil {
		return SearchArtistsResults{}, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return SearchArtistsResults{}, err
	}
	defer resp.Body.Close()

	apiResp := new(searchArtistsAPIResponse)
	err = json.NewDecoder(resp.Body).Decode(apiResp)
	if err != nil {
		return SearchArtistsResults{}, err
	}

	return apiResp.toResults(c, opts)
}

func makeRequest(ctx context.Context, params map[string]string) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, baseURL, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Set("_format", "json")
	q.Set("_marker", "0")
	q.Set("api_version", "4")
	q.Set("ctx", "web6dot0")

	for k, v := range params {
		q.Set(k, v)
	}

	req.URL.RawQuery = q.Encode()

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func buildSearchParams(opts *searchOptions) (map[string]string, error) {
	err := opts.validate()
	if err != nil {
		return nil, err
	}

	params := make(map[string]string)
	params["q"] = opts.query               // set search query
	params["p"] = strconv.Itoa(opts.page)  // set page
	params["n"] = strconv.Itoa(opts.limit) // set limit

	return params, nil
}
