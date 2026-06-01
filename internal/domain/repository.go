package domain

import "github.com/google/uuid"

// ===========================================
// Repository Interfaces
// ===========================================
// Interfaces are defined in domain layer
// Implementations are in repository layer

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *User) error
	FindByID(id uuid.UUID) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAll(filters UserFilters) ([]*User, int64, error)
	Update(user *User) error
	Delete(id uuid.UUID) error
	Restore(id uuid.UUID) error
	ExistsByEmail(email string) (bool, error)
}

// UserFilters contains filter options for listing users
type UserFilters struct {
	Name           string
	Email          string
	Role           string
	IsActive       *bool
	IncludeDeleted bool
	Page           int
	Limit          int
	SortBy         string
	SortOrder      string // "asc" or "desc"
}

// DefaultUserFilters returns default filter values
func DefaultUserFilters() UserFilters {
	return UserFilters{
		Page:      1,
		Limit:     10,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
}
