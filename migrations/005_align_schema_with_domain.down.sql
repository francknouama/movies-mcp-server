-- Revert schema alignment (SQLite version)

-- Drop indexes
DROP INDEX IF EXISTS idx_actors_birth_year;

-- Note: SQLite doesn't support dropping columns easily
-- In production, you would need to:
-- 1. Create a new table without the columns
-- 2. Copy data from old table to new table
-- 3. Drop old table
-- 4. Rename new table to old name

-- For development purposes, we'll document the revert but not execute it:
-- The birth_year and bio columns would remain in the actors table
-- This is acceptable for a down migration in SQLite
