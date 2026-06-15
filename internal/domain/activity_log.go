package domain

import "time"

type ActivityLog struct {
	ID           string
	UserID       *string
	Action       string
	DocumentID   *string
	DocumentName *string
	IPAddress    *string
	CreatedAt    time.Time
}
