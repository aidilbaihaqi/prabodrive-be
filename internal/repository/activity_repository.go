package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
)

type activityRepo struct {
	db *pgxpool.Pool
}

func NewActivityRepository(db *pgxpool.Pool) domain.ActivityRepository {
	return &activityRepo{db: db}
}

func (r *activityRepo) Log(ctx context.Context, entry *domain.ActivityLog) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO activity_logs (id, user_id, action, document_id, document_name, ip_address, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		entry.ID, entry.UserID, entry.Action, entry.DocumentID, entry.DocumentName, entry.IPAddress, entry.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("activityRepo.Log: %w", err)
	}
	return nil
}

func (r *activityRepo) List(ctx context.Context, userID string, page, limit int) ([]*domain.ActivityLog, int, error) {
	var total int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM activity_logs WHERE user_id = $1`, userID).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("activityRepo.List count: %w", err)
	}

	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx,
		`SELECT id, user_id, action, document_id, document_name, ip_address::text, created_at
		 FROM activity_logs WHERE user_id = $1
		 ORDER BY created_at DESC LIMIT $2 OFFSET $3`,
		userID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("activityRepo.List query: %w", err)
	}
	defer rows.Close()

	var logs []*domain.ActivityLog
	for rows.Next() {
		l := &domain.ActivityLog{}
		if err := rows.Scan(&l.ID, &l.UserID, &l.Action, &l.DocumentID, &l.DocumentName, &l.IPAddress, &l.CreatedAt); err != nil {
			return nil, 0, fmt.Errorf("activityRepo.List scan: %w", err)
		}
		logs = append(logs, l)
	}
	return logs, total, rows.Err()
}
