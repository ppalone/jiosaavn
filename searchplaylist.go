package jiosaavn

import (
	"context"
	"fmt"
	"strconv"
)

// Search playlists results.
type SearchPlaylistsResults struct {
	Page      int
	Size      int
	Total     int
	HasNext   bool
	Playlists []Playlist

	// for next
	c             *Client
	searchOptions *searchOptions
}

// Search playlists API Response.
type searchPlaylistsAPIResponse struct {
	Total   int `json:"total"`
	Start   int `json:"start"`
	Results []struct {
		ID       string `json:"id"`
		Title    string `json:"title"`
		Subtitle string `json:"subtitle"`
		Type     string `json:"type"`
		Image    string `json:"image"`
		PermaURL string `json:"perma_url"`
		MoreInfo struct {
			UID            string `json:"uid"`
			Firstname      string `json:"firstname"`
			EntityType     string `json:"entity_type"`
			EntitySubType  string `json:"entity_sub_type"`
			VideoAvailable bool   `json:"video_available"`
			Lastname       string `json:"lastname"`
			SongCount      string `json:"song_count"`
			Language       string `json:"language"`
		} `json:"more_info"`
		ExplicitContent string `json:"explicit_content"`
		MiniObj         bool   `json:"mini_obj"`
	} `json:"results"`
}

func (resp *searchPlaylistsAPIResponse) toResults(c *Client, opts *searchOptions) (SearchPlaylistsResults, error) {
	playlists := make([]Playlist, 0)

	for _, result := range resp.Results {
		count, _ := strconv.Atoi(result.MoreInfo.SongCount)

		playlists = append(playlists, Playlist{
			ID:              result.ID,
			Title:           result.Title,
			Image:           result.Image,
			PermanentURL:    result.PermaURL,
			SongCount:       count,
			Language:        result.MoreInfo.Language,
			ExplicitContent: result.ExplicitContent == "1",
		})
	}

	hasNext := ((resp.Start - 1) + len(resp.Results)) < resp.Total
	if !hasNext {
		return SearchPlaylistsResults{
			Size:      len(playlists),
			Page:      opts.page,
			HasNext:   hasNext,
			Total:     resp.Total,
			Playlists: playlists,
		}, nil
	}

	return SearchPlaylistsResults{
		Size:          len(playlists),
		Page:          opts.page,
		HasNext:       hasNext,
		Total:         resp.Total,
		Playlists:     playlists,
		c:             c,
		searchOptions: opts,
	}, nil
}

func (results *SearchPlaylistsResults) Next(ctx context.Context) (SearchPlaylistsResults, error) {
	if !results.HasNext {
		return SearchPlaylistsResults{}, fmt.Errorf("doesn't have further results")
	}

	// next page is available
	results.searchOptions.page += 1
	return results.c.searchPlaylists(ctx, results.searchOptions.query, results.searchOptions)
}
