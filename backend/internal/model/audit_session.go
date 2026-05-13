package model

import "time"

type AuditSession struct {
	ID          string         `gorm:"primaryKey;size:36" json:"id"`
	StartedAt   time.Time      `gorm:"not null" json:"startedAt"`
	EndedAt     *time.Time     `json:"endedAt"`
	Status      AuditStatus    `gorm:"size:20;default:PLANNED" json:"status"`
	Notes       string         `gorm:"size:2000" json:"notes"`
	LocationID  string         `gorm:"size:36;not null" json:"locationId"`
	Location    Location       `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	PerformedBy string         `gorm:"size:36;not null" json:"performedBy"`
	Performer   User           `gorm:"foreignKey:PerformedBy" json:"performer,omitempty"`
	Findings    []AuditFinding `gorm:"foreignKey:AuditSessionID" json:"findings,omitempty"`
	Summaries   []AuditSummary `gorm:"foreignKey:AuditSessionID" json:"summaries,omitempty"`
}
