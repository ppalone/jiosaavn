package jiosaavn

import (
	"context"
	"encoding/json"
	"fmt"
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

// SearchSongs
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

// SearchPlaylists
func (c *Client) SearchPlaylists(ctx context.Context, q string, opts ...SearchOption) (SearchPlaylistsResults, error) {
	searchOpts := defaultSearchOpts()
	for _, opt := range opts {
		opt(searchOpts)
	}

	return c.searchPlaylists(ctx, q, searchOpts)
}

// SearchAlbums
func (c *Client) SearchAlbums(ctx context.Context, q string, opts ...SearchOption) (SearchAlbumsResults, error) {
	searchOpts := defaultSearchOpts()
	for _, opt := range opts {
		opt(searchOpts)
	}

	return c.searchAlbums(ctx, q, searchOpts)
}

// GetSongById
func (c *Client) GetSongById(ctx context.Context, id string) (Song, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return Song{}, fmt.Errorf("song id cannot be empty")
	}

	params := make(map[string]string)
	params["pids"] = id
	params[callEndpoint] = getSongById

	apiResponse := new(getSongAPIResponse)
	err := c.makeRequestAndUnmarshal(ctx, params, apiResponse)
	if err != nil {
		return Song{}, err
	}

	return apiResponse.toSong()
}

// GetPlaylistById
func (c *Client) GetPlaylistById(ctx context.Context, id string) (PlaylistInfo, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return PlaylistInfo{}, fmt.Errorf("playlist id cannot be empty")
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		return PlaylistInfo{}, fmt.Errorf("playlist id must be a number")
	}

	// TODO: add p(page) and n(limit) pagination
	params := make(map[string]string)
	params["listid"] = id
	params[callEndpoint] = getPlaylistById

	apiResponse := new(getPlaylistAPIResponse)
	err = c.makeRequestAndUnmarshal(ctx, params, apiResponse)
	if err != nil {
		return PlaylistInfo{}, err
	}

	return apiResponse.toPlaylistInfo()
}

// GetAlbumById
func (c *Client) GetAlbumById(ctx context.Context, id string) (AlbumInfo, error) {
	id = strings.TrimSpace(id)
	if len(id) == 0 {
		return AlbumInfo{}, fmt.Errorf("album id cannot be empty")
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		return AlbumInfo{}, fmt.Errorf("album id must be a number")
	}

	params := make(map[string]string)
	params["albumid"] = id
	params[callEndpoint] = getAlbumById

	apiResponse := new(getAlbumAPIResponse)
	err = c.makeRequestAndUnmarshal(ctx, params, apiResponse)
	if err != nil {
		return AlbumInfo{}, err
	}

	return apiResponse.toAlbumInfo()
}

func (c *Client) searchSongs(ctx context.Context, q string, opts *searchOptions) (SearchSongsResults, error) {
	opts.query = strings.TrimSpace(q)

	params, err := buildSearchParams(opts)
	if err != nil {
		return SearchSongsResults{}, err
	}
	params[callEndpoint] = searchSongsEndpoint

	apiResp := new(searchSongsAPIResponse)
	err = c.makeRequestAndUnmarshal(ctx, params, apiResp)
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

	apiResp := new(searchArtistsAPIResponse)
	err = c.makeRequestAndUnmarshal(ctx, params, apiResp)
	if err != nil {
		return SearchArtistsResults{}, err
	}

	return apiResp.toResults(c, opts)
}

func (c *Client) searchPlaylists(ctx context.Context, q string, opts *searchOptions) (SearchPlaylistsResults, error) {
	opts.query = strings.TrimSpace(q)
	params, err := buildSearchParams(opts)
	if err != nil {
		return SearchPlaylistsResults{}, err
	}
	params[callEndpoint] = searchPlaylistsEndpoint

	apiResponse := new(searchPlaylistsAPIResponse)
	err = c.makeRequestAndUnmarshal(ctx, params, apiResponse)
	if err != nil {
		return SearchPlaylistsResults{}, err
	}

	return apiResponse.toResults(c, opts)
}

func (c *Client) searchAlbums(ctx context.Context, q string, opts *searchOptions) (SearchAlbumsResults, error) {
	opts.query = strings.TrimSpace(q)
	params, err := buildSearchParams(opts)
	if err != nil {
		return SearchAlbumsResults{}, err
	}
	params[callEndpoint] = searchAlbumsEndpoint

	apiResponse := new(searchAlbumAPIResponse)
	err = c.makeRequestAndUnmarshal(ctx, params, apiResponse)
	if err != nil {
		return SearchAlbumsResults{}, err
	}

	return apiResponse.toResults(c, opts)
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

func (c *Client) makeRequestAndUnmarshal(ctx context.Context, params map[string]string, v any) error {
	req, err := makeRequest(ctx, params)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(v)
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
