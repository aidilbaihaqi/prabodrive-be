package domain

import "time"

type ActivityLog struct {
	ID         string
	UserID     *string
	Action     string
	DocumentID *string
	IPAddress  *string
	CreatedAt  time.Time
}
