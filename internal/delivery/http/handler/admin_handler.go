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

type AdminHandler struct {
	adminUC usecase.AdminUsecase
}

func NewAdminHandler(adminUC usecase.AdminUsecase) *AdminHandler {
	return &AdminHandler{adminUC: adminUC}
}

func (h *AdminHandler) ListUsers(c *gin.Context) {
	var q request.PaginationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	page, limit := clampPage(q.Page, q.Limit)
	users, total, err := h.adminUC.ListUsers(c.Request.Context(), page, limit)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.OKList(c, "users fetched", toUserList(users), page, limit, total)
}

func (h *AdminHandler) GetUser(c *gin.Context) {
	user, err := h.adminUC.GetUser(c.Request.Context(), c.Param("id"))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			response.NotFound(c)
			return
		}
		response.InternalError(c, err)
		return
	}
	response.OK(c, "user fetched", toUserResponse(user))
}

func (h *AdminHandler) UpdateRole(c *gin.Context) {
	var req request.UpdateUserRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	targetID := c.Param("id")
	requesterID := c.GetString("user_id")

	if err := h.adminUC.UpdateRole(c.Request.Context(), targetID, requesterID, req.Role); err != nil {
		switch {
		case errors.Is(err, domain.ErrForbidden):
			response.Forbidden(c, "cannot change your own role")
		case errors.Is(err, domain.ErrUserNotFound):
			response.NotFound(c)
		default:
			response.InternalError(c, err)
		}
		return
	}

	response.Updated(c, "user role updated", gin.H{"id": targetID})
}

func (h *AdminHandler) DeleteUser(c *gin.Context) {
	targetID := c.Param("id")
	requesterID := c.GetString("user_id")

	if err := h.adminUC.DeleteUser(c.Request.Context(), targetID, requesterID); err != nil {
		switch {
		case errors.Is(err, domain.ErrForbidden):
			response.Forbidden(c, "cannot delete your own account")
		case errors.Is(err, domain.ErrUserNotFound):
			response.NotFound(c)
		default:
			response.InternalError(c, err)
		}
		return
	}

	response.Deleted(c, "user deleted", gin.H{"id": targetID})
}

type userResponse struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	QuotaUsed int64  `json:"quota_used"`
	QuotaMax  int64  `json:"quota_max"`
	CreatedAt string `json:"created_at"`
}

func toUserResponse(u *domain.User) userResponse {
	return userResponse{
		ID:        u.ID,
		Email:     u.Email,
		Name:      u.Name,
		Role:      u.Role,
		QuotaUsed: u.QuotaUsed,
		QuotaMax:  u.QuotaMax,
		CreatedAt: u.CreatedAt.Format(time.RFC3339),
	}
}

func toUserList(users []*domain.User) []userResponse {
	out := make([]userResponse, 0, len(users))
	for _, u := range users {
		out = append(out, toUserResponse(u))
	}
	return out
}
