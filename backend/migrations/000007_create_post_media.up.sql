CREATE TABLE IF NOT EXISTS post_media (
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    media_id INTEGER NOT NULL REFERENCES media_files(id) ON DELETE CASCADE,
    sort_order INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (post_id, media_id)
);

CREATE INDEX idx_post_media_post_id ON post_media(post_id);
CREATE INDEX idx_post_media_media_id ON post_media(media_id);
CREATE INDEX idx_post_media_sort_order ON post_media(post_id, sort_order);
