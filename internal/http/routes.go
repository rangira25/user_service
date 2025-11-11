package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rangira25/user_service/internal/pkg"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes(r *gin.Engine, h *Handler) {
	v1 := r.Group("/api/v1")
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1.POST("/login", h.Login)

	{
		u := v1.Group("/users")

		u.Use(pkg.JWTMiddleware())
		{

			u.GET("/:id", h.GetUser)
			u.GET("", h.ListUsers)
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
