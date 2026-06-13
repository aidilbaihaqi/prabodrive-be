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

type folderRepo struct {
	db *pgxpool.Pool
}

func NewFolderRepository(db *pgxpool.Pool) domain.FolderRepository {
	return &folderRepo{db: db}
}

func (r *folderRepo) Create(ctx context.Context, f *domain.Folder) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO folders (id, user_id, parent_id, name, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		f.ID, f.UserID, f.ParentID, f.Name, f.CreatedAt, f.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("folderRepo.Create: %w", err)
	}
	return nil
}

func (r *folderRepo) FindByID(ctx context.Context, id, userID string) (*domain.Folder, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, user_id, parent_id, name, created_at, updated_at
		 FROM folders WHERE id = $1 AND user_id = $2`, id, userID)
	return scanFolder(row)
}

func (r *folderRepo) List(ctx context.Context, userID string, parentID *string) ([]*domain.Folder, error) {
	var query string
	var args []any

	const cols = `SELECT id, user_id, parent_id, name, created_at, updated_at FROM folders`
	switch {
	case parentID == nil:
		query = cols + ` WHERE user_id = $1 ORDER BY name`
		args = []any{userID}
	case *parentID == "":
		query = cols + ` WHERE user_id = $1 AND parent_id IS NULL ORDER BY name`
		args = []any{userID}
	default:
		query = cols + ` WHERE user_id = $1 AND parent_id = $2 ORDER BY name`
		args = []any{userID, *parentID}
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("folderRepo.List: %w", err)
	}
	defer rows.Close()

	var folders []*domain.Folder
	for rows.Next() {
		f := &domain.Folder{}
		if err := rows.Scan(&f.ID, &f.UserID, &f.ParentID, &f.Name, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, fmt.Errorf("folderRepo.List scan: %w", err)
		}
		folders = append(folders, f)
	}
	return folders, rows.Err()
}

func (r *folderRepo) Update(ctx context.Context, id, userID, name string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE folders SET name = $1, updated_at = $2 WHERE id = $3 AND user_id = $4`,
		name, time.Now(), id, userID,
	)
	if err != nil {
		return fmt.Errorf("folderRepo.Update: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrFolderNotFound
	}
	return nil
}

func (r *folderRepo) Delete(ctx context.Context, id, userID string) error {
	// Move documents inside this folder to root (folder_id = NULL)
	_, err := r.db.Exec(ctx,
		`UPDATE documents SET folder_id = NULL WHERE folder_id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("folderRepo.Delete move docs: %w", err)
	}

	tag, err := r.db.Exec(ctx,
		`DELETE FROM folders WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("folderRepo.Delete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrFolderNotFound
	}
	return nil
}

func scanFolder(row pgx.Row) (*domain.Folder, error) {
	f := &domain.Folder{}
	err := row.Scan(&f.ID, &f.UserID, &f.ParentID, &f.Name, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrFolderNotFound
		}
		return nil, fmt.Errorf("folderRepo scan: %w", err)
	}
	return f, nil
}
