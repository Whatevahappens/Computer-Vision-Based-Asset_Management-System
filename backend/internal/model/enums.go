package model

type UserStatus string

const (
	Active   UserStatus = "ACTIVE"
	Inactive UserStatus = "INACTIVE"
	Banned   UserStatus = "BANNED"
)

type Role string

const (
	Admin          Role = "ADMIN"
	Accountant     Role = "ACCOUNTANT"
	AssetCustodian Role = "ASSET_CUSTODIAN"
	Employee       Role = "EMPLOYEE"
)
