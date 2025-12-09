DROP INDEX IF EXISTS idx_posts_likes_count;
DROP INDEX IF EXISTS idx_posts_cover_image_id;

ALTER TABLE posts
  DROP COLUMN IF EXISTS comments_count,
  DROP COLUMN IF EXISTS likes_count,
  DROP COLUMN IF EXISTS read_time_minutes,
  DROP COLUMN IF EXISTS cover_image_id;
