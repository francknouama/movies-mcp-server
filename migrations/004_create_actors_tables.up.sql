-- Create actors table (SQLite version)
CREATE TABLE IF NOT EXISTS actors (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    birth_date TEXT,
    biography TEXT,
    photo_url TEXT,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    updated_at TEXT DEFAULT CURRENT_TIMESTAMP
);

-- Create movie_actors junction table
CREATE TABLE IF NOT EXISTS movie_actors (
    movie_id INTEGER NOT NULL,
    actor_id INTEGER NOT NULL,
    role TEXT,
    billing_order INTEGER,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (movie_id, actor_id),
    FOREIGN KEY (movie_id) REFERENCES movies(id) ON DELETE CASCADE,
    FOREIGN KEY (actor_id) REFERENCES actors(id) ON DELETE CASCADE
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_actors_name ON actors(name);
CREATE INDEX IF NOT EXISTS idx_movie_actors_movie_id ON movie_actors(movie_id);
CREATE INDEX IF NOT EXISTS idx_movie_actors_actor_id ON movie_actors(actor_id);
CREATE INDEX IF NOT EXISTS idx_movie_actors_billing_order ON movie_actors(billing_order);

-- Add trigger to update updated_at timestamp
CREATE TRIGGER update_actors_updated_at
AFTER UPDATE ON actors
FOR EACH ROW
BEGIN
    UPDATE actors SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
