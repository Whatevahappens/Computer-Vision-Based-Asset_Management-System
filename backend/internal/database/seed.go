package database

import (
	"backend/internal/config"
	"backend/internal/model"
	"fmt"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SeedAdmin(cfg *config.Config) {
	var count int64
	DB.Model(&model.User{}).Where("role = ?", model.Admin).Count(&count)
	if count > 0 {
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(cfg.AdminPassword), 12)
	if err != nil {
		log.Fatalf("Failed to hash admin password: %v", err)
	}

	admin := model.User{
		ID:           uuid.New().String(),
		FirstName:    "System",
		LastName:     "Admin",
		Email:        cfg.AdminEmail,
		Username:     "admin",
		PasswordHash: string(hash),
		Phone:        "99001122",
		Status:       model.UserActive,
		Role:         model.Admin,
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Fatalf("Failed to seed admin: %v", err)
	}
	fmt.Println("✓ Admin user seeded (admin / " + cfg.AdminPassword + ")")
}
