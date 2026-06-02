package handler

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
)

func Health(c *gin.Context) {
	isMaintenance := os.Getenv("MAINTENANCE_MODE") == "true"
	response.OK(c, "ok", gin.H{
		"version":        "1.0.0",
		"is_maintenance": isMaintenance,
	})
}
