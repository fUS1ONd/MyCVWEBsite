package domain

import "time"

// MediaFile represents an uploaded media file
type MediaFile struct {
	ID         int       `json:"id"`
	Filename   string    `json:"filename" validate:"required"`
	Path       string    `json:"path" validate:"required"`
	MimeType   string    `json:"mime_type" validate:"required"`
	Size       int64     `json:"size" validate:"required,gt=0"`
	UploaderID int       `json:"uploader_id" validate:"required"`
	CreatedAt  time.Time `json:"created_at"`
}
