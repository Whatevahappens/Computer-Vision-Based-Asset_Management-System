package model

import "time"

type Notification struct {
	ID        string
	Title     string
	Message   string
	Type      NotificationType
	IsRead    bool
	CreatedAt time.Time
}
