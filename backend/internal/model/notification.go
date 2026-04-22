package model

import "time"

type Notification struct {
	ID        string           `gorm:"primaryKey;size:36" json:"id"`
	Title     string           `gorm:"size:300;not null" json:"title"`
	Message   string           `gorm:"size:2000" json:"message"`
	Type      NotificationType `gorm:"size:20;default:INFO" json:"type"`
	IsRead    bool             `gorm:"default:false" json:"isRead"`
	CreatedAt time.Time        `json:"createdAt"`
	UserID    string           `gorm:"size:36;not null;index" json:"userId"`
	User      User             `gorm:"foreignKey:UserID" json:"-"`
}
