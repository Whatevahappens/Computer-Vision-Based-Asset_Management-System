package model

type AuditFinding struct {
	ID              string          `gorm:"primaryKey;size:36" json:"id"`
	Type            FindingType     `gorm:"size:20;not null" json:"type"`
	Confidence      float64         `gorm:"default:0" json:"confidence"`
	Notes           string          `gorm:"size:1000" json:"notes"`
	AuditSessionID  string          `gorm:"size:36;not null;index" json:"auditSessionId"`
	ExpectedAssetID *string         `gorm:"size:36" json:"expectedAssetId"`
	DetectedAssetID *string         `gorm:"size:36" json:"detectedAssetId"`
	Evidence        []AuditEvidence `gorm:"foreignKey:AuditFindingID" json:"evidence,omitempty"`
}
