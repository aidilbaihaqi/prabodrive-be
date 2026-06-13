package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/aidilbaihaqi/prabodrive-be/internal/config"
	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/constants"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/token"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/utils"
)

type AuthOutput struct {
	UserID       string
	AccessToken  string
	RefreshToken string
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type AuthUsecase interface {
	Register(ctx context.Context, email, password, name string) (*AuthOutput, error)
	Login(ctx context.Context, email, password, ip string) (*AuthOutput, error)
	Refresh(ctx context.Context, refreshToken string) (*TokenPair, error)
	Logout(ctx context.Context, userID, refreshToken, ip string) error
	GetProfile(ctx context.Context, userID string) (*domain.User, error)
	UpdateProfile(ctx context.Context, userID, name string) (*domain.User, error)
	ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error
}

type authUsecase struct {
	users    domain.UserRepository
	tokens   domain.RefreshTokenRepository
	activity domain.ActivityRepository
	cfg      config.JWTConfig
}

func NewAuthUsecase(
	users domain.UserRepository,
	tokens domain.RefreshTokenRepository,
	activity domain.ActivityRepository,
	cfg config.JWTConfig,
) AuthUsecase {
	return &authUsecase{users: users, tokens: tokens, activity: activity, cfg: cfg}
}

func (u *authUsecase) Register(ctx context.Context, email, password, name string) (*AuthOutput, error) {
	exists, err := u.users.ExistsByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, domain.ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	userID := uuid.New().String()
	user := &domain.User{
		ID:           userID,
		Email:        email,
		PasswordHash: string(hash),
		Name:         name,
		Role:         constants.RoleUser,
		QuotaMax:     constants.DefaultQuotaMax,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := u.users.Create(ctx, user); err != nil {
		return nil, err
	}

	access, refresh, err := u.issueTokenPair(ctx, userID, email, constants.RoleUser)
	if err != nil {
		return nil, err
	}

	return &AuthOutput{UserID: userID, AccessToken: access, RefreshToken: refresh}, nil
}

func (u *authUsecase) Login(ctx context.Context, email, password, ip string) (*AuthOutput, error) {
	user, err := u.users.FindByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrUnauthorized
		}
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, domain.ErrUnauthorized
	}

	access, refresh, err := u.issueTokenPair(ctx, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:        uuid.New().String(),
		UserID:    &user.ID,
		Action:    constants.ActionLogin,
		IPAddress: &ip,
		CreatedAt: time.Now(),
	})

	return &AuthOutput{UserID: user.ID, AccessToken: access, RefreshToken: refresh}, nil
}

func (u *authUsecase) Refresh(ctx context.Context, refreshToken string) (*TokenPair, error) {
	userIDFromJWT, err := token.ValidateRefresh(refreshToken, u.cfg.Secret)
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	oldHash := utils.SHA256(refreshToken)
	userID, err := u.tokens.Find(ctx, oldHash)
	if err != nil || userID != userIDFromJWT {
		return nil, domain.ErrInvalidToken
	}

	user, err := u.users.FindByID(ctx, userID)
	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	if err := u.tokens.Delete(ctx, oldHash); err != nil {
		return nil, err
	}

	access, refresh, err := u.issueTokenPair(ctx, user.ID, user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &TokenPair{AccessToken: access, RefreshToken: refresh}, nil
}

func (u *authUsecase) Logout(ctx context.Context, userID, refreshToken, ip string) error {
	_ = u.tokens.Delete(ctx, utils.SHA256(refreshToken))

	_ = u.activity.Log(ctx, &domain.ActivityLog{
		ID:        uuid.New().String(),
		UserID:    &userID,
		Action:    constants.ActionLogout,
		IPAddress: &ip,
		CreatedAt: time.Now(),
	})
	return nil
}

func (u *authUsecase) GetProfile(ctx context.Context, userID string) (*domain.User, error) {
	return u.users.FindByID(ctx, userID)
}

func (u *authUsecase) UpdateProfile(ctx context.Context, userID, name string) (*domain.User, error) {
	if err := u.users.UpdateProfile(ctx, userID, name); err != nil {
		return nil, err
	}
	return u.users.FindByID(ctx, userID)
}

func (u *authUsecase) ChangePassword(ctx context.Context, userID, currentPassword, newPassword string) error {
	user, err := u.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPassword)); err != nil {
		return domain.ErrUnauthorized
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}
	return u.users.UpdatePassword(ctx, userID, string(hash))
}

func (u *authUsecase) issueTokenPair(ctx context.Context, userID, email, role string) (access, refresh string, err error) {
	access, err = token.GenerateAccess(userID, email, role, u.cfg.Secret, u.cfg.AccessExpiry)
	if err != nil {
		return
	}
	refresh, err = token.GenerateRefresh(userID, u.cfg.Secret, u.cfg.RefreshExpiry)
	if err != nil {
		return
	}
	expiresAt := time.Now().Add(u.cfg.RefreshExpiry)
	err = u.tokens.Save(ctx, userID, utils.SHA256(refresh), expiresAt)
	return
}
