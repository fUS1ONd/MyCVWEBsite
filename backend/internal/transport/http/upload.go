package http

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

// uploadImage handles image upload for admin
func (h *Handler) uploadImage(w http.ResponseWriter, r *http.Request) {
	// 5MB limit
	r.Body = http.MaxBytesReader(w, r.Body, 5<<20)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		RespondBadRequest(w, "file too large or invalid format")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		RespondBadRequest(w, "invalid file")
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			h.log.Error("failed to close upload file", "error", err)
		}
	}()

	// Validate extension
	ext := strings.ToLower(filepath.Ext(header.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		RespondBadRequest(w, "unsupported file type (allowed: jpg, png, webp)")
		return
	}

	// Ensure upload dir exists
	uploadDir := "./uploads"
	if err := os.MkdirAll(uploadDir, 0750); err != nil {
		h.log.Error("failed to create upload dir", "error", err)
		RespondInternalError(w)
		return
	}

	// Generate filename
	newFilename := uuid.New().String() + ext
	dstPath := filepath.Join(uploadDir, newFilename)

	dst, err := os.Create(filepath.Clean(dstPath))
	if err != nil {
		h.log.Error("failed to create file", "error", err)
		RespondInternalError(w)
		return
	}
	defer func() {
		if err := dst.Close(); err != nil {
			h.log.Error("failed to close destination file", "error", err)
		}
	}()

	if _, err := io.Copy(dst, file); err != nil {
		h.log.Error("failed to save file", "error", err)
		RespondInternalError(w)
		return
	}

	RespondSuccess(w, map[string]string{
		"url": "/uploads/" + newFilename,
	})
}
