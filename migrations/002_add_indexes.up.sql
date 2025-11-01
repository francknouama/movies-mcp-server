-- Additional performance indexes (SQLite version)

-- Composite index for common queries
CREATE INDEX IF NOT EXISTS idx_movies_year_rating ON movies (year, rating DESC);

-- Index for updated_at to find recently modified movies
CREATE INDEX IF NOT EXISTS idx_movies_updated_at ON movies (updated_at DESC);

-- Partial index for high-rated movies (rating >= 8.0)
CREATE INDEX IF NOT EXISTS idx_movies_high_rated ON movies (rating DESC)
WHERE rating >= 8.0;

-- Index for director-year combinations
CREATE INDEX IF NOT EXISTS idx_movies_director_year ON movies (director, year);
