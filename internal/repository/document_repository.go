package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
)

type docRepo struct {
	db *pgxpool.Pool
}

func NewDocumentRepository(db *pgxpool.Pool) domain.DocumentRepository {
	return &docRepo{db: db}
}

func (r *docRepo) Create(ctx context.Context, doc *domain.Document) error {
	_, err := r.db.Exec(ctx,
		`INSERT INTO documents (id, user_id, folder_id, name, size, mime_type, s3_key, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		doc.ID, doc.UserID, doc.FolderID, doc.Name, doc.Size, doc.MIMEType, doc.S3Key, doc.CreatedAt, doc.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("docRepo.Create: %w", err)
	}
	return nil
}

func (r *docRepo) FindByID(ctx context.Context, id, userID string) (*domain.Document, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, user_id, folder_id, name, size, mime_type, s3_key, created_at, updated_at
		 FROM documents WHERE id = $1 AND user_id = $2`, id, userID)
	doc, err := scanDoc(row)
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			return nil, domain.ErrDocumentNotFound
		}
		return nil, err
	}
	return doc, nil
}

func (r *docRepo) List(ctx context.Context, userID string, folderID *string, search string, page, limit int) ([]*domain.Document, int, error) {
	args := []any{userID}
	where := []string{"user_id = $1"}
	idx := 2

	if folderID != nil {
		where = append(where, fmt.Sprintf("folder_id = $%d", idx))
		args = append(args, *folderID)
		idx++
	}
	if search != "" {
		where = append(where, fmt.Sprintf("name ILIKE $%d", idx))
		args = append(args, "%"+search+"%")
		idx++
	}

	clause := "WHERE " + strings.Join(where, " AND ")

	var total int
	err := r.db.QueryRow(ctx, "SELECT COUNT(*) FROM documents "+clause, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("docRepo.List count: %w", err)
	}

	offset := (page - 1) * limit
	query := fmt.Sprintf(
		`SELECT id, user_id, folder_id, name, size, mime_type, s3_key, created_at, updated_at
		 FROM documents %s ORDER BY created_at DESC LIMIT $%d OFFSET $%d`,
		clause, idx, idx+1,
	)
	args = append(args, limit, offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("docRepo.List query: %w", err)
	}
	defer rows.Close()

	var docs []*domain.Document
	for rows.Next() {
		doc := &domain.Document{}
		if err := rows.Scan(&doc.ID, &doc.UserID, &doc.FolderID, &doc.Name, &doc.Size, &doc.MIMEType, &doc.S3Key, &doc.CreatedAt, &doc.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("docRepo.List scan: %w", err)
		}
		docs = append(docs, doc)
	}
	return docs, total, rows.Err()
}

func (r *docRepo) Rename(ctx context.Context, id, userID, name, s3Key string) error {
	tag, err := r.db.Exec(ctx,
		`UPDATE documents SET name = $1, s3_key = $2, updated_at = NOW() WHERE id = $3 AND user_id = $4`,
		name, s3Key, id, userID,
	)
	if err != nil {
		return fmt.Errorf("docRepo.Rename: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return domain.ErrDocumentNotFound
	}
	return nil
}

func (r *docRepo) Delete(ctx context.Context, id, userID string) (*domain.Document, error) {
	row := r.db.QueryRow(ctx,
		`DELETE FROM documents WHERE id = $1 AND user_id = $2
		 RETURNING id, user_id, folder_id, name, size, mime_type, s3_key, created_at, updated_at`,
		id, userID)
	doc, err := scanDoc(row)
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			return nil, domain.ErrDocumentNotFound
		}
		return nil, err
	}
	return doc, nil
}

func scanDoc(row pgx.Row) (*domain.Document, error) {
	d := &domain.Document{}
	err := row.Scan(&d.ID, &d.UserID, &d.FolderID, &d.Name, &d.Size, &d.MIMEType, &d.S3Key, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrDocumentNotFound
		}
		return nil, fmt.Errorf("docRepo scan: %w", err)
	}
	return d, nil
}
