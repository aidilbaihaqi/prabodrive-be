package handler

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/request"
	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
)

type FolderHandler struct {
	folders domain.FolderRepository
}

func NewFolderHandler(folders domain.FolderRepository) *FolderHandler {
	return &FolderHandler{folders: folders}
}

func (h *FolderHandler) List(c *gin.Context) {
	userID := c.GetString("user_id")
	folders, err := h.folders.List(c.Request.Context(), userID)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, "folders fetched", toFolderList(folders))
}

func (h *FolderHandler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	folder, err := h.folders.FindByID(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		if errors.Is(err, domain.ErrFolderNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}
	response.OK(c, "folder fetched", toFolderResponse(folder))
}

func (h *FolderHandler) Create(c *gin.Context) {
	var req request.CreateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	now := time.Now()
	folder := &domain.Folder{
		ID:        uuid.New().String(),
		UserID:    userID,
		ParentID:  req.ParentID,
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := h.folders.Create(c.Request.Context(), folder); err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, "folder created", gin.H{"id": folder.ID})
}

func (h *FolderHandler) Update(c *gin.Context) {
	var req request.UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	folderID := c.Param("id")

	if err := h.folders.Update(c.Request.Context(), folderID, userID, req.Name); err != nil {
		if errors.Is(err, domain.ErrFolderNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Updated(c, "folder updated", gin.H{"id": folderID})
}

func (h *FolderHandler) Delete(c *gin.Context) {
	userID := c.GetString("user_id")
	folderID := c.Param("id")

	if err := h.folders.Delete(c.Request.Context(), folderID, userID); err != nil {
		if errors.Is(err, domain.ErrFolderNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}

	response.Deleted(c, "folder deleted", gin.H{"id": folderID})
}

type folderResponse struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	ParentID  *string `json:"parent_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

func toFolderResponse(f *domain.Folder) folderResponse {
	return folderResponse{
		ID:        f.ID,
		Name:      f.Name,
		ParentID:  f.ParentID,
		CreatedAt: f.CreatedAt.Format(time.RFC3339),
		UpdatedAt: f.UpdatedAt.Format(time.RFC3339),
	}
}

func toFolderList(folders []*domain.Folder) []folderResponse {
	out := make([]folderResponse, 0, len(folders))
	for _, f := range folders {
		out = append(out, toFolderResponse(f))
	}
	return out
}
