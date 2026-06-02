package domain

import "time"

type Folder struct {
	ID        string
	UserID    string
	ParentID  *string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
