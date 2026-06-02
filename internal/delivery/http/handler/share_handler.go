package handler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/request"
	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
	"github.com/aidilbaihaqi/prabodrive-be/internal/services"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/constants"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/utils"
)

type ShareHandler struct {
	shares   domain.ShareRepository
	docs     domain.DocumentRepository
	users    domain.UserRepository
	activity domain.ActivityRepository
	s3       *services.S3Service
	email    *services.EmailService
	baseURL  string
}

func NewShareHandler(
	shares domain.ShareRepository,
	docs domain.DocumentRepository,
	users domain.UserRepository,
	activity domain.ActivityRepository,
	s3 *services.S3Service,
	email *services.EmailService,
	baseURL string,
) *ShareHandler {
	return &ShareHandler{
		shares:   shares,
		docs:     docs,
		users:    users,
		activity: activity,
		s3:       s3,
		email:    email,
		baseURL:  baseURL,
	}
}

func (h *ShareHandler) Create(c *gin.Context) {
	var req request.CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")

	doc, err := h.docs.FindByID(c.Request.Context(), req.DocumentID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	tok, err := utils.GenerateToken(32)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	link := &domain.ShareLink{
		ID:         uuid.New().String(),
		DocumentID: req.DocumentID,
		Token:      tok,
		ExpiresAt:  req.ExpiresAt,
		CreatedBy:  userID,
		CreatedAt:  time.Now(),
	}

	if req.Password != nil && *req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), 12)
		if err != nil {
			response.InternalError(c, err)
			return
		}
		h := string(hash)
		link.PasswordHash = &h
	}

	if err := h.shares.Create(c.Request.Context(), link); err != nil {
		response.InternalError(c, err)
		return
	}

	shareURL := fmt.Sprintf("%s/api/v1/share/%s", h.baseURL, tok)

	if user, _ := h.users.FindByID(c.Request.Context(), userID); user != nil {
		docName, ownerEmail := doc.Name, user.Email
		go func() {
			_ = h.email.SendShareNotification(context.Background(), ownerEmail, docName, shareURL)
		}()
	}

	ip := c.ClientIP()
	linkID := link.ID
	_ = h.activity.Log(c.Request.Context(), &domain.ActivityLog{
		ID:         uuid.New().String(),
		UserID:     &userID,
		Action:     constants.ActionShareCreate,
		DocumentID: &req.DocumentID,
		IPAddress:  &ip,
		CreatedAt:  time.Now(),
	})

	response.Created(c, "share link created", gin.H{
		"id":         linkID,
		"token":      tok,
		"share_url":  shareURL,
		"expires_at": link.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *ShareHandler) Access(c *gin.Context) {
	var q request.AccessShareQuery
	_ = c.ShouldBindQuery(&q)

	tok := c.Param("token")
	link, err := h.shares.FindByToken(c.Request.Context(), tok)
	if err != nil {
		if errors.Is(err, domain.ErrShareNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	if time.Now().After(link.ExpiresAt) {
		response.Forbidden(c, domain.ErrShareExpired.Error())
		return
	}

	if link.PasswordHash != nil {
		if q.Password == "" {
			response.Forbidden(c, "password required")
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(*link.PasswordHash), []byte(q.Password)); err != nil {
			response.Forbidden(c, domain.ErrSharePasswordWrong.Error())
			return
		}
	}

	doc, err := h.docs.FindByID(c.Request.Context(), link.DocumentID, link.CreatedBy)
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	downloadURL, expiresAt, err := h.s3.GenerateGetURL(c.Request.Context(), doc.S3Key, 15*time.Minute)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	ip := c.ClientIP()
	_ = h.activity.Log(c.Request.Context(), &domain.ActivityLog{
		ID:         uuid.New().String(),
		Action:     constants.ActionShareAccess,
		DocumentID: &link.DocumentID,
		IPAddress:  &ip,
		CreatedAt:  time.Now(),
	})

	response.OK(c, "share accessed", gin.H{
		"download_url": downloadURL,
		"expires_at":   expiresAt.Format(time.RFC3339),
	})
}

func (h *ShareHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	linkID := c.Param("id")

	if err := h.shares.Delete(c.Request.Context(), linkID, userID); err != nil {
		if errors.Is(err, domain.ErrShareNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Deleted(c, "share link deleted", gin.H{"id": linkID})
}
