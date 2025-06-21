-- Drop trigger and function
DROP TRIGGER IF EXISTS update_movies_updated_at ON movies;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS idx_movies_search;
DROP INDEX IF EXISTS idx_movies_language;
DROP INDEX IF EXISTS idx_movies_country;
DROP INDEX IF EXISTS idx_movies_rating;
DROP INDEX IF EXISTS idx_movies_genre;
DROP INDEX IF EXISTS idx_movies_year;
DROP INDEX IF EXISTS idx_movies_director;
DROP INDEX IF EXISTS idx_movies_title;

-- Drop table
DROP TABLE IF EXISTS movies;