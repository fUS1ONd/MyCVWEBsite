DROP INDEX IF EXISTS idx_comments_likes_count;

ALTER TABLE comments
  DROP COLUMN IF EXISTS likes_count;
