package jiosaavn

import (
	"context"
	"fmt"
	"html"
)

// Search artists results.
type SearchArtistsResults struct {
	Page    int
	Size    int
	Total   int
	HasNext bool
	Artists []Artist

	// for next
	c             *Client
	searchOptions *searchOptions
}

// Search artists API response.
type searchArtistsAPIResponse struct {
	Total   int `json:"total"`
	Start   int `json:"start"`
	Results []struct {
		Name           string `json:"name"`
		ID             string `json:"id"`
		Ctr            int    `json:"ctr"`
		Entity         int    `json:"entity"`
		Image          string `json:"image"`
		Role           string `json:"role"`
		PermaURL       string `json:"perma_url"`
		Type           string `json:"type"`
		MiniObj        bool   `json:"mini_obj"`
		IsRadioPresent bool   `json:"isRadioPresent"`
		IsFollowed     bool   `json:"is_followed"`
	} `json:"results"`
}

func (resp *searchArtistsAPIResponse) toResults(c *Client, opts *searchOptions) (SearchArtistsResults, error) {
	artists := make([]Artist, 0)

	for _, result := range resp.Results {
		artist := Artist{
			ID:           result.ID,
			Name:         html.EscapeString(result.Name),
			Image:        result.Image,
			PermanentURL: result.PermaURL,
		}

		artists = append(artists, artist)
	}

	hasNext := ((resp.Start - 1) + len(resp.Results)) < resp.Total
	if !hasNext {
		return SearchArtistsResults{
			Page:    opts.page,
			Size:    len(artists),
			Total:   resp.Total,
			Artists: artists,
			HasNext: hasNext,
		}, nil
	}

	return SearchArtistsResults{
		Page:          opts.page,
		Size:          len(artists),
		Total:         resp.Total,
		Artists:       artists,
		HasNext:       hasNext,
		c:             c,
		searchOptions: opts,
	}, nil
}

func (results *SearchArtistsResults) Next(ctx context.Context) (SearchArtistsResults, error) {
	if !results.HasNext {
		return SearchArtistsResults{}, fmt.Errorf("no further results")
	}

	results.searchOptions.page += 1
	return results.c.searchArtists(ctx, results.searchOptions.query, results.searchOptions)
}
