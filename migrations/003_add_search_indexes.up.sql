-- Add indexes for search optimization (SQLite version)

-- Index for title search (case-insensitive using COLLATE NOCASE)
CREATE INDEX IF NOT EXISTS idx_movies_title_lower ON movies (title COLLATE NOCASE);

-- Index for director search (case-insensitive using COLLATE NOCASE)
CREATE INDEX IF NOT EXISTS idx_movies_director_lower ON movies (director COLLATE NOCASE);

-- Index for year search
CREATE INDEX IF NOT EXISTS idx_movies_year ON movies (year);

-- Index for rating (for sorting top movies)
CREATE INDEX IF NOT EXISTS idx_movies_rating ON movies (rating DESC);

-- Composite index for common queries
CREATE INDEX IF NOT EXISTS idx_movies_year_rating ON movies (year, rating DESC);
