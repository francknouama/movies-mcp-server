-- Drop search optimization indexes

DROP INDEX IF EXISTS idx_movies_title_lower;
DROP INDEX IF EXISTS idx_movies_director_lower;
DROP INDEX IF EXISTS idx_movies_year;
DROP INDEX IF EXISTS idx_movies_rating;
DROP INDEX IF EXISTS idx_movies_genre;
DROP INDEX IF EXISTS idx_movies_fulltext;
DROP INDEX IF EXISTS idx_movies_year_rating;