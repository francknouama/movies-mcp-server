-- Reverse the schema alignment changes

-- Remove new columns from actors table
ALTER TABLE actors 
    DROP COLUMN IF EXISTS birth_year,
    DROP COLUMN IF EXISTS bio;

-- Remove new columns from movies table  
ALTER TABLE movies 
    DROP COLUMN IF EXISTS poster_url;

-- Remove new indexes
DROP INDEX IF EXISTS idx_actors_birth_year;

-- Remove check constraint
ALTER TABLE actors 
    DROP CONSTRAINT IF EXISTS chk_actors_birth_year;

-- Remove comments
COMMENT ON TABLE movies IS NULL;
COMMENT ON TABLE actors IS NULL;
COMMENT ON TABLE movie_actors IS NULL;
COMMENT ON COLUMN movies.genre IS NULL;
COMMENT ON COLUMN movies.poster_data IS NULL;
COMMENT ON COLUMN movies.poster_type IS NULL;
COMMENT ON COLUMN movies.poster_url IS NULL;
COMMENT ON COLUMN actors.birth_year IS NULL;
COMMENT ON COLUMN actors.bio IS NULL;