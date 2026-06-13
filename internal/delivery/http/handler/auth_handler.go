package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/request"
	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
	"github.com/aidilbaihaqi/prabodrive-be/internal/usecase"
)

type AuthHandler struct {
	auth usecase.AuthUsecase
}

func NewAuthHandler(auth usecase.AuthUsecase) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	out, err := h.auth.Register(c.Request.Context(), req.Email, req.Password, req.Name)
	if err != nil {
		if errors.Is(err, domain.ErrEmailExists) {
			response.BadRequest(c, err.Error())
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Created(c, "register successful", gin.H{
		"id":            out.UserID,
		"access_token":  out.AccessToken,
		"refresh_token": out.RefreshToken,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	out, err := h.auth.Login(c.Request.Context(), req.Email, req.Password, c.ClientIP())
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			response.Unauthorized(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.OK(c, "login successful", gin.H{
		"id":            out.UserID,
		"access_token":  out.AccessToken,
		"refresh_token": out.RefreshToken,
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	var req request.RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	pair, err := h.auth.Refresh(c.Request.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidToken) || errors.Is(err, domain.ErrUnauthorized) {
			response.Unauthorized(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.OK(c, "token refreshed", gin.H{
		"access_token":  pair.AccessToken,
		"refresh_token": pair.RefreshToken,
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	user, err := h.auth.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, "profile retrieved", gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"quota_used": user.QuotaUsed,
		"quota_max":  user.QuotaMax,
		"created_at": user.CreatedAt,
		"updated_at": user.UpdatedAt,
	})
}

func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	var req request.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")

	if req.Name == "" && req.NewPassword == "" {
		response.BadRequest(c, "at least one of name or new_password is required")
		return
	}

	if req.NewPassword != "" {
		if req.CurrentPassword == "" {
			response.BadRequest(c, "current_password is required when changing password")
			return
		}
	}

	if req.NewPassword != "" {
		if err := h.auth.ChangePassword(c.Request.Context(), userID, req.CurrentPassword, req.NewPassword); err != nil {
			if errors.Is(err, domain.ErrUnauthorized) {
				response.BadRequest(c, "current password is incorrect")
				return
			}
			response.InternalError(c, err)
			return
		}
	}

	if req.Name != "" {
		if _, err := h.auth.UpdateProfile(c.Request.Context(), userID, req.Name); err != nil {
			response.InternalError(c, err)
			return
		}
	}

	user, err := h.auth.GetProfile(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.OK(c, "profile updated", gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"role":       user.Role,
		"quota_used": user.QuotaUsed,
		"quota_max":  user.QuotaMax,
		"updated_at": user.UpdatedAt,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	var req request.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	_ = h.auth.Logout(c.Request.Context(), userID, req.RefreshToken, c.ClientIP())

	response.OK(c, "logout successful", gin.H{"id": userID})
}
