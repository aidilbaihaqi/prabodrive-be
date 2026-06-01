package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/yourname/yourapp/internal/domain"
)

// ListUsersUsecase handles listing users with pagination
type ListUsersUsecase struct {
	userRepo domain.UserRepository
}

// NewListUsersUsecase creates a new ListUsersUsecase instance
func NewListUsersUsecase(userRepo domain.UserRepository) *ListUsersUsecase {
	return &ListUsersUsecase{
		userRepo: userRepo,
	}
}

// ListUsersInput is the input DTO
type ListUsersInput struct {
	Name           string
	Email          string
	Role           string
	IsActive       *bool
	IncludeDeleted bool
	Page           int
	Limit          int
	SortBy         string
	SortOrder      string
}

// ListUsersOutput is the output DTO
type ListUsersOutput struct {
	Users      []UserItem `json:"users"`
	Pagination Pagination `json:"pagination"`
}

// UserItem represents a user in the list
type UserItem struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// Pagination contains pagination metadata
type Pagination struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

// Execute retrieves a paginated list of users
func (uc *ListUsersUsecase) Execute(input ListUsersInput) (*ListUsersOutput, error) {
	// Set defaults
	if input.Page < 1 {
		input.Page = 1
	}
	if input.Limit < 1 {
		input.Limit = 10
	}
	if input.Limit > 100 {
		input.Limit = 100
	}

	// Build filters
	filters := domain.UserFilters{
		Name:           input.Name,
		Email:          input.Email,
		Role:           input.Role,
		IsActive:       input.IsActive,
		IncludeDeleted: input.IncludeDeleted,
		Page:           input.Page,
		Limit:          input.Limit,
		SortBy:         input.SortBy,
		SortOrder:      input.SortOrder,
	}

	// Get users from repository
	users, total, err := uc.userRepo.FindAll(filters)
	if err != nil {
		return nil, err
	}

	// Convert to output format
	items := make([]UserItem, len(users))
	for i, user := range users {
		items[i] = UserItem{
			ID:        user.ID,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
			IsActive:  user.IsActive,
			CreatedAt: user.CreatedAt,
		}
	}

	// Calculate total pages
	totalPages := int(total) / input.Limit
	if int(total)%input.Limit > 0 {
		totalPages++
	}

	return &ListUsersOutput{
		Users: items,
		Pagination: Pagination{
			Page:       input.Page,
			Limit:      input.Limit,
			Total:      total,
			TotalPages: totalPages,
		},
	}, nil
}
