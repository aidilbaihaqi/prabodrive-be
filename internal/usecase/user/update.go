package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourname/yourapp/internal/domain"
)

// UpdateUserUsecase handles user update business logic
type UpdateUserUsecase struct {
	userRepo domain.UserRepository
}

// NewUpdateUserUsecase creates a new UpdateUserUsecase instance
func NewUpdateUserUsecase(userRepo domain.UserRepository) *UpdateUserUsecase {
	return &UpdateUserUsecase{
		userRepo: userRepo,
	}
}

// UpdateUserInput is the input DTO
type UpdateUserInput struct {
	ID       uuid.UUID
	Email    *string
	Password *string
	Name     *string
	Role     *string
	IsActive *bool
}

// UpdateUserOutput is the output DTO
type UpdateUserOutput struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"is_active"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Execute performs the user update
func (uc *UpdateUserUsecase) Execute(input UpdateUserInput) (*UpdateUserOutput, error) {
	// 1. Get existing user
	user, err := uc.userRepo.FindByID(input.ID)
	if err != nil {
		return nil, err
	}

	// 2. Update fields if provided
	if input.Email != nil {
		// Check if email is being changed and if new email exists
		if *input.Email != user.Email {
			exists, err := uc.userRepo.ExistsByEmail(*input.Email)
			if err != nil {
				return nil, err
			}
			if exists {
				return nil, domain.ErrEmailExists
			}
		}
		user.Email = *input.Email
	}

	if input.Password != nil {
		if len(*input.Password) < 8 {
			return nil, domain.ErrPasswordTooShort
		}
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*input.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.Password = string(hashedPassword)
	}

	if input.Name != nil {
		user.Name = *input.Name
	}

	if input.Role != nil {
		user.Role = *input.Role
	}

	if input.IsActive != nil {
		user.IsActive = *input.IsActive
	}

	user.UpdatedAt = time.Now()

	// 3. Validate entity
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// 4. Save to repository
	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	// 5. Return output
	return &UpdateUserOutput{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		IsActive:  user.IsActive,
		UpdatedAt: user.UpdatedAt,
	}, nil
}
