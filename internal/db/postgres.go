package db

import (
	"fmt"
	"log"

	"github.com/rangira25/user_service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgres(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode, cfg.DBTimeZone,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf(" Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf(" Failed to get database instance: %v", err)
	}

	if err := sqlDB.Ping(); err != nil {
		log.Fatalf(" Database unreachable: %v", err)
	}

	log.Println("Connected to PostgreSQL")

	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	db.Exec(`CREATE EXTENSION IF NOT EXISTS citext;`)

	return db
}
