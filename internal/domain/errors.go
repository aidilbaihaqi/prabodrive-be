package domain

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrUserNotFound     = errors.New("user not found")
	ErrDocumentNotFound = errors.New("document not found")
	ErrFolderNotFound   = errors.New("folder not found")
	ErrShareNotFound    = errors.New("share link not found")

	ErrEmailExists      = errors.New("email already registered")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrPasswordTooShort = errors.New("password must be at least 8 characters")

	ErrUnauthorized      = errors.New("unauthorized")
	ErrForbidden         = errors.New("forbidden")
	ErrOwnership         = errors.New("resource belongs to another user")
	ErrShareExpired      = errors.New("share link has expired")
	ErrSharePasswordWrong = errors.New("incorrect share link password")

	ErrQuotaExceeded = errors.New("quota exceeded")
	ErrFileTooLarge  = errors.New("file exceeds 5 MB limit")
	ErrMIMENotAllowed = errors.New("file type not allowed")
	ErrMIMEMismatch  = errors.New("declared MIME type does not match file content")

	ErrInvalidToken = errors.New("invalid or expired token")
)
