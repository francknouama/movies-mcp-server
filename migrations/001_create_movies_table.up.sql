-- Create movies table with image support (SQLite version)
CREATE TABLE IF NOT EXISTS movies (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    director TEXT NOT NULL,
    year INTEGER NOT NULL CHECK (year >= 1888 AND year <= 2100),
    genre TEXT NOT NULL DEFAULT '[]', -- JSON array
    rating REAL CHECK (rating >= 0 AND rating <= 10),
    description TEXT,
    duration INTEGER CHECK (duration > 0), -- minutes
    language TEXT,
    country TEXT,
    poster_data BLOB, -- Binary image data
    poster_type TEXT, -- MIME type (e.g., 'image/jpeg')
    poster_url TEXT, -- URL to poster image
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_movies_title ON movies (title);
CREATE INDEX idx_movies_director ON movies (director);
CREATE INDEX idx_movies_year ON movies (year);
CREATE INDEX idx_movies_rating ON movies (rating DESC);
CREATE INDEX idx_movies_country ON movies (country);
CREATE INDEX idx_movies_language ON movies (language);

-- Create update trigger for updated_at
CREATE TRIGGER update_movies_updated_at
AFTER UPDATE ON movies
FOR EACH ROW
BEGIN
    UPDATE movies SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;