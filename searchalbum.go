package jiosaavn

import (
	"context"
	"fmt"
	"strconv"
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
	Total   int `json:"total"`
	Start   int `json:"start"`
	Results []struct {
		ID              string `json:"id"`
		Title           string `json:"title"`
		Subtitle        string `json:"subtitle"`
		HeaderDesc      string `json:"header_desc"`
		Type            string `json:"type"`
		PermaURL        string `json:"perma_url"`
		Image           string `json:"image"`
		Language        string `json:"language"`
		Year            string `json:"year"`
		PlayCount       string `json:"play_count"`
		ExplicitContent string `json:"explicit_content"`
		ListCount       string `json:"list_count"`
		ListType        string `json:"list_type"`
		List            string `json:"list"`
		MoreInfo        struct {
			Query     string `json:"query"`
			Text      string `json:"text"`
			Music     string `json:"music"`
			SongCount string `json:"song_count"`
			ArtistMap struct {
				PrimaryArtists  []artistAPIResponse `json:"primary_artists"`
				FeaturedArtists []artistAPIResponse `json:"featured_artists"`
				Artists         []artistAPIResponse `json:"artists"`
			} `json:"artistMap"`
		} `json:"more_info"`
	} `json:"results"`
}

func (res *searchAlbumAPIResponse) toResults(c *Client, opts *searchOptions) (SearchAlbumsResults, error) {
	albums := make([]Album, 0)

	for _, result := range res.Results {
		year, _ := strconv.Atoi(result.Year)
		playCount, _ := strconv.Atoi(result.PlayCount)
		songCount, _ := strconv.Atoi(result.MoreInfo.SongCount)
		album := Album{
			ID:           result.ID,
			Title:        result.Title,
			Subtitle:     result.Subtitle,
			PermanentURL: result.PermaURL,
			Image:        result.Image,
			Language:     result.Language,
			Year:         year,
			PlayCount:    playCount,
			SongCount:    songCount,
		}

		primaryArtists := make([]Artist, 0)
		for _, artist := range result.MoreInfo.ArtistMap.PrimaryArtists {
			primaryArtists = append(primaryArtists, artist.toArtist())
		}
		album.PrimaryArtists = primaryArtists

		featuredArtists := make([]Artist, 0)
		for _, artist := range result.MoreInfo.ArtistMap.FeaturedArtists {
			featuredArtists = append(featuredArtists, artist.toArtist())
		}
		album.FeaturedArtists = featuredArtists

		albums = append(albums, album)
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
