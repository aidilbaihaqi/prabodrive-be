package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
	"github.com/aidilbaihaqi/prabodrive-be/internal/services"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/constants"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/utils"
)

type ShareOutput struct {
	ID        string
	Token     string
	ShareURL  string
	ExpiresAt time.Time
}

type ShareListItem struct {
	ID          string
	DocumentID  string
	Token       string
	ShareURL    string
	HasPassword bool
	ExpiresAt   time.Time
	CreatedAt   time.Time
}

type ShareUsecase interface {
	Create(ctx context.Context, userID, documentID string, expiresAt time.Time, password *string, ip, baseURL string) (*ShareOutput, error)
	ListByUser(ctx context.Context, userID string, page, limit int, baseURL string) ([]*ShareListItem, int, error)
	Access(ctx context.Context, token, password, ip string) (downloadURL string, expiresAt time.Time, err error)
	Delete(ctx context.Context, id, userID string) error
}

type shareUsecase struct {
	shares   domain.ShareRepository
	docs     domain.DocumentRepository
	users    domain.UserRepository
	activity domain.ActivityRepository
	s3       *services.S3Service
	email    *services.EmailService
}

func NewShareUsecase(
	shares domain.ShareRepository,
	docs domain.DocumentRepository,
	users domain.UserRepository,
	activity domain.ActivityRepository,
	s3 *services.S3Service,
	email *services.EmailService,
) ShareUsecase {
	return &shareUsecase{
		shares:   shares,
		docs:     docs,
		users:    users,
		activity: activity,
		s3:       s3,
		email:    email,
	}
}

func (u *shareUsecase) Create(ctx context.Context, userID, documentID string, expiresAt time.Time, password *string, ip, baseURL string) (*ShareOutput, error) {
	doc, err := u.docs.FindByID(ctx, documentID, userID)
	if err != nil {
		return nil, err
	}

	tok, err := utils.GenerateToken(32)
	if err != nil {
		return nil, err
	}

	link := &domain.ShareLink{
		ID:         uuid.New().String(),
		DocumentID: documentID,
		Token:      tok,
		ExpiresAt:  expiresAt,
		CreatedBy:  userID,
		CreatedAt:  time.Now(),
	}

	if password != nil && *password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*password), 12)
		if err != nil {
			return nil, err
		}
		h := string(hash)
		link.PasswordHash = &h
	}

	if err := u.shares.Create(ctx, link); err != nil {
		return nil, err
	}

	shareURL := fmt.Sprintf("%s/api/v1/share/%s", baseURL, tok)

	if user, _ := u.users.FindByID(ctx, userID); user != nil {
		docName, ownerEmail := doc.Name, user.Email
		go func() {
			_ = u.email.SendShareNotification(context.Background(), ownerEmail, docName, shareURL)
		}()
	}

	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:         uuid.New().String(),
		UserID:     &userID,
		Action:     constants.ActionShareCreate,
		DocumentID: &documentID,
		IPAddress:  &ip,
		CreatedAt:  time.Now(),
	})

	return &ShareOutput{
		ID:        link.ID,
		Token:     tok,
		ShareURL:  shareURL,
		ExpiresAt: link.ExpiresAt,
	}, nil
}

func (u *shareUsecase) ListByUser(ctx context.Context, userID string, page, limit int, baseURL string) ([]*ShareListItem, int, error) {
	links, total, err := u.shares.ListByUser(ctx, userID, page, limit)
	if err != nil {
		return nil, 0, err
	}

	out := make([]*ShareListItem, 0, len(links))
	for _, l := range links {
		out = append(out, &ShareListItem{
			ID:          l.ID,
			DocumentID:  l.DocumentID,
			Token:       l.Token,
			ShareURL:    fmt.Sprintf("%s/api/v1/share/%s", baseURL, l.Token),
			HasPassword: l.PasswordHash != nil,
			ExpiresAt:   l.ExpiresAt,
			CreatedAt:   l.CreatedAt,
		})
	}

	return out, total, nil
}

func (u *shareUsecase) Access(ctx context.Context, token, password, ip string) (string, time.Time, error) {
	link, err := u.shares.FindByToken(ctx, token)
	if err != nil {
		if errors.Is(err, domain.ErrShareNotFound) {
			return "", time.Time{}, domain.ErrShareNotFound
		}
		return "", time.Time{}, err
	}

	if time.Now().After(link.ExpiresAt) {
		return "", time.Time{}, domain.ErrShareExpired
	}

	if link.PasswordHash != nil {
		if password == "" {
			return "", time.Time{}, fmt.Errorf("password required")
		}
		if err := bcrypt.CompareHashAndPassword([]byte(*link.PasswordHash), []byte(password)); err != nil {
			return "", time.Time{}, domain.ErrSharePasswordWrong
		}
	}

	doc, err := u.docs.FindByID(ctx, link.DocumentID, link.CreatedBy)
	if err != nil {
		return "", time.Time{}, err
	}

	downloadURL, expiresAt, err := u.s3.GenerateGetURL(ctx, doc.S3Key, 15*time.Minute)
	if err != nil {
		return "", time.Time{}, err
	}

	docName := doc.Name
	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:           uuid.New().String(),
		Action:       constants.ActionShareAccess,
		DocumentID:   &link.DocumentID,
		DocumentName: &docName,
		IPAddress:    &ip,
		CreatedAt:    time.Now(),
	})

	return downloadURL, expiresAt, nil
}

func (u *shareUsecase) Delete(ctx context.Context, id, userID string) error {
	share, err := u.shares.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if err := u.shares.Delete(ctx, id, userID); err != nil {
		return err
	}

	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:         uuid.New().String(),
		UserID:     &userID,
		Action:     constants.ActionShareDelete,
		DocumentID: &share.DocumentID,
		CreatedAt:  time.Now(),
	})

	return nil
}
