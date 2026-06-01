package user

import (
	"github.com/google/uuid"

	"github.com/yourname/yourapp/internal/domain"
)

// DeleteUserUsecase handles user deletion business logic
type DeleteUserUsecase struct {
	userRepo domain.UserRepository
}

// NewDeleteUserUsecase creates a new DeleteUserUsecase instance
func NewDeleteUserUsecase(userRepo domain.UserRepository) *DeleteUserUsecase {
	return &DeleteUserUsecase{
		userRepo: userRepo,
	}
}

// DeleteUserInput is the input DTO
type DeleteUserInput struct {
	ID uuid.UUID
}

// Execute performs the user deletion (soft delete)
func (uc *DeleteUserUsecase) Execute(input DeleteUserInput) error {
	// 1. Check if user exists
	_, err := uc.userRepo.FindByID(input.ID)
	if err != nil {
		return err
	}

	// 2. Soft delete user
	return uc.userRepo.Delete(input.ID)
}
