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

type ShareHandler struct {
	shareUC usecase.ShareUsecase
	baseURL string
}

func NewShareHandler(shareUC usecase.ShareUsecase, baseURL string) *ShareHandler {
	return &ShareHandler{shareUC: shareUC, baseURL: baseURL}
}

func (h *ShareHandler) Create(c *gin.Context) {
	var req request.CreateShareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	out, err := h.shareUC.Create(c.Request.Context(), userID, req.DocumentID, req.ExpiresAt, req.Password, c.ClientIP(), h.baseURL)
	if err != nil {
		if errors.Is(err, domain.ErrDocumentNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Created(c, "share link created", gin.H{
		"id":         out.ID,
		"token":      out.Token,
		"share_url":  out.ShareURL,
		"expires_at": out.ExpiresAt.Format(time.RFC3339),
	})
}

func (h *ShareHandler) List(c *gin.Context) {
	var q request.PaginationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	page, limit := clampPage(q.Page, q.Limit)

	items, total, err := h.shareUC.ListByUser(c.Request.Context(), userID, page, limit, h.baseURL)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.OKList(c, "share links fetched", toShareList(items), page, limit, total)
}

func (h *ShareHandler) Access(c *gin.Context) {
	var q request.AccessShareQuery
	_ = c.ShouldBindQuery(&q)

	tok := c.Param("token")
	downloadURL, expiresAt, err := h.shareUC.Access(c.Request.Context(), tok, q.Password, c.ClientIP())
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrShareNotFound):
			response.NotFound(c)
		case errors.Is(err, domain.ErrShareExpired):
			response.Forbidden(c, err.Error())
		case errors.Is(err, domain.ErrSharePasswordWrong):
			response.Forbidden(c, err.Error())
		case err.Error() == "password required":
			response.Forbidden(c, "password required")
		default:
			response.InternalError(c, err)
		}
		return
	}

	response.OK(c, "share accessed", gin.H{
		"download_url": downloadURL,
		"expires_at":   expiresAt.Format(time.RFC3339),
	})
}

func (h *ShareHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	linkID := c.Param("id")

	if err := h.shareUC.Delete(c.Request.Context(), linkID, userID); err != nil {
		if errors.Is(err, domain.ErrShareNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Deleted(c, "share link deleted", gin.H{"id": linkID})
}

type shareListResponse struct {
	ID          string `json:"id"`
	DocumentID  string `json:"document_id"`
	Token       string `json:"token"`
	ShareURL    string `json:"share_url"`
	HasPassword bool   `json:"has_password"`
	ExpiresAt   string `json:"expires_at"`
	CreatedAt   string `json:"created_at"`
}

func toShareList(items []*usecase.ShareListItem) []shareListResponse {
	out := make([]shareListResponse, 0, len(items))
	for _, item := range items {
		out = append(out, shareListResponse{
			ID:          item.ID,
			DocumentID:  item.DocumentID,
			Token:       item.Token,
			ShareURL:    item.ShareURL,
			HasPassword: item.HasPassword,
			ExpiresAt:   item.ExpiresAt.Format(time.RFC3339),
			CreatedAt:   item.CreatedAt.Format(time.RFC3339),
		})
	}
	return out
}
