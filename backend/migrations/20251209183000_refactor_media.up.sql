ALTER TABLE posts ADD COLUMN cover_image VARCHAR(1024);
ALTER TABLE posts DROP COLUMN cover_image_id;
