package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
)

type shareRepo struct {
	db *pgxpool.Pool
}

func NewShareRepository(db *pgxpool.Pool) domain.ShareRepository {
	return &shareRepo{db: db}
}

func (r *shareRepo) Create(ctx context.Context, link *domain.ShareLink) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO share_links (id, document_id, token, password_hash, expires_at, created_by, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		link.ID, link.DocumentID, link.Token, link.PasswordHash,
		link.ExpiresAt, link.CreatedBy, link.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("shareRepo.Create: %w", err)
	}
	return nil
}

func (r *shareRepo) FindByToken(ctx context.Context, token string) (*domain.ShareLink, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, document_id, token, password_hash, expires_at, created_by, created_at
		 FROM share_links WHERE token = $1`, token)
	return scanShare(row)
}

func (r *shareRepo) FindByID(ctx context.Context, id string) (*domain.ShareLink, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, document_id, token, password_hash, expires_at, created_by, created_at
		 FROM share_links WHERE id = $1`, id)
	return scanShare(row)
}

func (r *shareRepo) Delete(ctx context.Context, id, createdBy string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM share_links WHERE id = $1 AND created_by = $2`, id, createdBy)
	if err != nil {
		return fmt.Errorf("shareRepo.Delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrShareNotFound
	}
	return nil
}

func scanShare(row pgx.Row) (*domain.ShareLink, error) {
	s := &domain.ShareLink{}
	err := row.Scan(&s.ID, &s.DocumentID, &s.Token, &s.PasswordHash, &s.ExpiresAt, &s.CreatedBy, &s.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrShareNotFound
		}
		return nil, fmt.Errorf("shareRepo scan: %w", err)
	}
	return s, nil
}
