package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rangira25/user_service/internal/config"
	"github.com/rangira25/user_service/internal/db"
	"github.com/rangira25/user_service/internal/domain"
)

func main() {
	_ = godotenv.Load()
	cfg := config.LoadConfig()
	d := db.ConnectPostgres(cfg)

	users := []domain.User{
		{FullName: "Admin User", Email: "admin@example.com", Role: "admin", Status: "active"},
		{FullName: "Regular User", Email: "user@example.com", Role: "user", Status: "active"},
	}

	for _, u := range users {
		if err := d.Create(&u).Error; err != nil {
			log.Printf("Skipping user %s: %v", u.Email, err)
		} else {
			log.Printf("Inserted user: %s", u.Email)
		}
	}
}
