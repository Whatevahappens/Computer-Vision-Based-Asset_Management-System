package model

import "time"

type AuditEvidence struct {
	ID             string    `gorm:"primaryKey;size:36" json:"id"`
	FilePath       string    `gorm:"size:500;not null" json:"filePath"`
	CapturedAt     time.Time `json:"capturedAt"`
	ModelName      string    `gorm:"size:100" json:"modelName"`
	ModelVersion   string    `gorm:"size:50" json:"modelVersion"`
	AuditFindingID string    `gorm:"size:36;not null;index" json:"auditFindingId"`
}
