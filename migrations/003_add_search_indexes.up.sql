-- Add indexes for search optimization

-- Index for title search (case-insensitive)
CREATE INDEX IF NOT EXISTS idx_movies_title_lower ON movies (LOWER(title));

-- Index for director search (case-insensitive)
CREATE INDEX IF NOT EXISTS idx_movies_director_lower ON movies (LOWER(director));

-- Index for year search
CREATE INDEX IF NOT EXISTS idx_movies_year ON movies (year);

-- Index for rating (for sorting top movies)
CREATE INDEX IF NOT EXISTS idx_movies_rating ON movies (rating DESC);

-- GIN index for genre array search
CREATE INDEX IF NOT EXISTS idx_movies_genre ON movies USING GIN (genre);

-- Full-text search index
CREATE INDEX IF NOT EXISTS idx_movies_fulltext ON movies 
USING GIN (to_tsvector('english', title || ' ' || director || ' ' || COALESCE(description, '')));

-- Composite index for common queries
CREATE INDEX IF NOT EXISTS idx_movies_year_rating ON movies (year, rating DESC);