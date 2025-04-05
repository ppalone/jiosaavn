package jiosaavn

import (
	"context"
	"fmt"
	"strconv"
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
			Music                string `json:"music"`
			AlbumID              string `json:"album_id"`
			Album                string `json:"album"`
			Label                string `json:"label"`
			LabelID              string `json:"label_id"`
			Origin               string `json:"origin"`
			IsDolbyContent       bool   `json:"is_dolby_content"`
			Three20Kbps          string `json:"320kbps"`
			EncryptedMediaURL    string `json:"encrypted_media_url"`
			EncryptedCacheURL    string `json:"encrypted_cache_url"`
			EncryptedDrmCacheURL string `json:"encrypted_drm_cache_url"`
			EncryptedDrmMediaURL string `json:"encrypted_drm_media_url"`
			AlbumURL             string `json:"album_url"`
			Duration             string `json:"duration"`
			Rights               struct {
				Code               string `json:"code"`
				Cacheable          string `json:"cacheable"`
				DeleteCachedObject string `json:"delete_cached_object"`
				Reason             string `json:"reason"`
			} `json:"rights"`
			CacheState    string `json:"cache_state"`
			HasLyrics     string `json:"has_lyrics"`
			LyricsSnippet string `json:"lyrics_snippet"`
			Starred       string `json:"starred"`
			CopyrightText string `json:"copyright_text"`
			ArtistMap     struct {
				PrimaryArtists []struct {
					ID       string `json:"id"`
					Name     string `json:"name"`
					Role     string `json:"role"`
					Image    string `json:"image"`
					Type     string `json:"type"`
					PermaURL string `json:"perma_url"`
				} `json:"primary_artists"`
				FeaturedArtists []struct {
					ID       string `json:"id"`
					Name     string `json:"name"`
					Role     string `json:"role"`
					Image    string `json:"image"`
					Type     string `json:"type"`
					PermaURL string `json:"perma_url"`
				} `json:"featured_artists"`
				Artists []struct {
					ID       string `json:"id"`
					Name     string `json:"name"`
					Role     string `json:"role"`
					Image    string `json:"image"`
					Type     string `json:"type"`
					PermaURL string `json:"perma_url"`
				} `json:"artists"`
			} `json:"artistMap"`
			ReleaseDate        string `json:"release_date"`
			LabelURL           string `json:"label_url"`
			Vcode              string `json:"vcode"`
			Vlink              string `json:"vlink"`
			TrillerAvailable   bool   `json:"triller_available"`
			RequestJiotuneFlag bool   `json:"request_jiotune_flag"`
			Webp               string `json:"webp"`
		} `json:"more_info"`
	} `json:"results"`
}

func (resp *searchSongsAPIResponse) toResults(c *Client, opts *searchOptions) (SearchSongsResults, error) {
	songs := make([]Song, 0)

	for _, result := range resp.Results {
		count, _ := strconv.Atoi(result.PlayCount)
		duration, _ := strconv.Atoi(result.MoreInfo.Duration)
		mediaURL, _ := generateMediaURL(result.MoreInfo.EncryptedMediaURL)
		song := Song{
			ID:              result.ID,
			Title:           result.Title,
			Subtitle:        result.Subtitle,
			PermanentURL:    result.PermaURL,
			Image:           result.Image,
			Language:        result.Language,
			Year:            result.Year,
			PlayCount:       count,
			ExplicitContent: result.ExplicitContent == "1",
			Music:           result.MoreInfo.Music,
			AlbumId:         result.MoreInfo.AlbumID,
			AlbumName:       result.MoreInfo.Album,
			AlbumURL:        result.MoreInfo.AlbumURL,
			Label:           result.MoreInfo.Label,
			MediaURL:        mediaURL,
			Duration:        duration,
		}

		primaryArtists := make([]Artist, 0)
		for _, artist := range result.MoreInfo.ArtistMap.PrimaryArtists {
			a := Artist{
				ID:           artist.ID,
				Name:         artist.Name,
				Image:        artist.Image,
				PermanentURL: artist.PermaURL,
			}

			primaryArtists = append(primaryArtists, a)
		}
		song.PrimaryArtists = primaryArtists

		featuredArtists := make([]Artist, 0)
		for _, artist := range result.MoreInfo.ArtistMap.FeaturedArtists {
			a := Artist{
				ID:           artist.ID,
				Name:         artist.Name,
				Image:        artist.Image,
				PermanentURL: artist.PermaURL,
			}

			featuredArtists = append(featuredArtists, a)
		}
		song.FeaturedArtists = featuredArtists

		songs = append(songs, song)
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
