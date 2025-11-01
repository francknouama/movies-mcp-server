-- Drop trigger
DROP TRIGGER IF EXISTS update_actors_updated_at;

-- Drop indexes
DROP INDEX IF EXISTS idx_movie_actors_billing_order;
DROP INDEX IF EXISTS idx_movie_actors_actor_id;
DROP INDEX IF EXISTS idx_movie_actors_movie_id;
DROP INDEX IF EXISTS idx_actors_name;

-- Drop tables
DROP TABLE IF EXISTS movie_actors;
DROP TABLE IF EXISTS actors;
