package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/response"
	"github.com/aidilbaihaqi/prabodrive-be/internal/shared/constants"
)

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.GetString("user_role") != constants.RoleAdmin {
			response.Forbidden(c, "admin access required")
			c.Abort()
			return
		}
		c.Next()
	}
}
