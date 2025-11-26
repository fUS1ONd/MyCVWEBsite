CREATE TABLE IF NOT EXISTS media_files (
    id SERIAL PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    path VARCHAR(500) NOT NULL,
    mime_type VARCHAR(100) NOT NULL,
    size BIGINT NOT NULL,
    uploader_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_media_files_uploader_id ON media_files(uploader_id);
CREATE INDEX idx_media_files_created_at ON media_files(created_at DESC);
CREATE INDEX idx_media_files_mime_type ON media_files(mime_type);
