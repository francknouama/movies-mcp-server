-- Align database schema with new domain model

-- Update actors table to match domain model
ALTER TABLE actors 
    ADD COLUMN IF NOT EXISTS birth_year INTEGER,
    ADD COLUMN IF NOT EXISTS bio TEXT;

-- Update existing birth_date data to birth_year if needed
UPDATE actors 
SET birth_year = EXTRACT(YEAR FROM birth_date) 
WHERE birth_date IS NOT NULL AND birth_year IS NULL;

-- The actors table structure should now support our domain model:
-- id, name, birth_year, bio, created_at, updated_at

-- Movies table already has most fields we need, but let's ensure poster_url is available
ALTER TABLE movies 
    ADD COLUMN IF NOT EXISTS poster_url TEXT;

-- Update movie_actors table to be simpler (we'll handle this through our domain)
-- The existing structure is fine: movie_id, actor_id as primary key

-- Add indexes for new fields
CREATE INDEX IF NOT EXISTS idx_actors_birth_year ON actors(birth_year);

-- Add check constraint for birth_year
ALTER TABLE actors 
    ADD CONSTRAINT chk_actors_birth_year 
    CHECK (birth_year IS NULL OR (birth_year >= 1850 AND birth_year <= EXTRACT(YEAR FROM CURRENT_DATE)));

-- Comments for documentation
COMMENT ON TABLE movies IS 'Movies with full metadata and image support';
COMMENT ON TABLE actors IS 'Actors with birth year and biography';
COMMENT ON TABLE movie_actors IS 'Many-to-many relationship between movies and actors';

COMMENT ON COLUMN movies.genre IS 'Array of genres for the movie';
COMMENT ON COLUMN movies.poster_data IS 'Binary image data for movie poster';
COMMENT ON COLUMN movies.poster_type IS 'MIME type of poster image';
COMMENT ON COLUMN movies.poster_url IS 'URL to poster image (alternative to poster_data)';
COMMENT ON COLUMN actors.birth_year IS 'Year the actor was born';
COMMENT ON COLUMN actors.bio IS 'Biography of the actor';