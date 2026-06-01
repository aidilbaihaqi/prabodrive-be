package user

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/yourname/yourapp/internal/domain"
)

// CreateUserUsecase handles user creation business logic
type CreateUserUsecase struct {
	userRepo domain.UserRepository
}

// NewCreateUserUsecase creates a new CreateUserUsecase instance
func NewCreateUserUsecase(userRepo domain.UserRepository) *CreateUserUsecase {
	return &CreateUserUsecase{
		userRepo: userRepo,
	}
}

// Input DTO for creating a user
type CreateUserInput struct {
	Email    string
	Password string
	Name     string
	Role     string
}

// Output DTO after creating a user
type CreateUserOutput struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Execute performs the user creation
func (uc *CreateUserUsecase) Execute(input CreateUserInput) (*CreateUserOutput, error) {
	// 1. Validate password length
	if len(input.Password) < 8 {
		return nil, domain.ErrPasswordTooShort
	}

	// 2. Check if email already exists
	exists, err := uc.userRepo.ExistsByEmail(input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailExists
	}

	// 3. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 4. Set default role if not provided
	role := input.Role
	if role == "" {
		role = "user"
	}

	// 5. Create domain entity
	user := &domain.User{
		ID:        uuid.New(),
		Email:     input.Email,
		Password:  string(hashedPassword),
		Name:      input.Name,
		Role:      role,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 6. Validate entity
	if err := user.Validate(); err != nil {
		return nil, err
	}

	// 7. Save to repository
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 8. Return output
	return &CreateUserOutput{
		ID:        user.ID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}
