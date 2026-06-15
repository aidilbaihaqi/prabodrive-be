package handler

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/request"
	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
	"github.com/aidilbaihaqi/prabodrive-be/internal/domain"
)

type ActivityHandler struct {
	activity domain.ActivityRepository
}

func NewActivityHandler(activity domain.ActivityRepository) *ActivityHandler {
	return &ActivityHandler{activity: activity}
}

func (h *ActivityHandler) List(c *gin.Context) {
	var q request.PaginationQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	userID := c.GetString("user_id")
	page, limit := clampPage(q.Page, q.Limit)

	logs, total, err := h.activity.List(c.Request.Context(), userID, page, limit)
	if err != nil {
		response.InternalError(c, err)
		return
	}

	response.OKList(c, "activity fetched", toLogList(logs), page, limit, total)
}

type activityResponse struct {
	ID           string  `json:"id"`
	Action       string  `json:"action"`
	DocumentID   *string `json:"document_id"`
	DocumentName *string `json:"document_name"`
	IPAddress    *string `json:"ip_address"`
	CreatedAt    string  `json:"created_at"`
}

func toLogList(logs []*domain.ActivityLog) []activityResponse {
	out := make([]activityResponse, 0, len(logs))
	for _, l := range logs {
		out = append(out, activityResponse{
			ID:           l.ID,
			Action:       l.Action,
			DocumentID:   l.DocumentID,
			DocumentName: l.DocumentName,
			IPAddress:    l.IPAddress,
			CreatedAt:    l.CreatedAt.Format(time.RFC3339),
		})
	}
	return out
}
