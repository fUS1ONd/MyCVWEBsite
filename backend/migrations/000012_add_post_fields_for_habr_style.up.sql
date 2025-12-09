ALTER TABLE posts
  ADD COLUMN cover_image_id INTEGER REFERENCES media_files(id) ON DELETE SET NULL,
  ADD COLUMN read_time_minutes INTEGER DEFAULT 0,
  ADD COLUMN likes_count INTEGER DEFAULT 0,
  ADD COLUMN comments_count INTEGER DEFAULT 0;

CREATE INDEX idx_posts_cover_image_id ON posts(cover_image_id);
CREATE INDEX idx_posts_likes_count ON posts(likes_count DESC);
