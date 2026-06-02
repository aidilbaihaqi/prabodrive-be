package response

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	TotalData  int `json:"total_data"`
	TotalPages int `json:"total_pages"`
}

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

func success(c *gin.Context, code int, message string, data any) {
	c.JSON(code, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func fail(c *gin.Context, code int, message string) {
	c.JSON(code, APIResponse{
		Status:  "error",
		Message: message,
	})
}

// OK sends 200 with data (single resource or custom payload).
func OK(c *gin.Context, message string, data any) {
	success(c, http.StatusOK, message, data)
}

// OKList sends 200 with an array payload and pagination meta.
func OKList(c *gin.Context, message string, data any, page, limit, total int) {
	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
		Meta: &Meta{
			Page:       page,
			Limit:      limit,
			TotalData:  total,
			TotalPages: totalPages(total, limit),
		},
	})
}

// Created sends 201.
func Created(c *gin.Context, message string, data any) {
	success(c, http.StatusCreated, message, data)
}

// Updated sends 200 for mutation responses.
func Updated(c *gin.Context, message string, data any) {
	success(c, http.StatusOK, message, data)
}

// Deleted sends 200 for delete responses.
func Deleted(c *gin.Context, message string, data any) {
	success(c, http.StatusOK, message, data)
}

// BadRequest sends 400.
func BadRequest(c *gin.Context, message string) {
	fail(c, http.StatusBadRequest, message)
}

// Unauthorized sends 401.
func Unauthorized(c *gin.Context) {
	fail(c, http.StatusUnauthorized, "unauthorized")
}

// Forbidden sends 403.
func Forbidden(c *gin.Context, message string) {
	fail(c, http.StatusForbidden, message)
}

// NotFound sends 404.
func NotFound(c *gin.Context) {
	fail(c, http.StatusNotFound, "not found")
}

// TooManyRequests sends 429.
func TooManyRequests(c *gin.Context) {
	fail(c, http.StatusTooManyRequests, "too many requests")
}

// Maintenance sends 503.
func Maintenance(c *gin.Context) {
	c.JSON(http.StatusServiceUnavailable, APIResponse{
		Status:  "error",
		Message: "service is under maintenance, please try again later",
	})
}

// InternalError sends 500 and logs the underlying error.
func InternalError(c *gin.Context, err error) {
	log.Printf("[ERROR] %v", err)
	fail(c, http.StatusInternalServerError, "internal server error")
}

func totalPages(total, limit int) int {
	if limit <= 0 {
		return 0
	}
	return (total + limit - 1) / limit
}
