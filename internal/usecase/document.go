package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
	"github.com/aidilbaihaqi/prabodrive-be/internal/services"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/constants"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/utils"
)

type PresignOutput struct {
	S3Key     string
	UploadURL string
	ExpiresAt time.Time
}

type DocumentUsecase interface {
	List(ctx context.Context, userID string, folderID *string, search string, page, limit int) ([]*domain.Document, int, error)
	Get(ctx context.Context, id, userID string) (*domain.Document, error)
	PresignUpload(ctx context.Context, userID, name string, size int64, mimeType string, folderID *string) (*PresignOutput, error)
	ConfirmUpload(ctx context.Context, userID, s3Key, name, mimeType string, size int64, folderID *string, ip string) (docID string, err error)
	Rename(ctx context.Context, id, userID, name string) error
	Delete(ctx context.Context, id, userID, ip string) (docID string, err error)
	Download(ctx context.Context, id, userID, ip string) (url string, expiresAt time.Time, err error)
}

type documentUsecase struct {
	docs     domain.DocumentRepository
	users    domain.UserRepository
	activity domain.ActivityRepository
	s3       *services.S3Service
	quota    *services.QuotaService
}

func NewDocumentUsecase(
	docs domain.DocumentRepository,
	users domain.UserRepository,
	activity domain.ActivityRepository,
	s3 *services.S3Service,
	quota *services.QuotaService,
) DocumentUsecase {
	return &documentUsecase{docs: docs, users: users, activity: activity, s3: s3, quota: quota}
}

func (u *documentUsecase) List(ctx context.Context, userID string, folderID *string, search string, page, limit int) ([]*domain.Document, int, error) {
	return u.docs.List(ctx, userID, folderID, search, page, limit)
}

func (u *documentUsecase) Get(ctx context.Context, id, userID string) (*domain.Document, error) {
	return u.docs.FindByID(ctx, id, userID)
}

func (u *documentUsecase) PresignUpload(ctx context.Context, userID, name string, size int64, mimeType string, folderID *string) (*PresignOutput, error) {
	if !utils.IsAllowedMIME(mimeType) {
		return nil, domain.ErrMIMENotAllowed
	}
	if size > constants.MaxFileSize {
		return nil, domain.ErrFileTooLarge
	}
	if err := u.quota.Check(ctx, userID, size); err != nil {
		return nil, err
	}

	folderStr := ""
	if folderID != nil {
		folderStr = *folderID
	}
	docID := uuid.New().String()
	s3Key := services.S3Key(userID, folderStr, docID, utils.SanitizeFilename(name))

	uploadURL, expiresAt, err := u.s3.GeneratePutURL(ctx, s3Key, mimeType)
	if err != nil {
		return nil, err
	}

	return &PresignOutput{S3Key: s3Key, UploadURL: uploadURL, ExpiresAt: expiresAt}, nil
}

func (u *documentUsecase) ConfirmUpload(ctx context.Context, userID, s3Key, name, mimeType string, size int64, folderID *string, ip string) (string, error) {
	if !utils.IsAllowedMIME(mimeType) {
		return "", domain.ErrMIMENotAllowed
	}
	if size > constants.MaxFileSize {
		return "", domain.ErrFileTooLarge
	}

	now := time.Now()
	docID := uuid.New().String()
	doc := &domain.Document{
		ID:        docID,
		UserID:    userID,
		FolderID:  folderID,
		Name:      name,
		Size:      size,
		MIMEType:  mimeType,
		S3Key:     s3Key,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := u.docs.Create(ctx, doc); err != nil {
		return "", err
	}
	if err := u.users.AddQuota(ctx, userID, size); err != nil {
		return "", err
	}

	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:           uuid.New().String(),
		UserID:       &userID,
		Action:       constants.ActionUpload,
		DocumentID:   &docID,
		DocumentName: &name,
		IPAddress:    &ip,
		CreatedAt:    now,
	})

	return docID, nil
}

func (u *documentUsecase) Rename(ctx context.Context, id, userID, name string) error {
	doc, err := u.docs.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}

	// Rename S3 object: copy to new key, delete old key
	oldKey := doc.S3Key
	folderStr := ""
	if doc.FolderID != nil {
		folderStr = *doc.FolderID
	}
	newKey := services.S3Key(userID, folderStr, id, utils.SanitizeFilename(name))

	if oldKey != newKey {
		if err := u.s3.CopyObject(ctx, oldKey, newKey); err != nil {
			return err
		}
		if err := u.s3.DeleteObject(ctx, oldKey); err != nil {
			return err
		}
	}

	if err := u.docs.Rename(ctx, id, userID, name, newKey); err != nil {
		return err
	}

	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:           uuid.New().String(),
		UserID:       &userID,
		Action:       constants.ActionRename,
		DocumentID:   &id,
		DocumentName: &name,
		CreatedAt:    time.Now(),
	})

	return nil
}

func (u *documentUsecase) Delete(ctx context.Context, id, userID, ip string) (string, error) {
	doc, err := u.docs.FindByID(ctx, id, userID)
	if err != nil {
		return "", err
	}

	// Log activity BEFORE actual delete so it's always recorded
	docName := doc.Name
	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:           uuid.New().String(),
		UserID:       &userID,
		Action:       constants.ActionDelete,
		DocumentID:   &id,
		DocumentName: &docName,
		IPAddress:    &ip,
		CreatedAt:    time.Now(),
	})

	if err := u.s3.DeleteObject(ctx, doc.S3Key); err != nil {
		return "", fmt.Errorf("s3 delete: %w", err)
	}

	if _, err := u.docs.Delete(ctx, id, userID); err != nil {
		return "", err
	}

	if err := u.users.AddQuota(ctx, userID, -doc.Size); err != nil {
		return "", err
	}

	return doc.ID, nil
}

func (u *documentUsecase) Download(ctx context.Context, id, userID, ip string) (string, time.Time, error) {
	doc, err := u.docs.FindByID(ctx, id, userID)
	if err != nil {
		return "", time.Time{}, err
	}

	url, expiresAt, err := u.s3.GenerateGetURL(ctx, doc.S3Key, 15*time.Minute)
	if err != nil {
		return "", time.Time{}, err
	}

	docID := doc.ID
	docName := doc.Name
	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:           uuid.New().String(),
		UserID:       &userID,
		Action:       constants.ActionDownload,
		DocumentID:   &docID,
		DocumentName: &docName,
		IPAddress:    &ip,
		CreatedAt:    time.Now(),
	})

	return url, expiresAt, nil
}

// compile-time check
var _ DocumentUsecase = (*documentUsecase)(nil)
