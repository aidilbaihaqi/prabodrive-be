package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/aidilbaihaqi/prabodrive-be/internal/delivery/http/handler"
	"github.com/aidilbaihaqi/prabodrive-be/internal/middleware"
)

func Register(
	r *gin.Engine,
	jwtSecret string,
	authH *handler.AuthHandler,
	docH *handler.DocumentHandler,
	folderH *handler.FolderHandler,
	shareH *handler.ShareHandler,
	activityH *handler.ActivityHandler,
	adminH *handler.AdminHandler,
) {
	r.GET("/health", handler.Health)

	api := r.Group("/api/v1")

	// Auth
	auth := api.Group("/auth")
	{
		auth.POST("/register", authH.Register)
		auth.POST("/login", authH.Login)
		auth.POST("/refresh", authH.Refresh)
		auth.POST("/logout", middleware.Auth(jwtSecret), authH.Logout)
	}

	// Protected user routes
	protected := api.Group("/", middleware.Auth(jwtSecret))
	{
		protected.GET("/me", authH.GetProfile)
		protected.PATCH("/me", authH.UpdateProfile)

		protected.GET("/documents", docH.List)
		protected.GET("/documents/:id", docH.Get)
		protected.POST("/documents/presign-upload", docH.PresignUpload)
		protected.POST("/documents/confirm-upload", docH.ConfirmUpload)
		protected.PATCH("/documents/:id", docH.Rename)
		protected.DELETE("/documents/:id", docH.Delete)
		protected.GET("/documents/:id/download", docH.Download)

		protected.GET("/folders", folderH.List)
		protected.GET("/folders/:id", folderH.Get)
		protected.POST("/folders", folderH.Create)
		protected.PATCH("/folders/:id", folderH.Update)
		protected.DELETE("/folders/:id", folderH.Delete)

		protected.GET("/share", shareH.List)
		protected.POST("/share", shareH.Create)
		protected.DELETE("/share/:id", shareH.Delete)

		protected.GET("/activity", activityH.List)
	}

	// Public share access
	api.GET("/share/:token", shareH.Access)

	// Admin routes — requires auth + admin role
	admin := api.Group("/admin", middleware.Auth(jwtSecret), middleware.RequireAdmin())
	{
		admin.GET("/users", adminH.ListUsers)
		admin.GET("/users/:id", adminH.GetUser)
		admin.PATCH("/users/:id/role", adminH.UpdateRole)
		admin.DELETE("/users/:id", adminH.DeleteUser)
	}
}
