-- Revert column renames
ALTER TABLE media_files RENAME COLUMN storage_path TO path;
ALTER TABLE media_files RENAME COLUMN uploaded_at TO created_at;

-- Revert index rename
ALTER INDEX idx_media_files_uploaded_at RENAME TO idx_media_files_created_at;
