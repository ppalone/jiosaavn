package jiosaavn

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
