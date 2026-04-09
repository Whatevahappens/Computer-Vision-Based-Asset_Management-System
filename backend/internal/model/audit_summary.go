package model

type AuditSummary struct {
	ID              string
	Category        AssetCategory
	RegisteredCount int
	DetectedCount   int
	Difference      int
}
