package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/constants"
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
	folders  domain.FolderRepository
	activity domain.ActivityRepository
}

func NewFolderUsecase(folders domain.FolderRepository, activity domain.ActivityRepository) FolderUsecase {
	return &folderUsecase{folders: folders, activity: activity}
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

	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:           uuid.New().String(),
		UserID:       &userID,
		Action:       constants.ActionCreateFolder,
		DocumentID:   &folder.ID,
		DocumentName: &name,
		CreatedAt:    now,
	})

	return folder, nil
}

func (u *folderUsecase) Update(ctx context.Context, id, userID, name string) error {
	folder, err := u.folders.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}

	if err := u.folders.Update(ctx, id, userID, name); err != nil {
		return err
	}

	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:           uuid.New().String(),
		UserID:       &userID,
		Action:       constants.ActionRenameFolder,
		DocumentID:   &id,
		DocumentName: &folder.Name,
		CreatedAt:    time.Now(),
	})

	return nil
}

func (u *folderUsecase) Delete(ctx context.Context, id, userID string) error {
	folder, err := u.folders.FindByID(ctx, id, userID)
	if err != nil {
		return err
	}

	if err := u.folders.Delete(ctx, id, userID); err != nil {
		return err
	}

	folderName := folder.Name
	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:           uuid.New().String(),
		UserID:       &userID,
		Action:       constants.ActionDeleteFolder,
		DocumentID:   &id,
		DocumentName: &folderName,
		CreatedAt:    time.Now(),
	})

	return nil
}
