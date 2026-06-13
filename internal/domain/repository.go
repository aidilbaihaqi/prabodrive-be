package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindByID(ctx context.Context, id string) (*User, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	AddQuota(ctx context.Context, userID string, delta int64) error

	// Admin operations
	ListAll(ctx context.Context, page, limit int) ([]*User, int, error)
	UpdateRole(ctx context.Context, id, role string) error
	DeleteUser(ctx context.Context, id string) error
}

type RefreshTokenRepository interface {
	Save(ctx context.Context, userID, tokenHash string, expiresAt interface{}) error
	Find(ctx context.Context, tokenHash string) (userID string, err error)
	Delete(ctx context.Context, tokenHash string) error
	DeleteByUserID(ctx context.Context, userID string) error
}

type DocumentRepository interface {
	Create(ctx context.Context, doc *Document) error
	FindByID(ctx context.Context, id, userID string) (*Document, error)
	List(ctx context.Context, userID string, folderID *string, search string, page, limit int) ([]*Document, int, error)
	Rename(ctx context.Context, id, userID, name string) error
	Delete(ctx context.Context, id, userID string) (*Document, error)
}

type FolderRepository interface {
	Create(ctx context.Context, folder *Folder) error
	FindByID(ctx context.Context, id, userID string) (*Folder, error)
	// parentID nil = all folders; ptr to "" = root only; ptr to uuid = children of that folder
	List(ctx context.Context, userID string, parentID *string) ([]*Folder, error)
	Update(ctx context.Context, id, userID, name string) error
	Delete(ctx context.Context, id, userID string) error
}

type ShareRepository interface {
	Create(ctx context.Context, link *ShareLink) error
	FindByToken(ctx context.Context, token string) (*ShareLink, error)
	FindByID(ctx context.Context, id string) (*ShareLink, error)
	ListByUser(ctx context.Context, userID string, page, limit int) ([]*ShareLink, int, error)
	Delete(ctx context.Context, id, createdBy string) error
}

type ActivityRepository interface {
	Log(ctx context.Context, entry *ActivityLog) error
	List(ctx context.Context, userID string, page, limit int) ([]*ActivityLog, int, error)
}
