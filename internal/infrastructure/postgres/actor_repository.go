package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"movies-mcp-server/internal/domain/actor"
	"movies-mcp-server/internal/domain/shared"
)

// ActorRepository implements the actor.Repository interface for PostgreSQL
type ActorRepository struct {
	db *sql.DB
}

// NewActorRepository creates a new PostgreSQL actor repository
func NewActorRepository(db *sql.DB) *ActorRepository {
	return &ActorRepository{db: db}
}

// dbActor represents the database model for actors
type dbActor struct {
	ID        int            `db:"id"`
	Name      string         `db:"name"`
	BirthYear sql.NullInt64  `db:"birth_year"`
	Bio       sql.NullString `db:"bio"`
	CreatedAt sql.NullTime   `db:"created_at"`
	UpdatedAt sql.NullTime   `db:"updated_at"`
}

// Save persists an actor (insert or update)
func (r *ActorRepository) Save(ctx context.Context, domainActor *actor.Actor) error {
	dbActor := r.toDBModel(domainActor)
	
	if domainActor.ID().IsZero() {
		return r.insert(ctx, dbActor, domainActor)
	}
	return r.update(ctx, dbActor, domainActor)
}

func (r *ActorRepository) insert(ctx context.Context, dbActor *dbActor, domainActor *actor.Actor) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Insert actor
	query := `
		INSERT INTO actors (name, birth_year, bio, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	
	var id int
	err = tx.QueryRowContext(ctx, query,
		dbActor.Name,
		dbActor.BirthYear,
		dbActor.Bio,
		dbActor.CreatedAt.Time,
		dbActor.UpdatedAt.Time,
	).Scan(&id)
	
	if err != nil {
		return fmt.Errorf("failed to insert actor: %w", err)
	}
	
	// Update domain actor with the new ID
	actorID, err := shared.NewActorID(id)
	if err != nil {
		return fmt.Errorf("failed to create actor ID: %w", err)
	}
	domainActor.SetID(actorID)
	
	// Insert movie relationships
	if err := r.insertMovieRelationships(ctx, tx, domainActor); err != nil {
		return fmt.Errorf("failed to insert movie relationships: %w", err)
	}
	
	return tx.Commit()
}

func (r *ActorRepository) update(ctx context.Context, dbActor *dbActor, domainActor *actor.Actor) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Update actor
	query := `
		UPDATE actors 
		SET name = $2, birth_year = $3, bio = $4, updated_at = $5
		WHERE id = $1`
	
	result, err := tx.ExecContext(ctx, query,
		domainActor.ID().Value(),
		dbActor.Name,
		dbActor.BirthYear,
		dbActor.Bio,
		dbActor.UpdatedAt.Time,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update actor: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("actor not found")
	}
	
	// Update movie relationships
	if err := r.updateMovieRelationships(ctx, tx, domainActor); err != nil {
		return fmt.Errorf("failed to update movie relationships: %w", err)
	}
	
	return tx.Commit()
}

func (r *ActorRepository) insertMovieRelationships(ctx context.Context, tx *sql.Tx, domainActor *actor.Actor) error {
	for _, movieID := range domainActor.MovieIDs() {
		query := `
			INSERT INTO movie_actors (movie_id, actor_id, created_at)
			VALUES ($1, $2, CURRENT_TIMESTAMP)
			ON CONFLICT (movie_id, actor_id) DO NOTHING`
		
		_, err := tx.ExecContext(ctx, query, movieID.Value(), domainActor.ID().Value())
		if err != nil {
			return fmt.Errorf("failed to insert movie relationship: %w", err)
		}
	}
	return nil
}

func (r *ActorRepository) updateMovieRelationships(ctx context.Context, tx *sql.Tx, domainActor *actor.Actor) error {
	// Delete existing relationships
	deleteQuery := "DELETE FROM movie_actors WHERE actor_id = $1"
	_, err := tx.ExecContext(ctx, deleteQuery, domainActor.ID().Value())
	if err != nil {
		return fmt.Errorf("failed to delete existing movie relationships: %w", err)
	}
	
	// Insert new relationships
	return r.insertMovieRelationships(ctx, tx, domainActor)
}

// FindByID retrieves an actor by their ID
func (r *ActorRepository) FindByID(ctx context.Context, id shared.ActorID) (*actor.Actor, error) {
	query := `
		SELECT id, name, birth_year, bio, created_at, updated_at
		FROM actors 
		WHERE id = $1`
	
	var dbActor dbActor
	err := r.db.QueryRowContext(ctx, query, id.Value()).Scan(
		&dbActor.ID,
		&dbActor.Name,
		&dbActor.BirthYear,
		&dbActor.Bio,
		&dbActor.CreatedAt,
		&dbActor.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("actor not found")
		}
		return nil, fmt.Errorf("failed to find actor: %w", err)
	}
	
	// Get movie relationships
	movieIDs, err := r.getActorMovieIDs(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get actor movie IDs: %w", err)
	}
	
	return r.toDomainModel(&dbActor, movieIDs)
}

func (r *ActorRepository) getActorMovieIDs(ctx context.Context, actorID shared.ActorID) ([]shared.MovieID, error) {
	query := "SELECT movie_id FROM movie_actors WHERE actor_id = $1"
	rows, err := r.db.QueryContext(ctx, query, actorID.Value())
	if err != nil {
		return nil, fmt.Errorf("failed to query movie relationships: %w", err)
	}
	defer rows.Close()
	
	var movieIDs []shared.MovieID
	for rows.Next() {
		var movieIDValue int
		if err := rows.Scan(&movieIDValue); err != nil {
			return nil, fmt.Errorf("failed to scan movie ID: %w", err)
		}
		
		movieID, err := shared.NewMovieID(movieIDValue)
		if err != nil {
			return nil, fmt.Errorf("failed to create movie ID: %w", err)
		}
		
		movieIDs = append(movieIDs, movieID)
	}
	
	return movieIDs, nil
}

// FindByCriteria retrieves actors based on search criteria
func (r *ActorRepository) FindByCriteria(ctx context.Context, criteria actor.SearchCriteria) ([]*actor.Actor, error) {
	query, args := r.buildSearchQuery(criteria)
	
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to search actors: %w", err)
	}
	defer rows.Close()
	
	var actors []*actor.Actor
	for rows.Next() {
		var dbActor dbActor
		err := rows.Scan(
			&dbActor.ID,
			&dbActor.Name,
			&dbActor.BirthYear,
			&dbActor.Bio,
			&dbActor.CreatedAt,
			&dbActor.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan actor: %w", err)
		}
		
		actorID, err := shared.NewActorID(dbActor.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to create actor ID: %w", err)
		}
		
		// Get movie relationships
		movieIDs, err := r.getActorMovieIDs(ctx, actorID)
		if err != nil {
			return nil, fmt.Errorf("failed to get actor movie IDs: %w", err)
		}
		
		domainActor, err := r.toDomainModel(&dbActor, movieIDs)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to domain model: %w", err)
		}
		
		actors = append(actors, domainActor)
	}
	
	return actors, nil
}

func (r *ActorRepository) buildSearchQuery(criteria actor.SearchCriteria) (string, []interface{}) {
	query := `
		SELECT DISTINCT a.id, a.name, a.birth_year, a.bio, a.created_at, a.updated_at
		FROM actors a`
	
	var args []interface{}
	argIndex := 1
	var conditions []string
	
	// Join with movie_actors if searching by movie
	if !criteria.MovieID.IsZero() {
		query += " INNER JOIN movie_actors ma ON a.id = ma.actor_id"
		conditions = append(conditions, fmt.Sprintf("ma.movie_id = $%d", argIndex))
		args = append(args, criteria.MovieID.Value())
		argIndex++
	}
	
	// Add WHERE conditions
	if criteria.Name != "" {
		conditions = append(conditions, fmt.Sprintf("a.name ILIKE $%d", argIndex))
		args = append(args, "%"+criteria.Name+"%")
		argIndex++
	}
	
	if criteria.MinBirthYear > 0 {
		conditions = append(conditions, fmt.Sprintf("a.birth_year >= $%d", argIndex))
		args = append(args, criteria.MinBirthYear)
		argIndex++
	}
	
	if criteria.MaxBirthYear > 0 {
		conditions = append(conditions, fmt.Sprintf("a.birth_year <= $%d", argIndex))
		args = append(args, criteria.MaxBirthYear)
		argIndex++
	}
	
	if len(conditions) > 0 {
		query += " WHERE " + fmt.Sprintf("%s", conditions[0])
		for i := 1; i < len(conditions); i++ {
			query += " AND " + conditions[i]
		}
	}
	
	// Add ORDER BY
	orderField := "a.name"
	switch criteria.OrderBy {
	case actor.OrderByBirthYear:
		orderField = "a.birth_year"
	case actor.OrderByCreatedAt:
		orderField = "a.created_at"
	case actor.OrderByUpdatedAt:
		orderField = "a.updated_at"
	}
	
	orderDir := "ASC"
	if criteria.OrderDir == actor.OrderDesc {
		orderDir = "DESC"
	}
	
	query += fmt.Sprintf(" ORDER BY %s %s", orderField, orderDir)
	
	// Add LIMIT and OFFSET
	if criteria.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, criteria.Limit)
		argIndex++
	}
	
	if criteria.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, criteria.Offset)
	}
	
	return query, args
}

// FindByName searches actors by name (partial match)
func (r *ActorRepository) FindByName(ctx context.Context, name string) ([]*actor.Actor, error) {
	criteria := actor.SearchCriteria{
		Name:  name,
		Limit: 100, // Default limit
	}
	return r.FindByCriteria(ctx, criteria)
}

// FindByMovieID retrieves actors who appeared in a specific movie
func (r *ActorRepository) FindByMovieID(ctx context.Context, movieID shared.MovieID) ([]*actor.Actor, error) {
	criteria := actor.SearchCriteria{
		MovieID: movieID,
		Limit:   100, // Default limit
	}
	return r.FindByCriteria(ctx, criteria)
}

// CountAll returns the total number of actors
func (r *ActorRepository) CountAll(ctx context.Context) (int, error) {
	query := "SELECT COUNT(*) FROM actors"
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count actors: %w", err)
	}
	return count, nil
}

// Delete removes an actor by ID
func (r *ActorRepository) Delete(ctx context.Context, id shared.ActorID) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Delete movie relationships first (foreign key constraints)
	deleteRelQuery := "DELETE FROM movie_actors WHERE actor_id = $1"
	_, err = tx.ExecContext(ctx, deleteRelQuery, id.Value())
	if err != nil {
		return fmt.Errorf("failed to delete actor movie relationships: %w", err)
	}
	
	// Delete actor
	deleteQuery := "DELETE FROM actors WHERE id = $1"
	result, err := tx.ExecContext(ctx, deleteQuery, id.Value())
	if err != nil {
		return fmt.Errorf("failed to delete actor: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("actor not found")
	}
	
	return tx.Commit()
}

// DeleteAll removes all actors (for testing)
func (r *ActorRepository) DeleteAll(ctx context.Context) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Delete all movie relationships first
	_, err = tx.ExecContext(ctx, "DELETE FROM movie_actors")
	if err != nil {
		return fmt.Errorf("failed to delete all movie relationships: %w", err)
	}
	
	// Delete all actors
	_, err = tx.ExecContext(ctx, "DELETE FROM actors")
	if err != nil {
		return fmt.Errorf("failed to delete all actors: %w", err)
	}
	
	return tx.Commit()
}

// toDBModel converts a domain actor to a database model
func (r *ActorRepository) toDBModel(domainActor *actor.Actor) *dbActor {
	dbActor := &dbActor{
		ID:   domainActor.ID().Value(),
		Name: domainActor.Name(),
	}
	
	// Handle optional birth year
	if !domainActor.BirthYear().IsZero() {
		dbActor.BirthYear = sql.NullInt64{
			Int64: int64(domainActor.BirthYear().Value()),
			Valid: true,
		}
	}
	
	// Handle optional bio
	if domainActor.Bio() != "" {
		dbActor.Bio = sql.NullString{
			String: domainActor.Bio(),
			Valid:  true,
		}
	}
	
	// Handle timestamps
	dbActor.CreatedAt = sql.NullTime{
		Time:  domainActor.CreatedAt(),
		Valid: true,
	}
	dbActor.UpdatedAt = sql.NullTime{
		Time:  domainActor.UpdatedAt(),
		Valid: true,
	}
	
	return dbActor
}

// toDomainModel converts a database model to a domain actor
func (r *ActorRepository) toDomainModel(dbActor *dbActor, movieIDs []shared.MovieID) (*actor.Actor, error) {
	actorID, err := shared.NewActorID(dbActor.ID)
	if err != nil {
		return nil, fmt.Errorf("invalid actor ID: %w", err)
	}
	
	birthYear := 1900 // Default if not set
	if dbActor.BirthYear.Valid {
		birthYear = int(dbActor.BirthYear.Int64)
	}
	
	domainActor, err := actor.NewActorWithID(actorID, dbActor.Name, birthYear)
	if err != nil {
		return nil, fmt.Errorf("failed to create domain actor: %w", err)
	}
	
	// Set bio if present
	if dbActor.Bio.Valid {
		domainActor.SetBio(dbActor.Bio.String)
	}
	
	// Add movie relationships
	for _, movieID := range movieIDs {
		if err := domainActor.AddMovie(movieID); err != nil {
			return nil, fmt.Errorf("failed to add movie to actor: %w", err)
		}
	}
	
	return domainActor, nil
}