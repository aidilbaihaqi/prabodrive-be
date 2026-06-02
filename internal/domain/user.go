package domain

import "time"

type User struct {
	ID           string
	Email        string
	PasswordHash string
	Name         string
	Role         string
	QuotaUsed    int64
	QuotaMax     int64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
