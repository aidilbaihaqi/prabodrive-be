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

type FolderHandler struct {
	folderUC usecase.FolderUsecase
}

func NewFolderHandler(folderUC usecase.FolderUsecase) *FolderHandler {
	return &FolderHandler{folderUC: folderUC}
}

func (h *FolderHandler) List(c *gin.Context) {
	var q request.ListFoldersQuery
	_ = c.ShouldBindQuery(&q)

	userID := c.GetString("user_id")

	// Resolve parent filter: nil = all, "" = root, uuid = children
	var parentFilter *string
	if q.ParentID != nil {
		if *q.ParentID == "root" {
			empty := ""
			parentFilter = &empty
		} else {
			parentFilter = q.ParentID
		}
	}

	folders, err := h.folderUC.List(c.Request.Context(), userID, parentFilter)
	if err != nil {
		response.InternalError(c, err)
		return
	}
	response.OK(c, "folders fetched", toFolderList(folders))
}

func (h *FolderHandler) Get(c *gin.Context) {
	userID := c.GetString("user_id")
	folder, err := h.folderUC.Get(c.Request.Context(), c.Param("id"), userID)
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
	folder, err := h.folderUC.Create(c.Request.Context(), userID, req.Name, req.ParentID)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.Created(c, "folder created", toFolderResponse(folder))
}

func (h *FolderHandler) Update(c *gin.Context) {
	var req request.UpdateFolderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	folderID := c.Param("id")

	if err := h.folderUC.Update(c.Request.Context(), folderID, userID, req.Name); err != nil {
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

	if err := h.folderUC.Delete(c.Request.Context(), folderID, userID); err != nil {
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
