package usecase

import (
	"context"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
)

type AdminUsecase interface {
	ListUsers(ctx context.Context, page, limit int) ([]*domain.User, int, error)
	GetUser(ctx context.Context, id string) (*domain.User, error)
	UpdateRole(ctx context.Context, targetID, requesterID, role string) error
	DeleteUser(ctx context.Context, targetID, requesterID string) error
}

type adminUsecase struct {
	users domain.UserRepository
}

func NewAdminUsecase(users domain.UserRepository) AdminUsecase {
	return &adminUsecase{users: users}
}

func (u *adminUsecase) ListUsers(ctx context.Context, page, limit int) ([]*domain.User, int, error) {
	return u.users.ListAll(ctx, page, limit)
}

func (u *adminUsecase) GetUser(ctx context.Context, id string) (*domain.User, error) {
	return u.users.FindByID(ctx, id)
}

func (u *adminUsecase) UpdateRole(ctx context.Context, targetID, requesterID, role string) error {
	if targetID == requesterID {
		return domain.ErrForbidden
	}
	return u.users.UpdateRole(ctx, targetID, role)
}

func (u *adminUsecase) DeleteUser(ctx context.Context, targetID, requesterID string) error {
	if targetID == requesterID {
		return domain.ErrForbidden
	}
	return u.users.DeleteUser(ctx, targetID)
}
