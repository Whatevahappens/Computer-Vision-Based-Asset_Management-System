package model

import "time"

type AuditSession struct {
	ID        string
	StartedAt time.Time
	EndedAt   time.Time
	Status    AuditStatus
	Notes     string
}
