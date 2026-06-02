package utils

import (
	"fmt"

	"github.com/gabriel-vasile/mimetype"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
)

var allowedMIMETypes = map[string]bool{
	"application/pdf": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
	"application/msword":                                                       true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":       true,
	"application/vnd.ms-excel":                                                 true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
	"application/vnd.ms-powerpoint": true,
	"text/plain":                    true,
	"image/jpeg":                    true,
	"image/png":                     true,
	"image/webp":                    true,
}

func IsAllowedMIME(mimeType string) bool {
	return allowedMIMETypes[mimeType]
}

func ValidateMIME(data []byte, declared string) error {
	if !allowedMIMETypes[declared] {
		return fmt.Errorf("%w: %s", domain.ErrMIMENotAllowed, declared)
	}

	detected := mimetype.Detect(data)
	detectedStr := detected.String()

	// Strip parameters (e.g. "text/plain; charset=utf-8" → "text/plain")
	for i, c := range detectedStr {
		if c == ';' {
			detectedStr = detectedStr[:i]
			break
		}
	}

	if !allowedMIMETypes[detectedStr] {
		return fmt.Errorf("%w: detected %s", domain.ErrMIMENotAllowed, detectedStr)
	}
	if detectedStr != declared {
		return fmt.Errorf("%w: declared %s, detected %s", domain.ErrMIMEMismatch, declared, detectedStr)
	}
	return nil
}
