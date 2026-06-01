package domain

import (
	"time"

	"github.com/google/uuid"
)

// User represents the user entity in the domain layer
// This is a pure business entity without any framework dependencies
type User struct {
	ID        uuid.UUID
	Email     string
	Password  string // hashed password
	Name      string
	Role      string
	IsActive  bool
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

// Validate performs business rule validation for User
func (u *User) Validate() error {
	if u.Email == "" {
		return ErrEmailRequired
	}
	if !isValidEmail(u.Email) {
		return ErrInvalidEmail
	}
	if u.Name == "" {
		return ErrNameRequired
	}
	if len(u.Name) < 2 {
		return ErrNameTooShort
	}
	if u.Password == "" {
		return ErrPasswordRequired
	}
	return nil
}

// IsAdmin checks if user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == "admin"
}

// CanAccess checks if user can access a resource
func (u *User) CanAccess(resourceOwnerID uuid.UUID) bool {
	return u.IsAdmin() || u.ID == resourceOwnerID
}

// Deactivate soft deletes the user
func (u *User) Deactivate() {
	now := time.Now()
	u.IsActive = false
	u.DeletedAt = &now
}

// Activate restores a soft-deleted user
func (u *User) Activate() {
	u.IsActive = true
	u.DeletedAt = nil
}

// Helper function for email validation
func isValidEmail(email string) bool {
	// Simple email validation - in production use a proper regex or library
	if len(email) < 3 || len(email) > 254 {
		return false
	}
	atIndex := -1
	for i, c := range email {
		if c == '@' {
			atIndex = i
			break
		}
	}
	return atIndex > 0 && atIndex < len(email)-1
}
