package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"regexp"
)

var unsafeChars = regexp.MustCompile(`[^a-zA-Z0-9._-]`)

func SanitizeFilename(name string) string {
	return unsafeChars.ReplaceAllString(name, "_")
}

func GenerateToken(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func SHA256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
