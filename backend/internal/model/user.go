package model

import "time"

type User struct {
	ID           string      `gorm:"primaryKey;size:36" json:"id"`
	FirstName    string      `gorm:"size:100;not null" json:"firstName"`
	LastName     string      `gorm:"size:100;not null" json:"lastName"`
	Email        string      `gorm:"size:255;uniqueIndex;not null" json:"email"`
	Username     string      `gorm:"size:100;uniqueIndex;not null" json:"username"`
	PasswordHash string      `gorm:"size:255;not null" json:"-"`
	Phone        string      `gorm:"size:20" json:"phone"`
	Status       UserStatus  `gorm:"size:20;default:ACTIVE" json:"status"`
	Role         Role        `gorm:"size:30;not null" json:"role"`
	DepartmentID *string     `gorm:"size:36" json:"departmentId"`
	Department   *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	CreatedAt    time.Time   `json:"createdAt"`
	UpdatedAt    time.Time   `json:"updatedAt"`
}
