package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
)

func Health(c *gin.Context) {
	response.OK(c, "ok", gin.H{"version": "1.0.0"})
}
