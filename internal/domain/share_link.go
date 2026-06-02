package domain

import "time"

type ShareLink struct {
	ID           string
	DocumentID   string
	Token        string
	PasswordHash *string
	ExpiresAt    time.Time
	CreatedBy    string
	CreatedAt    time.Time
}
