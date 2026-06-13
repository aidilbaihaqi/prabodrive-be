package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/request"
	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
	"github.com/aidilbaihaqi/prabodrive-be/internal/usecase"
)

type DocumentHandler struct {
	docUC usecase.DocumentUsecase
}

func NewDocumentHandler(docUC usecase.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{docUC: docUC}
}

func (h *DocumentHandler) List(c *gin.Context) {
	var q request.ListDocumentsQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	page, limit := clampPage(q.Page, q.Limit)

	docs, total, err := h.docUC.List(c.Request.Context(), userID, q.FolderID, q.Search, page, limit)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.OKList(c, "documents fetched", toDocList(docs), page, limit, total)
}

func (h *DocumentHandler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	doc, err := h.docUC.Get(c.Request.Context(), c.Param("id"), userID)
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

	userID := c.GetString("user_id")
	out, err := h.docUC.PresignUpload(c.Request.Context(), userID, req.Name, req.Size, req.MIMEType, req.FolderID)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMIMENotAllowed):
			response.BadRequest(c, err.Error())
		case errors.Is(err, domain.ErrFileTooLarge):
			response.BadRequest(c, err.Error())
		case errors.Is(err, domain.ErrQuotaExceeded):
			response.Forbidden(c, err.Error())
		default:
			response.InternalError(c, err)
		}
		return
	}

	response.OK(c, "presigned URL generated", gin.H{
		"upload_url": out.UploadURL,
		"s3_key":     out.S3Key,
		"expires_at": out.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *DocumentHandler) ConfirmUpload(c *gin.Context) {
	var req request.ConfirmUploadRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	docID, err := h.docUC.ConfirmUpload(c.Request.Context(), userID, req.S3Key, req.Name, req.MIMEType, req.Size, req.FolderID, c.ClientIP())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrMIMENotAllowed):
			response.BadRequest(c, err.Error())
		case errors.Is(err, domain.ErrFileTooLarge):
			response.BadRequest(c, err.Error())
		default:
			response.InternalError(c, err)
		}
		return
	}

	response.Created(c, "document uploaded", gin.H{"id": docID})
}

func (h *DocumentHandler) Rename(c *gin.Context) {
	var req request.RenameDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	docID := c.Param("id")

	if err := h.docUC.Rename(c.Request.Context(), docID, userID, req.Name); err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Updated(c, "document renamed", gin.H{"id": docID})
}

func (h *DocumentHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	docID, err := h.docUC.Delete(c.Request.Context(), c.Param("id"), userID, c.ClientIP())
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Deleted(c, "document deleted", gin.H{"id": docID})
}

func (h *DocumentHandler) Download(c *gin.Context) {
	userID := c.GetString("user_id")
	url, expiresAt, err := h.docUC.Download(c.Request.Context(), c.Param("id"), userID, c.ClientIP())
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.OK(c, "download URL generated", gin.H{
		"url":        url,
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
