package model

import "time"

type Location struct {
	ID        string    `gorm:"primaryKey;size:36" json:"id"`
	Name      string    `gorm:"size:200;not null" json:"name"`
	Building  string    `gorm:"size:200" json:"building"`
	Floor     string    `gorm:"size:50" json:"floor"`
	Room      string    `gorm:"size:100" json:"room"`
	Capacity  int       `gorm:"default:0" json:"capacity"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
