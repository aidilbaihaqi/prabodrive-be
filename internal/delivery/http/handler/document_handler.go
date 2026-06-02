package handler

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/request"
	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
	"github.com/aidilbaihaqi/prabodrive-be/internal/services"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/constants"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/utils"
)

type DocumentHandler struct {
	docs     domain.DocumentRepository
	users    domain.UserRepository
	activity domain.ActivityRepository
	s3       *services.S3Service
	quota    *services.QuotaService
}

func NewDocumentHandler(
	docs domain.DocumentRepository,
	users domain.UserRepository,
	activity domain.ActivityRepository,
	s3 *services.S3Service,
	quota *services.QuotaService,
) *DocumentHandler {
	return &DocumentHandler{docs: docs, users: users, activity: activity, s3: s3, quota: quota}
}

func (h *DocumentHandler) List(c *gin.Context) {
	var q request.ListDocumentsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	page, limit := clampPage(q.Page, q.Limit)

	docs, total, err := h.docs.List(c.Request.Context(), userID, q.FolderID, q.Search, page, limit)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.OKList(c, "documents fetched", toDocList(docs), page, limit, total)
}

func (h *DocumentHandler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	doc, err := h.docs.FindByID(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}
	response.OK(c, "document fetched", toDocResponse(doc))
}

func (h *DocumentHandler) PresignUpload(c *gin.Context) {
	var req request.PresignUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if !utils.IsAllowedMIME(req.MIMEType) {
		response.BadRequest(c, domain.ErrMIMENotAllowed.Error())
		return
	}
	if req.Size > constants.MaxFileSize {
		response.BadRequest(c, domain.ErrFileTooLarge.Error())
		return
	}

	userID := c.GetString("user_id")
	if err := h.quota.Check(c.Request.Context(), userID, req.Size); err != nil {
		if errors.Is(err, domain.ErrQuotaExceeded) {
			response.Forbidden(c, domain.ErrQuotaExceeded.Error())
			return
		}
		response.InternalError(c, err)
		return
	}

	folderID := ""
	if req.FolderID != nil {
		folderID = *req.FolderID
	}
	docID := uuid.New().String()
	s3Key := services.S3Key(userID, folderID, docID, utils.SanitizeFilename(req.Name))

	uploadURL, expiresAt, err := h.s3.GeneratePutURL(c.Request.Context(), s3Key, req.MIMEType)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.OK(c, "presigned URL generated", gin.H{
		"upload_url": uploadURL,
		"s3_key":     s3Key,
		"expires_at": expiresAt.Format(time.RFC3339),
	})
}

func (h *DocumentHandler) ConfirmUpload(c *gin.Context) {
	var req request.ConfirmUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if !utils.IsAllowedMIME(req.MIMEType) {
		response.BadRequest(c, domain.ErrMIMENotAllowed.Error())
		return
	}
	if req.Size > constants.MaxFileSize {
		response.BadRequest(c, domain.ErrFileTooLarge.Error())
		return
	}

	userID := c.GetString("user_id")
	now := time.Now()
	docID := uuid.New().String()

	doc := &domain.Document{
		ID:        docID,
		UserID:    userID,
		FolderID:  req.FolderID,
		Name:      req.Name,
		Size:      req.Size,
		MIMEType:  req.MIMEType,
		S3Key:     req.S3Key,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.docs.Create(c.Request.Context(), doc); err != nil {
		response.InternalError(c, err)
		return
	}

	if err := h.users.AddQuota(c.Request.Context(), userID, req.Size); err != nil {
		response.InternalError(c, err)
		return
	}

	ip := c.ClientIP()
	_ = h.activity.Log(c.Request.Context(), &domain.ActivityLog{
		ID:         uuid.New().String(),
		UserID:     &userID,
		Action:     constants.ActionUpload,
		DocumentID: &docID,
		IPAddress:  &ip,
		CreatedAt:  now,
	})

	response.Created(c, "document uploaded", gin.H{"id": docID})
}

func (h *DocumentHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	docID := c.Param("id")

	doc, err := h.docs.Delete(c.Request.Context(), docID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	if err := h.s3.DeleteObject(c.Request.Context(), doc.S3Key); err != nil {
		response.InternalError(c, fmt.Errorf("s3 delete: %w", err))
		return
	}

	if err := h.users.AddQuota(c.Request.Context(), userID, -doc.Size); err != nil {
		response.InternalError(c, err)
		return
	}

	ip := c.ClientIP()
	_ = h.activity.Log(c.Request.Context(), &domain.ActivityLog{
		ID:         uuid.New().String(),
		UserID:     &userID,
		Action:     constants.ActionDelete,
		DocumentID: &docID,
		IPAddress:  &ip,
		CreatedAt:  time.Now(),
	})

	response.Deleted(c, "document deleted", gin.H{"id": doc.ID})
}

func (h *DocumentHandler) Download(c *gin.Context) {
	userID := c.GetString("user_id")
	doc, err := h.docs.FindByID(c.Request.Context(), c.Param("id"), userID)
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
	docID := doc.ID
	_ = h.activity.Log(c.Request.Context(), &domain.ActivityLog{
		ID:         uuid.New().String(),
		UserID:     &userID,
		Action:     constants.ActionDownload,
		DocumentID: &docID,
		IPAddress:  &ip,
		CreatedAt:  time.Now(),
	})

	response.OK(c, "download URL generated", gin.H{
		"url":        downloadURL,
		"expires_at": expiresAt.Format(time.RFC3339),
	})
}

type docResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Size      int64   `json:"size"`
	MIMEType  string  `json:"mime_type"`
	FolderID  *string `json:"folder_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func toDocResponse(d *domain.Document) docResponse {
	return docResponse{
		ID:        d.ID,
		Name:      d.Name,
		Size:      d.Size,
		MIMEType:  d.MIMEType,
		FolderID:  d.FolderID,
		CreatedAt: d.CreatedAt.Format(time.RFC3339),
		UpdatedAt: d.UpdatedAt.Format(time.RFC3339),
	}
}

func toDocList(docs []*domain.Document) []docResponse {
	out := make([]docResponse, 0, len(docs))
	for _, d := range docs {
		out = append(out, toDocResponse(d))
	}
	return out
}
