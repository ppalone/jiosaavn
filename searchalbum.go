package jiosaavn

import (
	"context"
	"fmt"
)

// Search Album Results
type SearchAlbumsResults struct {
	Page    int
	Size    int
	Total   int
	HasNext bool
	Albums  []Album

	// for next
	c             *Client
	searchOptions *searchOptions
}

// Search Album API Response.
type searchAlbumAPIResponse struct {
	Total   int                   `json:"total"`
	Start   int                   `json:"start"`
	Results []getAlbumAPIResponse `json:"results"`
}

func (res *searchAlbumAPIResponse) toResults(c *Client, opts *searchOptions) (SearchAlbumsResults, error) {
	albums := make([]Album, 0)

	for _, result := range res.Results {
		albums = append(albums, result.toAlbum())
	}

	hasNext := ((res.Start - 1) + len(res.Results)) < res.Total
	if !hasNext {
		return SearchAlbumsResults{
			Page:    opts.page,
			Size:    len(albums),
			Total:   res.Total,
			Albums:  albums,
			HasNext: hasNext,
		}, nil
	}

	return SearchAlbumsResults{
		Page:          opts.page,
		Size:          len(albums),
		Total:         res.Total,
		Albums:        albums,
		HasNext:       hasNext,
		c:             c,
		searchOptions: opts,
	}, nil
}

func (results *SearchAlbumsResults) Next(ctx context.Context) (SearchAlbumsResults, error) {
	if !results.HasNext {
		return SearchAlbumsResults{}, fmt.Errorf("doesn't have further results")
	}

	// next page is available
	results.searchOptions.page += 1
	return results.c.searchAlbums(ctx, results.searchOptions.query, results.searchOptions)
}
