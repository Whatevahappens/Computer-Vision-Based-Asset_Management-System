package model

import "time"

type AssetHistory struct {
	ID          string     `gorm:"primaryKey;size:36" json:"id"`
	ChangeType  ChangeType `gorm:"size:30;not null" json:"changeType"`
	ChangedAt   time.Time  `gorm:"not null" json:"changedAt"`
	Description string     `gorm:"size:1000" json:"description"`
	AssetID     string     `gorm:"size:36;not null;index" json:"assetId"`
	Asset       Asset      `gorm:"foreignKey:AssetID" json:"-"`
	UserID      string     `gorm:"size:36;not null" json:"userId"`
	User        User       `gorm:"foreignKey:UserID" json:"-"`
}
