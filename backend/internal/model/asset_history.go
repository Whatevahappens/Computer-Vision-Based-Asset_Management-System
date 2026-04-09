package model

import "time"

type ChangeType string

const (
	Created     ChangeType = "CREATED"
	Updated     ChangeType = "UPDATED"
	Assigned    ChangeType = "ASSIGNED"
	Transferred ChangeType = "TRANSFERRED"
	Revalued    ChangeType = "REVALUED"
	Depreciated ChangeType = "DEPRECIATED"
	Disposed    ChangeType = "DISPOSED"
)

type AssetHistory struct {
	ID          string
	ChangeType  ChangeType
	ChangedAt   time.Time
	Description string
}
