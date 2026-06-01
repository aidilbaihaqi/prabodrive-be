package user

import (
	"time"

	"github.com/google/uuid"

	"github.com/yourname/yourapp/internal/domain"
)

// GetUserUsecase handles getting a single user
type GetUserUsecase struct {
	userRepo domain.UserRepository
}

// NewGetUserUsecase creates a new GetUserUsecase instance
func NewGetUserUsecase(userRepo domain.UserRepository) *GetUserUsecase {
	return &GetUserUsecase{
		userRepo: userRepo,
	}
}

// GetUserInput is the input DTO
type GetUserInput struct {
	ID uuid.UUID
}

// GetUserOutput is the output DTO
type GetUserOutput struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Execute retrieves a user by ID
func (uc *GetUserUsecase) Execute(input GetUserInput) (*GetUserOutput, error) {
	user, err := uc.userRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	return &GetUserOutput{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
