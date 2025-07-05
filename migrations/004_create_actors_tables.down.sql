-- Drop triggers
DROP TRIGGER IF EXISTS update_actors_updated_at ON actors;

-- Drop indexes
DROP INDEX IF EXISTS idx_movie_actors_billing_order;
DROP INDEX IF EXISTS idx_movie_actors_actor_id;
DROP INDEX IF EXISTS idx_movie_actors_movie_id;
DROP INDEX IF EXISTS idx_actors_name;

-- Drop tables
DROP TABLE IF EXISTS movie_actors;
DROP TABLE IF EXISTS actors;

-- Drop function if no other tables use it
-- Note: Only drop if this is the last table using this function
-- DROP FUNCTION IF EXISTS update_updated_at_column();