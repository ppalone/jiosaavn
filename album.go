package jiosaavn

import (
	"fmt"
	"strconv"
)

// Album.
type Album struct {
	ID              string
	Title           string
	Subtitle        string
	PermanentURL    string
	Image           string
	Language        string
	Year            int
	PlayCount       int
	SongCount       int
	PrimaryArtists  []Artist
	FeaturedArtists []Artist
}

// Album Info
type AlbumInfo struct {
	Album
	Songs []Song
}

// Get Album API Response.
type getAlbumAPIResponse struct {
	ID              string   `json:"id"`
	Title           string   `json:"title"`
	Subtitle        string   `json:"subtitle"`
	HeaderDesc      string   `json:"header_desc"`
	Type            string   `json:"type"`
	PermaURL        string   `json:"perma_url"`
	Image           string   `json:"image"`
	Language        string   `json:"language"`
	Year            string   `json:"year"`
	PlayCount       string   `json:"play_count"`
	ExplicitContent string   `json:"explicit_content"`
	ListCount       string   `json:"list_count"`
	ListType        string   `json:"list_type"`
	List            songList `json:"list"`
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
}

func (res *getAlbumAPIResponse) toAlbum() Album {
	year, _ := strconv.Atoi(res.Year)
	playCount, _ := strconv.Atoi(res.PlayCount)
	songCount, _ := strconv.Atoi(res.MoreInfo.SongCount)
	album := Album{
		ID:           res.ID,
		Title:        res.Title,
		Subtitle:     res.Subtitle,
		PermanentURL: res.PermaURL,
		Image:        res.Image,
		Language:     res.Language,
		Year:         year,
		PlayCount:    playCount,
		SongCount:    songCount,
	}

	primaryArtists := make([]Artist, 0)
	for _, artist := range res.MoreInfo.ArtistMap.PrimaryArtists {
		primaryArtists = append(primaryArtists, artist.toArtist())
	}
	album.PrimaryArtists = primaryArtists

	featuredArtists := make([]Artist, 0)
	for _, artist := range res.MoreInfo.ArtistMap.FeaturedArtists {
		featuredArtists = append(featuredArtists, artist.toArtist())
	}
	album.FeaturedArtists = featuredArtists

	return album
}

func (res *getAlbumAPIResponse) toAlbumInfo() (AlbumInfo, error) {
	if len(res.Title) == 0 || len(res.List) == 0 {
		return AlbumInfo{}, fmt.Errorf("invalid album id")
	}

	album := res.toAlbum()

	songs := make([]Song, 0)
	for _, entry := range res.List {
		songs = append(songs, entry.toSong())
	}

	return AlbumInfo{
		Album: album,
		Songs: songs,
	}, nil
}
