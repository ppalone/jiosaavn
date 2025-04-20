package jiosaavn

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
