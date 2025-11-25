package http

import (
	"github.com/gin-gonic/gin"
	"github.com/rangira25/user_service/internal/pkg"
 
)

func RegisterRoutes(r *gin.Engine, h *Handler) {
    v1 := r.Group("/api/v1")

    // Swagger
    

    // Public route
    v1.POST("/login", h.Login)

    // Protected routes
    u := v1.Group("/users")
    u.Use(pkg.JWTMiddleware())
    {
        // User actions
        u.GET("/:id", h.GetUser)
        u.GET("", h.ListUsers)
        u.POST("", h.CreateUser)

        // Admin-only prefix
        admin := u.Group("/admin")
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


