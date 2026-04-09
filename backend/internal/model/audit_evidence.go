package model

import "time"

type AuditEvidence struct {
	ID           string
	FilePath     string
	CapturedAt   time.Time
	ModelName    string
	ModelVersion string
}
