package request

import "github.com/google/uuid"

// CreateUserRequest is the request body for creating a user
type CreateUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=2,max=100"`
	Role     string `json:"role" binding:"omitempty,oneof=user admin"`
}

// UpdateUserRequest is the request body for updating a user
type UpdateUserRequest struct {
	Email    *string `json:"email" binding:"omitempty,email"`
	Password *string `json:"password" binding:"omitempty,min=8"`
	Name     *string `json:"name" binding:"omitempty,min=2,max=100"`
	Role     *string `json:"role" binding:"omitempty,oneof=user admin"`
	IsActive *bool   `json:"is_active"`
}

// ListUsersRequest is the query parameters for listing users
type ListUsersRequest struct {
	Name           string `form:"name"`
	Email          string `form:"email"`
	Role           string `form:"role"`
	IsActive       *bool  `form:"is_active"`
	IncludeDeleted bool   `form:"include_deleted"`
	Page           int    `form:"page" binding:"omitempty,min=1"`
	Limit          int    `form:"limit" binding:"omitempty,min=1,max=100"`
	SortBy         string `form:"sort_by" binding:"omitempty,oneof=name email created_at updated_at"`
	SortOrder      string `form:"sort_order" binding:"omitempty,oneof=asc desc"`
}

// IDRequest is used for path parameter ID
type IDRequest struct {
	ID uuid.UUID `uri:"id" binding:"required"`
}

// ParseID parses string ID to UUID
func ParseID(id string) (uuid.UUID, error) {
	return uuid.Parse(id)
}
