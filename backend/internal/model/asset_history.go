package model

import "time"

type AssetHistory struct {
	ID              string
	AssetChangeType ChangeType
	ChangedAt       time.Time
	Description     string
}
