// @title User Service API
// @version 1.0
// @description This is the User Service API for managing users.
// @BasePath /api/v1

package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"github.com/swaggo/files"
    "github.com/swaggo/gin-swagger"
    _ "github.com/rangira25/user_service/docs"

	"github.com/rangira25/user_service/internal/config"
	"github.com/rangira25/user_service/internal/db"
	"github.com/rangira25/user_service/internal/domain"
	"github.com/rangira25/user_service/internal/http"
	"github.com/rangira25/user_service/internal/pkg"
	"github.com/rangira25/user_service/internal/repository"
	"github.com/rangira25/user_service/internal/service"
)

func main() {
	_ = godotenv.Load() 

	cfg := config.LoadConfig() 

	d := db.ConnectPostgres(cfg) 

	// Ensure PostgreSQL extensions exist (run once)
	_ = d.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";").Error
	_ = d.Exec("CREATE EXTENSION IF NOT EXISTS citext;").Error

	// AutoMigrate model
	if err := d.AutoMigrate(&domain.User{}); err != nil {
		log.Fatalf("migrate: %v", err)
	}

	repo := repository.NewUserRepository(d)
	svc := service.NewUserService(repo)
	handler := http.NewHandler(svc)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(pkg.LoggingMiddleware())


	// register routes
	http.RegisterRoutes(r, handler)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	port := cfg.Port
	if port == "" {
		port = os.Getenv("PORT")
	}
	log.Printf("listening on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("server exit: %v", err)
	}
}
