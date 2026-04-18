package model

import "time"

type User struct {
	ID             uint       `json:"id" gorm:"primaryKey"`
	Firstname      string     `json:"firstname" gorm:"not null"`
	Lastname       string     `json:"lastname" gorm:"not null"`
	Email          string     `json:"email" gorm:"uniqueIndex;not null"`
	Username       string     `json:"username" gorm:"uniqueIndex;not null"`
	Password       string     `json:"-"`
	Phone          string     `json:"phone"`
	Status         UserStatus `json:"status"`
	UserRole       Role       `json:"user_role"`
	ResetToken     string     `json:"-"`
	ResetExpiresAt *time.Time `json:"-"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}
