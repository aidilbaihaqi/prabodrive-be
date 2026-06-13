package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
)

type userRepo struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) domain.UserRepository {
	return &userRepo{db: db}
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, name, role, quota_used, quota_max, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		u.ID, u.Email, u.PasswordHash, u.Name, u.Role, u.QuotaUsed, u.QuotaMax, u.CreatedAt, u.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("userRepo.Create: %w", err)
	}
	return nil
}

func (r *userRepo) FindByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, name, role, quota_used, quota_max, created_at, updated_at
		 FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, email, password_hash, name, role, quota_used, quota_max, created_at, updated_at
		 FROM users WHERE email = $1`, email)
	return scanUser(row)
}

func (r *userRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`, email).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("userRepo.ExistsByEmail: %w", err)
	}
	return exists, nil
}

func (r *userRepo) AddQuota(ctx context.Context, userID string, delta int64) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET quota_used = quota_used + $1, updated_at = $2 WHERE id = $3`,
		delta, time.Now(), userID,
	)
	if err != nil {
		return fmt.Errorf("userRepo.AddQuota: %w", err)
	}
	return nil
}

func (r *userRepo) ListAll(ctx context.Context, page, limit int) ([]*domain.User, int, error) {
	var total int
	if err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("userRepo.ListAll count: %w", err)
	}

	offset := (page - 1) * limit
	rows, err := r.db.Query(ctx,
		`SELECT id, email, password_hash, name, role, quota_used, quota_max, created_at, updated_at
		 FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("userRepo.ListAll query: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		u := &domain.User{}
		if err := rows.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role,
			&u.QuotaUsed, &u.QuotaMax, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("userRepo.ListAll scan: %w", err)
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

func (r *userRepo) UpdateProfile(ctx context.Context, id, name string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE users SET name = $1, updated_at = $2 WHERE id = $3`,
		name, time.Now(), id,
	)
	if err != nil {
		return fmt.Errorf("userRepo.UpdateProfile: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepo) UpdatePassword(ctx context.Context, id, passwordHash string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE users SET password_hash = $1, updated_at = $2 WHERE id = $3`,
		passwordHash, time.Now(), id,
	)
	if err != nil {
		return fmt.Errorf("userRepo.UpdatePassword: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepo) UpdateRole(ctx context.Context, id, role string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE users SET role = $1, updated_at = $2 WHERE id = $3`,
		role, time.Now(), id,
	)
	if err != nil {
		return fmt.Errorf("userRepo.UpdateRole: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *userRepo) DeleteUser(ctx context.Context, id string) error {
	tag, err := r.db.Exec(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("userRepo.DeleteUser: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func scanUser(row pgx.Row) (*domain.User, error) {
	u := &domain.User{}
	err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Name, &u.Role,
		&u.QuotaUsed, &u.QuotaMax, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("userRepo scan: %w", err)
	}
	return u, nil
}
