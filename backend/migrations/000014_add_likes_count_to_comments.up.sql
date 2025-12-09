ALTER TABLE comments
  ADD COLUMN likes_count INTEGER DEFAULT 0;

CREATE INDEX idx_comments_likes_count ON comments(likes_count DESC);
