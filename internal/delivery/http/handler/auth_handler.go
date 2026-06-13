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
