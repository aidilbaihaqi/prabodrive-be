package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
)

type FolderUsecase interface {
	// parentID nil = all; ptr to "" = root; ptr to uuid = children
	List(ctx context.Context, userID string, parentID *string) ([]*domain.Folder, error)
	Get(ctx context.Context, id, userID string) (*domain.Folder, error)
	Create(ctx context.Context, userID, name string, parentID *string) (*domain.Folder, error)
	Update(ctx context.Context, id, userID, name string) error
	Delete(ctx context.Context, id, userID string) error
}

type folderUsecase struct {
	folders domain.FolderRepository
}

func NewFolderUsecase(folders domain.FolderRepository) FolderUsecase {
	return &folderUsecase{folders: folders}
}

func (u *folderUsecase) List(ctx context.Context, userID string, parentID *string) ([]*domain.Folder, error) {
	return u.folders.List(ctx, userID, parentID)
}

func (u *folderUsecase) Get(ctx context.Context, id, userID string) (*domain.Folder, error) {
	return u.folders.FindByID(ctx, id, userID)
}

func (u *folderUsecase) Create(ctx context.Context, userID, name string, parentID *string) (*domain.Folder, error) {
	now := time.Now()
	folder := &domain.Folder{
		ID:        uuid.New().String(),
		UserID:    userID,
		ParentID:  parentID,
		Name:      name,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := u.folders.Create(ctx, folder); err != nil {
		return nil, err
	}
	return folder, nil
}

func (u *folderUsecase) Update(ctx context.Context, id, userID, name string) error {
	return u.folders.Update(ctx, id, userID, name)
}

func (u *folderUsecase) Delete(ctx context.Context, id, userID string) error {
	return u.folders.Delete(ctx, id, userID)
}
