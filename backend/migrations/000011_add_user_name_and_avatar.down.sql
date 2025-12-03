-- Remove name and avatar_url columns from users table
ALTER TABLE users
DROP COLUMN IF EXISTS name,
DROP COLUMN IF EXISTS avatar_url;
