// Package dto provides data transfer objects for the movies MCP server.
package dto

// CreateActorRequest represents the MCP request to create an actor.
type CreateActorRequest struct {
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	Bio       string `json:"bio,omitempty"`
}

// UpdateActorRequest represents the MCP request to update an actor.
type UpdateActorRequest struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	Bio       string `json:"bio,omitempty"`
}

// LinkActorToMovieRequest represents the MCP request to link an actor to a movie.
type LinkActorToMovieRequest struct {
	ActorID int `json:"actor_id"`
	MovieID int `json:"movie_id"`
}

// SearchActorsRequest represents the MCP request to search actors.
type SearchActorsRequest struct {
	Name         string `json:"name,omitempty"`
	MinBirthYear int    `json:"min_birth_year,omitempty"`
	MaxBirthYear int    `json:"max_birth_year,omitempty"`
	MovieID      int    `json:"movie_id,omitempty"`
	Limit        int    `json:"limit,omitempty"`
	Offset       int    `json:"offset,omitempty"`
	OrderBy      string `json:"order_by,omitempty"`
	OrderDir     string `json:"order_dir,omitempty"`
}

// ActorResponse represents the MCP response for an actor.
type ActorResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	BirthYear int    `json:"birth_year"`
	Bio       string `json:"bio,omitempty"`
	MovieIDs  []int  `json:"movie_ids"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ActorsListResponse represents the MCP response for a list of actors.
type ActorsListResponse struct {
	Actors      []*ActorResponse `json:"actors"`
	Total       int              `json:"total,omitempty"`
	Description string           `json:"description,omitempty"`
}
