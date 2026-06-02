package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/aidilbaihaqi/prabodrive-be/internal/config"
	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/request"
	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/constants"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/token"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/utils"
)

type AuthHandler struct {
	users    domain.UserRepository
	tokens   domain.RefreshTokenRepository
	activity domain.ActivityRepository
	cfg      config.JWTConfig
}

func NewAuthHandler(
	users domain.UserRepository,
	tokens domain.RefreshTokenRepository,
	activity domain.ActivityRepository,
	cfg config.JWTConfig,
) *AuthHandler {
	return &AuthHandler{users: users, tokens: tokens, activity: activity, cfg: cfg}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	exists, err := h.users.ExistsByEmail(c.Request.Context(), req.Email)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	if exists {
		response.BadRequest(c, domain.ErrEmailExists.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	now := time.Now()
	userID := uuid.New().String()
	user := &domain.User{
		ID:           userID,
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
		Role:         constants.RoleUser,
		QuotaUsed:    0,
		QuotaMax:     constants.DefaultQuotaMax,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := h.users.Create(c.Request.Context(), user); err != nil {
		response.InternalError(c, err)
		return
	}

	access, refresh, err := h.issueTokenPair(c, userID, req.Email, constants.RoleUser)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, "register successful", gin.H{
		"id":            userID,
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	user, err := h.users.FindByEmail(c.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.Unauthorized(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		response.Unauthorized(c)
		return
	}

	access, refresh, err := h.issueTokenPair(c, user.ID, user.Email, user.Role)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	ip := c.ClientIP()
	_ = h.activity.Log(c.Request.Context(), &domain.ActivityLog{
		ID:        uuid.New().String(),
		UserID:    &user.ID,
		Action:    constants.ActionLogin,
		IPAddress: &ip,
		CreatedAt: time.Now(),
	})

	response.OK(c, "login successful", gin.H{
		"id":            user.ID,
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req request.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userIDFromJWT, err := token.ValidateRefresh(req.RefreshToken, h.cfg.Secret)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	oldHash := utils.SHA256(req.RefreshToken)
	userID, err := h.tokens.Find(c.Request.Context(), oldHash)
	if err != nil || userID != userIDFromJWT {
		response.Unauthorized(c)
		return
	}

	user, err := h.users.FindByID(c.Request.Context(), userID)
	if err != nil {
		response.Unauthorized(c)
		return
	}

	if err := h.tokens.Delete(c.Request.Context(), oldHash); err != nil {
		response.InternalError(c, err)
		return
	}

	access, refresh, err := h.issueTokenPair(c, user.ID, user.Email, user.Role)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.OK(c, "token refreshed", gin.H{
		"access_token":  access,
		"refresh_token": refresh,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req request.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	tokenHash := utils.SHA256(req.RefreshToken)
	_ = h.tokens.Delete(c.Request.Context(), tokenHash)

	ip := c.ClientIP()
	_ = h.activity.Log(c.Request.Context(), &domain.ActivityLog{
		ID:        uuid.New().String(),
		UserID:    &userID,
		Action:    constants.ActionLogout,
		IPAddress: &ip,
		CreatedAt: time.Now(),
	})

	response.OK(c, "logout successful", gin.H{"id": userID})
}

func (h *AuthHandler) issueTokenPair(c *gin.Context, userID, email, role string) (access, refresh string, err error) {
	access, err = token.GenerateAccess(userID, email, role, h.cfg.Secret, h.cfg.AccessExpiry)
	if err != nil {
		return
	}

	refresh, err = token.GenerateRefresh(userID, h.cfg.Secret, h.cfg.RefreshExpiry)
	if err != nil {
		return
	}

	expiresAt := time.Now().Add(h.cfg.RefreshExpiry)
	err = h.tokens.Save(c.Request.Context(), userID, utils.SHA256(refresh), expiresAt)
	return
}
