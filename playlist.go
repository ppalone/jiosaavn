package jiosaavn

import (
	"fmt"
	"html"
	"strconv"
)

// Playlist.
type Playlist struct {
	ID              string
	Title           string
	Image           string
	PermanentURL    string
	SongCount       int
	Language        string
	ExplicitContent bool
}

// Playlist Info.
type PlaylistInfo struct {
	ID              string
	Title           string
	Image           string
	PermanentURL    string
	SongCount       int
	PlayCount       int
	Language        string
	ExplicitContent bool
	Songs           []Song
	Artists         []Artist
}

// Get Playlist API Response.
type getPlaylistAPIResponse struct {
	ID              string            `json:"id"`
	Title           string            `json:"title"`
	Subtitle        string            `json:"subtitle"`
	HeaderDesc      string            `json:"header_desc"`
	Type            string            `json:"type"`
	PermaURL        string            `json:"perma_url"`
	Image           string            `json:"image"`
	Language        string            `json:"language"`
	Year            string            `json:"year"`
	PlayCount       string            `json:"play_count"`
	ExplicitContent string            `json:"explicit_content"`
	ListCount       string            `json:"list_count"`
	ListType        string            `json:"list_type"`
	List            []songAPIResponse `json:"list"`
	MoreInfo        struct {
		UID            string              `json:"uid"`
		Contents       string              `json:"contents"`
		IsDolbyContent bool                `json:"is_dolby_content"`
		LastUpdated    string              `json:"last_updated"`
		Username       string              `json:"username"`
		Firstname      string              `json:"firstname"`
		Lastname       string              `json:"lastname"`
		IsFollowed     string              `json:"is_followed"`
		IsFY           bool                `json:"isFY"`
		FollowerCount  string              `json:"follower_count"`
		FanCount       string              `json:"fan_count"`
		PlaylistType   string              `json:"playlist_type"`
		Share          string              `json:"share"`
		VideoCount     string              `json:"video_count"`
		Artists        []artistAPIResponse `json:"artists"`
		SubtitleDesc   []string            `json:"subtitle_desc"`
	} `json:"more_info"`
}

func (res *getPlaylistAPIResponse) toPlaylistInfo() (PlaylistInfo, error) {
	if len(res.Title) == 0 && len(res.List) == 0 {
		return PlaylistInfo{}, fmt.Errorf("invalid playlist id")
	}

	songCount, _ := strconv.Atoi(res.ListCount)
	playCount, _ := strconv.Atoi(res.PlayCount)
	playlistInfo := PlaylistInfo{
		ID:              res.ID,
		Title:           html.UnescapeString(res.Title),
		Image:           res.Image,
		PermanentURL:    res.PermaURL,
		SongCount:       songCount,
		PlayCount:       playCount,
		Language:        res.Language,
		ExplicitContent: res.ExplicitContent == "1",
	}

	songs := make([]Song, 0)
	for _, s := range res.List {
		songs = append(songs, s.toSong())
	}
	playlistInfo.Songs = songs

	artists := make([]Artist, 0)
	for _, a := range res.MoreInfo.Artists {
		artists = append(artists, a.toArtist())
	}
	playlistInfo.Artists = artists

	return playlistInfo, nil
}
