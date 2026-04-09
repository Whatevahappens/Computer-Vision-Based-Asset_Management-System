package model

import "time"

type Asset struct {
	ID               string
	Barcode          string
	SerialNumber     string
	AssetName        string
	AcquisitionPrice int
	AcquisitionDate  time.Time
	UsefulLifeMonths int
	CurrentValue     int
	Status           AssetStatus
	Nature           AssetNature
	Description      string
}
