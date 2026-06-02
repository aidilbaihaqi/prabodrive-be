package request

import "time"

type CreateShareRequest struct {
	DocumentID string     `json:"document_id" binding:"required,uuid"`
	ExpiresAt  time.Time  `json:"expires_at" binding:"required"`
	Password   *string    `json:"password"`
}

type AccessShareQuery struct {
	Password string `form:"password"`
}

type PaginationQuery struct {
	Page  int `form:"page,default=1"`
	Limit int `form:"limit,default=20"`
}
