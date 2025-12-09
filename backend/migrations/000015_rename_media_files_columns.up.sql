-- Rename columns to match Go code expectations
ALTER TABLE media_files RENAME COLUMN path TO storage_path;
ALTER TABLE media_files RENAME COLUMN created_at TO uploaded_at;

-- Rename index to match new column name
ALTER INDEX idx_media_files_created_at RENAME TO idx_media_files_uploaded_at;
