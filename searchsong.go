package jiosaavn

import (
	"context"
	"fmt"
)

// Search songs results.
type SearchSongsResults struct {
	Page    int
	Size    int
	Total   int
	HasNext bool
	Songs   []Song

	// for next
	c             *Client
	searchOptions *searchOptions
}

// Search songs API response.
type searchSongsAPIResponse struct {
	Total   int               `json:"total"`
	Start   int               `json:"start"`
	Results []songAPIResponse `json:"results"`
}

func (resp *searchSongsAPIResponse) toResults(c *Client, opts *searchOptions) (SearchSongsResults, error) {
	songs := make([]Song, 0)

	for _, result := range resp.Results {
		songs = append(songs, result.toSong())
	}

	hasNext := ((resp.Start - 1) + len(resp.Results)) < resp.Total
	if !hasNext {
		return SearchSongsResults{
			Page:    opts.page,
			Size:    len(songs),
			Total:   resp.Total,
			HasNext: hasNext,
			Songs:   songs,
		}, nil
	}

	return SearchSongsResults{
		Page:          opts.page,
		Size:          len(songs),
		Total:         resp.Total,
		HasNext:       hasNext,
		Songs:         songs,
		c:             c,
		searchOptions: opts,
	}, nil
}

func (results *SearchSongsResults) Next(ctx context.Context) (SearchSongsResults, error) {
	if !results.HasNext {
		return SearchSongsResults{}, fmt.Errorf("doesn't have further results")
	}

	// next page is available
	results.searchOptions.page += 1
	return results.c.searchSongs(ctx, results.searchOptions.query, results.searchOptions)
}
