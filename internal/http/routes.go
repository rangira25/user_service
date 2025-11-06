package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rangira25/user_service/internal/pkg"
)

func RegisterRoutes(r *gin.Engine, h *Handler) {
	v1 := r.Group("/api/v1")
	{
		u := v1.Group("/users")
		{
			// Public / user-accessible routes
			u.GET("", h.ListUsers)
			u.GET("/:id", h.GetUser)
			u.POST("", h.CreateUser)

			// Admin-only routes
			admin := u.Group("")
			admin.Use(pkg.RequireRole("admin")) 
			{
				
				admin.PATCH("/:id", h.UpdateUser)
				admin.PUT("/:id/status", h.SetStatus)
				admin.DELETE("/:id", h.DeleteUser)
				admin.POST("/:id/restore", h.RestoreUser)
				admin.POST("/:id/reset-password", h.AdminResetPassword)
			}
		}
	}
}
