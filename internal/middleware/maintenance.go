package middleware

import (
	"os"

	"github.com/gin-gonic/gin"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
)

func MaintenanceMode() gin.HandlerFunc {
	return func(c *gin.Context) {
		if os.Getenv("MAINTENANCE_MODE") == "true" {
			if c.FullPath() == "/health" {
				c.Next()
				return
			}
			response.Maintenance(c)
			c.Abort()
			return
		}
		c.Next()
	}
}
