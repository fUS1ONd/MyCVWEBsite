-- Add name and avatar_url columns to users table
ALTER TABLE users
ADD COLUMN name VARCHAR(255) NOT NULL DEFAULT '',
ADD COLUMN avatar_url VARCHAR(500);

-- Add comment to explain the columns
COMMENT ON COLUMN users.name IS 'User full name from OAuth provider';
COMMENT ON COLUMN users.avatar_url IS 'User avatar URL from OAuth provider';
