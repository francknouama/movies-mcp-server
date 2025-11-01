package dto

// CreateMovieRequest represents the MCP request to create a movie.
type CreateMovieRequest struct {
	Title     string   `json:"title"`
	Director  string   `json:"director"`
	Year      int      `json:"year"`
	Rating    float64  `json:"rating,omitempty"`
	Genres    []string `json:"genres,omitempty"`
	PosterURL string   `json:"poster_url,omitempty"`
}

// UpdateMovieRequest represents the MCP request to update a movie.
type UpdateMovieRequest struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Director  string   `json:"director"`
	Year      int      `json:"year"`
	Rating    float64  `json:"rating,omitempty"`
	Genres    []string `json:"genres,omitempty"`
	PosterURL string   `json:"poster_url,omitempty"`
}

// SearchMoviesRequest represents the MCP request to search movies.
type SearchMoviesRequest struct {
	Query     string  `json:"query,omitempty"`
	Title     string  `json:"title,omitempty"`
	Director  string  `json:"director,omitempty"`
	Genre     string  `json:"genre,omitempty"`
	MinYear   int     `json:"min_year,omitempty"`
	MaxYear   int     `json:"max_year,omitempty"`
	MinRating float64 `json:"min_rating,omitempty"`
	MaxRating float64 `json:"max_rating,omitempty"`
	Limit     int     `json:"limit,omitempty"`
	Offset    int     `json:"offset,omitempty"`
	OrderBy   string  `json:"order_by,omitempty"`
	OrderDir  string  `json:"order_dir,omitempty"`
}

// SearchByDecadeRequest represents the MCP request to search movies by decade.
type SearchByDecadeRequest struct {
	Decade string `json:"decade"`
}

// SearchByRatingRangeRequest represents the MCP request to search movies by rating range.
type SearchByRatingRangeRequest struct {
	MinRating float64 `json:"min_rating,omitempty"`
	MaxRating float64 `json:"max_rating,omitempty"`
}

// SearchSimilarMoviesRequest represents the MCP request to find similar movies.
type SearchSimilarMoviesRequest struct {
	MovieID int `json:"movie_id"`
	Limit   int `json:"limit,omitempty"`
}

// MovieResponse represents the MCP response for a movie.
type MovieResponse struct {
	ID        int      `json:"id"`
	Title     string   `json:"title"`
	Director  string   `json:"director"`
	Year      int      `json:"year"`
	Rating    float64  `json:"rating,omitempty"`
	Genres    []string `json:"genres"`
	PosterURL string   `json:"poster_url,omitempty"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
}

// MoviesListResponse represents the MCP response for a list of movies.
type MoviesListResponse struct {
	Movies      []*MovieResponse `json:"movies"`
	Total       int              `json:"total,omitempty"`
	Description string           `json:"description,omitempty"`
}
