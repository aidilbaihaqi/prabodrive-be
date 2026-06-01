package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/yourname/yourapp/internal/delivery/http/handler"
	"github.com/yourname/yourapp/internal/middleware"
)

// RegisterUserRoutes registers all user-related routes
func RegisterUserRoutes(r *gin.Engine, h *handler.UserHandler, jwtSecret string) {
	// API version group
	v1 := r.Group("/api/v1")

	// Public routes (no authentication required)
	// Add public routes here if needed

	// Protected routes (authentication required)
	users := v1.Group("/users")
	users.Use(middleware.Auth(jwtSecret))
	{
		users.POST("", h.Create)
		users.GET("", h.List)
		users.GET("/:id", h.Get)
		users.PUT("/:id", h.Update)
		users.DELETE("/:id", h.Delete)
	}
}

// RegisterHealthRoutes registers health check routes
func RegisterHealthRoutes(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "Service is running",
		})
	})

	r.GET("/ready", func(c *gin.Context) {
		// Add readiness checks here (db connection, etc.)
		c.JSON(200, gin.H{
			"status": "ready",
		})
	})
}
