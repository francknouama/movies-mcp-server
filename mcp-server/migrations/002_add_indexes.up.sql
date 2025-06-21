-- Additional performance indexes

-- Composite index for common queries
CREATE INDEX idx_movies_year_rating ON movies (year, rating DESC);

-- Index for poster queries
CREATE INDEX idx_movies_has_poster ON movies ((poster_data IS NOT NULL));

-- Index for updated_at to find recently modified movies
CREATE INDEX idx_movies_updated_at ON movies (updated_at DESC);

-- Partial index for high-rated movies (rating >= 8.0)
CREATE INDEX idx_movies_high_rated ON movies (rating DESC) 
WHERE rating >= 8.0;

-- Index for director-year combinations
CREATE INDEX idx_movies_director_year ON movies (director, year);