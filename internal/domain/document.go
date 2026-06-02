package domain

import "time"

type Document struct {
	ID        string
	UserID    string
	FolderID  *string
	Name      string
	Size      int64
	MIMEType  string
	S3Key     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
