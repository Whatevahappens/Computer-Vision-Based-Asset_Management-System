package database

import (
	"backend/internal/config"
	"backend/internal/model"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect(cfg *config.Config) {
	var err error
	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Database connected")
}

func Migrate() {
	err := DB.AutoMigrate(
		&model.User{},
		&model.Department{},
		&model.Location{},
		&model.AssetModel{},
		&model.Asset{},
		&model.AssetHistory{},
		&model.Document{},
		&model.AuditSession{},
		&model.AuditFinding{},
		&model.AuditEvidence{},
		&model.AuditSummary{},
		&model.Notification{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}
	fmt.Println("Database migrated")
}
