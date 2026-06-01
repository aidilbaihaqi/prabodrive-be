package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/yourname/yourapp/internal/domain"
	"github.com/yourname/yourapp/internal/repository/models"
)

type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) domain.UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user in the database
func (r *userRepository) Create(user *domain.User) error {
	model := r.toModel(user)
	if err := r.db.Create(model).Error; err != nil {
		return r.handleError(err)
	}
	return nil
}

// FindByID finds a user by their ID
func (r *userRepository) FindByID(id uuid.UUID) (*domain.User, error) {
	var model models.User
	err := r.db.Where("id = ?", id).First(&model).Error
	if err != nil {
		return nil, r.handleError(err)
	}
	return r.toDomain(&model), nil
}

// FindByEmail finds a user by their email
func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	var model models.User
	err := r.db.Where("email = ?", email).First(&model).Error
	if err != nil {
		return nil, r.handleError(err)
	}
	return r.toDomain(&model), nil
}

// FindAll returns a paginated list of users with filters
func (r *userRepository) FindAll(filters domain.UserFilters) ([]*domain.User, int64, error) {
	var models []models.User
	var total int64

	query := r.db.Model(&models)

	// Apply filters
	if filters.Name != "" {
		query = query.Where("name ILIKE ?", "%"+filters.Name+"%")
	}
	if filters.Email != "" {
		query = query.Where("email ILIKE ?", "%"+filters.Email+"%")
	}
	if filters.Role != "" {
		query = query.Where("role = ?", filters.Role)
	}
	if filters.IsActive != nil {
		query = query.Where("is_active = ?", *filters.IsActive)
	}
	if !filters.IncludeDeleted {
		query = query.Where("deleted_at IS NULL")
	} else {
		query = query.Unscoped()
	}

	// Count total before pagination
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply sorting
	if filters.SortBy != "" {
		order := filters.SortBy
		if filters.SortOrder == "desc" {
			order += " DESC"
		} else {
			order += " ASC"
		}
		query = query.Order(order)
	}

	// Apply pagination
	if filters.Limit > 0 {
		offset := (filters.Page - 1) * filters.Limit
		query = query.Offset(offset).Limit(filters.Limit)
	}

	// Execute query
	if err := query.Find(&models).Error; err != nil {
		return nil, 0, err
	}

	// Convert to domain entities
	users := make([]*domain.User, len(models))
	for i, model := range models {
		users[i] = r.toDomain(&model)
	}

	return users, total, nil
}

// Update updates a user in the database
func (r *userRepository) Update(user *domain.User) error {
	return r.db.Model(&models.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]interface{}{
			"email":     user.Email,
			"password":  user.Password,
			"name":      user.Name,
			"role":      user.Role,
			"is_active": user.IsActive,
		}).Error
}

// Delete soft deletes a user
func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Where("id = ?", id).Delete(&models.User{}).Error
}

// Restore restores a soft-deleted user
func (r *userRepository) Restore(id uuid.UUID) error {
	return r.db.Model(&models.User{}).
		Unscoped().
		Where("id = ?", id).
		Update("deleted_at", nil).Error
}

// ExistsByEmail checks if a user with the given email exists
func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

// ===========================================
// Helper Methods
// ===========================================

// toModel converts domain entity to GORM model
func (r *userRepository) toModel(user *domain.User) *models.User {
	return &models.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
		Name:     user.Name,
		Role:     user.Role,
		IsActive: user.IsActive,
	}
}

// toDomain converts GORM model to domain entity
func (r *userRepository) toDomain(model *models.User) *domain.User {
	user := &domain.User{
		ID:        model.ID,
		Email:     model.Email,
		Password:  model.Password,
		Name:      model.Name,
		Role:      model.Role,
		IsActive:  model.IsActive,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
	if model.DeletedAt.Valid {
		user.DeletedAt = &model.DeletedAt.Time
	}
	return user
}

// handleError converts database errors to domain errors
func (r *userRepository) handleError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return domain.ErrUserNotFound
	}
	return err
}
