package model

type AuditFinding struct {
	ID         string
	Type       FindingType
	Confidence float64
	Notes      string
}
