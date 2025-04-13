package jiosaavn

import (
	"fmt"
	"strconv"
)

// Song.
type Song struct {
	ID              string
	Title           string
	Subtitle        string
	PermanentURL    string
	Image           string
	Language        string
	Year            string
	PlayCount       int
	ExplicitContent bool
	Music           string
	AlbumId         string
	AlbumName       string
	AlbumURL        string
	Label           string
	MediaURL        string
	Duration        int
	PrimaryArtists  []Artist
	FeaturedArtists []Artist
}

// Song API Response.
type songAPIResponse struct {
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
}

// Get Song API Response.
type getSongAPIResponse struct {
	Songs []songAPIResponse `json:"songs"`
}

func (res *songAPIResponse) toSong() Song {
	count, _ := strconv.Atoi(res.PlayCount)
	duration, _ := strconv.Atoi(res.MoreInfo.Duration)
	mediaURL, _ := generateMediaURL(res.MoreInfo.EncryptedMediaURL)
	song := Song{
		ID:              res.ID,
		Title:           res.Title,
		Subtitle:        res.Subtitle,
		PermanentURL:    res.PermaURL,
		Image:           res.Image,
		Language:        res.Language,
		Year:            res.Year,
		PlayCount:       count,
		ExplicitContent: res.ExplicitContent == "1",
		Music:           res.MoreInfo.Music,
		AlbumId:         res.MoreInfo.AlbumID,
		AlbumName:       res.MoreInfo.Album,
		AlbumURL:        res.MoreInfo.AlbumURL,
		Label:           res.MoreInfo.Label,
		MediaURL:        mediaURL,
		Duration:        duration,
	}

	primaryArtists := make([]Artist, 0)
	for _, artist := range res.MoreInfo.ArtistMap.PrimaryArtists {
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
	for _, artist := range res.MoreInfo.ArtistMap.FeaturedArtists {
		a := Artist{
			ID:           artist.ID,
			Name:         artist.Name,
			Image:        artist.Image,
			PermanentURL: artist.PermaURL,
		}

		featuredArtists = append(featuredArtists, a)
	}
	song.FeaturedArtists = featuredArtists

	return song
}

func (res *getSongAPIResponse) toSong() (Song, error) {
	if len(res.Songs) == 0 {
		return Song{}, fmt.Errorf("invalid song id")
	}

	return res.Songs[0].toSong(), nil
}
