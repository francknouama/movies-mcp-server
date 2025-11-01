-- Drop trigger
DROP TRIGGER IF EXISTS update_movies_updated_at;

-- Drop indexes
DROP INDEX IF EXISTS idx_movies_language;
DROP INDEX IF EXISTS idx_movies_country;
DROP INDEX IF EXISTS idx_movies_rating;
DROP INDEX IF EXISTS idx_movies_year;
DROP INDEX IF EXISTS idx_movies_director;
DROP INDEX IF EXISTS idx_movies_title;

-- Drop table
DROP TABLE IF EXISTS movies;
