package model

import "time"

type Asset struct {
	ID               string      `gorm:"primaryKey;size:36" json:"id"`
	Barcode          string      `gorm:"size:50;uniqueIndex;not null" json:"barcode"`
	SerialNumber     string      `gorm:"size:100" json:"serialNumber"`
	AssetName        string      `gorm:"size:300;not null" json:"assetName"`
	AcquisitionPrice int         `gorm:"not null" json:"acquisitionPrice"`
	AcquisitionDate  time.Time   `gorm:"not null" json:"acquisitionDate"`
	UsefulLifeMonths int         `gorm:"not null" json:"usefulLifeMonths"`
	CurrentValue     int         `gorm:"not null" json:"currentValue"`
	Status           AssetStatus `gorm:"size:30;default:ACTIVE" json:"status"`
	Nature           AssetNature `gorm:"size:20;default:TANGIBLE" json:"nature"`
	Description      string      `gorm:"size:1000" json:"description"`
	AssetModelID     *string     `gorm:"size:36" json:"assetModelId"`
	AssetModel       *AssetModel `gorm:"foreignKey:AssetModelID" json:"assetModel,omitempty"`
	DepartmentID     *string     `gorm:"size:36" json:"departmentId"`
	Department       *Department `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	LocationID       *string     `gorm:"size:36" json:"locationId"`
	Location         *Location   `gorm:"foreignKey:LocationID" json:"location,omitempty"`
	AssignedUserID   *string     `gorm:"size:36" json:"assignedUserId"`
	AssignedUser     *User       `gorm:"foreignKey:AssignedUserID" json:"assignedUser,omitempty"`
	CreatedAt        time.Time   `json:"createdAt"`
	UpdatedAt        time.Time   `json:"updatedAt"`
}
