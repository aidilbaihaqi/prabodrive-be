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

type tokenRepo struct {
	db *pgxpool.Pool
}

func NewRefreshTokenRepository(db *pgxpool.Pool) domain.RefreshTokenRepository {
	return &tokenRepo{db: db}
}

func (r *tokenRepo) Save(ctx context.Context, userID, tokenHash string, expiresAt interface{}) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, tokenHash, expiresAt,
	)
	if err != nil {
		return fmt.Errorf("tokenRepo.Save: %w", err)
	}
	return nil
}

func (r *tokenRepo) Find(ctx context.Context, tokenHash string) (string, error) {
	var userID string
	var expiresAt time.Time
	err := r.db.QueryRow(ctx,
		`SELECT user_id, expires_at FROM refresh_tokens WHERE token_hash = $1`, tokenHash,
	).Scan(&userID, &expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrInvalidToken
		}
		return "", fmt.Errorf("tokenRepo.Find: %w", err)
	}
	if time.Now().After(expiresAt) {
		return "", domain.ErrInvalidToken
	}
	return userID, nil
}

func (r *tokenRepo) Delete(ctx context.Context, tokenHash string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE token_hash = $1`, tokenHash)
	if err != nil {
		return fmt.Errorf("tokenRepo.Delete: %w", err)
	}
	return nil
}

func (r *tokenRepo) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.db.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	if err != nil {
		return fmt.Errorf("tokenRepo.DeleteByUserID: %w", err)
	}
	return nil
}
