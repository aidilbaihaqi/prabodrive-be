package services

import (
	"context"
	"fmt"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
)

type QuotaService struct {
	users domain.UserRepository
}

func NewQuotaService(users domain.UserRepository) *QuotaService {
	return &QuotaService{users: users}
}

func (q *QuotaService) Check(ctx context.Context, userID string, fileSize int64) error {
	user, err := q.users.FindByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("quota check: %w", err)
	}
	if user.QuotaUsed+fileSize > user.QuotaMax {
		return domain.ErrQuotaExceeded
	}
	return nil
}
