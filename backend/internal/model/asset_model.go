package model

import "time"

type AssetModel struct {
	ID                      string             `gorm:"primaryKey;size:36" json:"id"`
	Brand                   string             `gorm:"size:200;not null" json:"brand"`
	ModelName               string             `gorm:"size:200;not null" json:"modelName"`
	AssetModelType          AssetType          `gorm:"size:30" json:"assetType"`
	Category                AssetCategory      `gorm:"size:30" json:"category"`
	DefaultUnitPrice        int                `gorm:"default:0" json:"defaultUnitPrice"`
	DefaultUsefulLifeMonths int                `gorm:"default:60" json:"defaultUsefulLifeMonths"`
	DepreciationMethod      DepreciationMethod `gorm:"size:30;default:STRAIGHT_LINE" json:"depreciationMethod"`
	CreatedAt               time.Time          `json:"createdAt"`
	UpdatedAt               time.Time          `json:"updatedAt"`
}
