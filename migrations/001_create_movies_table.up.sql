-- Create movies table with image support
CREATE TABLE IF NOT EXISTS movies (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    director VARCHAR(255) NOT NULL,
    year INTEGER NOT NULL CHECK (year >= 1888 AND year <= 2100),
    genre TEXT[] NOT NULL DEFAULT '{}',
    rating DECIMAL(3,1) CHECK (rating >= 0 AND rating <= 10),
    description TEXT,
    duration INTEGER CHECK (duration > 0), -- minutes
    language VARCHAR(50),
    country VARCHAR(100),
    poster_data BYTEA, -- Binary image data
    poster_type VARCHAR(50), -- MIME type (e.g., 'image/jpeg')
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_movies_title ON movies USING GIN (to_tsvector('english', title));
CREATE INDEX idx_movies_director ON movies (director);
CREATE INDEX idx_movies_year ON movies (year);
CREATE INDEX idx_movies_genre ON movies USING GIN (genre);
CREATE INDEX idx_movies_rating ON movies (rating DESC);
CREATE INDEX idx_movies_country ON movies (country);
CREATE INDEX idx_movies_language ON movies (language);

-- Full-text search index
CREATE INDEX idx_movies_search ON movies USING GIN (
    to_tsvector('english', 
        COALESCE(title, '') || ' ' || 
        COALESCE(director, '') || ' ' || 
        COALESCE(description, '')
    )
);

-- Create update trigger for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_movies_updated_at BEFORE UPDATE
    ON movies FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();