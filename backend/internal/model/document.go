package model

import "time"

type Document struct {
	ID         string       `gorm:"primaryKey;size:36" json:"id"`
	Name       string       `gorm:"size:300;not null" json:"name"`
	FilePath   string       `gorm:"size:500;not null" json:"filePath"`
	FileType   DocumentType `gorm:"size:20" json:"fileType"`
	FileSize   int64        `gorm:"default:0" json:"fileSize"`
	UploadDate time.Time    `json:"uploadDate"`
	AssetID    *string      `gorm:"size:36;index" json:"assetId"`
	Asset      *Asset       `gorm:"foreignKey:AssetID" json:"-"`
	UploadedBy string       `gorm:"size:36;not null" json:"uploadedBy"`
}
