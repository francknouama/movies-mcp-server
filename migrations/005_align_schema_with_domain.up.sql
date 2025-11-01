-- Align database schema with new domain model (SQLite version)

-- Update actors table to match domain model
ALTER TABLE actors ADD COLUMN birth_year INTEGER;
ALTER TABLE actors ADD COLUMN bio TEXT;

-- Update existing birth_date data to birth_year if needed
UPDATE actors 
SET birth_year = CAST(strftime('%Y', birth_date) AS INTEGER)
WHERE birth_date IS NOT NULL AND birth_year IS NULL;

-- The actors table structure should now support our domain model:
-- id, name, birth_year, bio, created_at, updated_at

-- Movies table already has poster_url from migration 001

-- Add indexes for new fields
CREATE INDEX IF NOT EXISTS idx_actors_birth_year ON actors(birth_year);

-- Note: SQLite doesn't support adding CHECK constraints to existing tables
-- The constraint will be enforced at the application level

-- SQLite doesn't support COMMENT ON, but we'll document the schema here:
-- movies: Movies with full metadata and image support
-- actors: Actors with birth year and biography  
-- movie_actors: Many-to-many relationship between movies and actors
-- movies.genre: JSON array of genres for the movie
-- movies.poster_data: Binary image data for movie poster
-- movies.poster_type: MIME type of poster image
-- movies.poster_url: URL to poster image (alternative to poster_data)
-- actors.birth_year: Year the actor was born
-- actors.bio: Biography of the actor
