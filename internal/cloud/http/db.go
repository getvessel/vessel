package http

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"vessl.dev/vessl/internal/models"
)

func InitDatabase() *gorm.DB {
	dsn := os.Getenv("CLOUD_DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=vessl password=vessl dbname=vesslcloud port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to cloud database: %v", err)
	}

	if err := db.AutoMigrate(
		&models.CloudTeam{},
		&models.CloudServer{},
		&models.CloudUsageLog{},
		&models.CloudTelemetryLog{},
	); err != nil {
		log.Fatalf("Failed to run cloud database migrations: %v", err)
	}

	return db
}
