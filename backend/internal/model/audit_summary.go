package model

type AuditSummary struct {
	ID              string        `gorm:"primaryKey;size:36" json:"id"`
	Category        AssetCategory `gorm:"size:30" json:"category"`
	RegisteredCount int           `json:"registeredCount"`
	DetectedCount   int           `json:"detectedCount"`
	Difference      int           `json:"difference"`
	AuditSessionID  string        `gorm:"size:36;not null;index" json:"auditSessionId"`
}
