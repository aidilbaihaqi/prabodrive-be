package utils

import (
	"fmt"

	"github.com/gabriel-vasile/mimetype"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
)

// blockedMIMETypes contains only dangerous file types that should never be uploaded.
// Everything else is allowed — this is a blocklist approach so that new file types
// (video, audio, archives, etc.) work without code changes.
var blockedMIMETypes = map[string]bool{
	// Executables
	"application/x-executable":    true,
	"application/x-msdownload":    true,
	"application/x-ms-installer":  true,
	"application/vnd.microsoft.portable-executable": true,
	"application/x-dosexec":       true,
	"application/x-elf":           true,
	"application/x-mach-binary":   true,
	// Scripts
	"application/x-sh":            true,
	"application/x-bash":          true,
	"application/x-csh":           true,
	"application/x-perl":          true,
	"application/x-python":        true,
	"application/x-ruby":          true,
	// HTML (XSS risk when served directly)
	"text/html":                   true,
	"application/xhtml+xml":       true,
	// Java archives
	"application/java-archive":    true,
	"application/x-java-archive":  true,
}

func IsAllowedMIME(mimeType string) bool {
	return !blockedMIMETypes[mimeType]
}

func ValidateMIME(data []byte, declared string) error {
	if blockedMIMETypes[declared] {
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

	if blockedMIMETypes[detectedStr] {
		return fmt.Errorf("%w: detected %s", domain.ErrMIMENotAllowed, detectedStr)
	}
	if detectedStr != declared {
		return fmt.Errorf("%w: declared %s, detected %s", domain.ErrMIMEMismatch, detectedStr)
	}
	return nil
}
