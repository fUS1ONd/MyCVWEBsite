package domain

import "time"

// MediaFile represents an uploaded media file (image or video).
type MediaFile struct {
	ID          int       `json:"id" db:"id"`
	Filename    string    `json:"filename" db:"filename"`
	MimeType    string    `json:"mime_type" db:"mime_type"`
	Size        int64     `json:"size" db:"size"`
	UploaderID  int       `json:"uploader_id" db:"uploader_id"`
	StoragePath string    `json:"-" db:"storage_path"`
	URL         string    `json:"url" db:"-"`
	UploadedAt  time.Time `json:"uploaded_at" db:"uploaded_at"`
}

// UploadMediaRequest represents a request to upload a media file.
type UploadMediaRequest struct {
	Filename string
	MimeType string
	Size     int64
}
