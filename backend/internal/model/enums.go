package model

type UserStatus string

const (
	UserActive   UserStatus = "ACTIVE"
	UserInactive UserStatus = "INACTIVE"
	UserBanned   UserStatus = "BANNED"
)

type Role string

const (
	Admin          Role = "ADMIN"
	Accountant     Role = "ACCOUNTANT"
	AssetCustodian Role = "ASSET_CUSTODIAN"
	Employee       Role = "EMPLOYEE"
)

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

type DocumentType string

const (
	Invoice  DocumentType = "INVOICE"
	Receipt  DocumentType = "RECEIPT"
	Warranty DocumentType = "WARRANTY"
	Image    DocumentType = "IMAGE"
	Report   DocumentType = "REPORT"
	OtherDoc DocumentType = "OTHER"
)

type AuditStatus string

const (
	Planned    AuditStatus = "PLANNED"
	InProgress AuditStatus = "IN_PROGRESS"
	Completed  AuditStatus = "COMPLETED"
	Cancelled  AuditStatus = "CANCELLED"
)

type FindingType string

const (
	Matched      FindingType = "MATCHED"
	Missing      FindingType = "MISSING"
	Unregistered FindingType = "UNREGISTERED"
	Damaged      FindingType = "DAMAGED"
	Misplaced    FindingType = "MISPLACED"
)

type NotificationType string

const (
	Info    NotificationType = "INFO"
	Warning NotificationType = "WARNING"
	Alert   NotificationType = "ALERT"
)

type DepreciationMethod string

const (
	StraightLine     DepreciationMethod = "STRAIGHT_LINE"
	DecliningBalance DepreciationMethod = "DECLINING_BALANCE"
)

type AssetCategory string

const (
	ItEquipment     AssetCategory = "IT_EQUIPMENT"
	OfficeEquipment AssetCategory = "OFFICE_EQUIPMENT"
	Furniture       AssetCategory = "FURNITURE"
	Vehicle         AssetCategory = "VEHICLE"
	OtherAsset      AssetCategory = "OTHER"
)

type AssetStatus string

const (
	AssetActive           AssetStatus = "ACTIVE"
	AssetInactive         AssetStatus = "INACTIVE"
	AssetDisposed         AssetStatus = "DISPOSED"
	AssetLost             AssetStatus = "LOST"
	AssetUnderMaintenance AssetStatus = "UNDER_MAINTENANCE"
)

type AssetNature string

const (
	Tangible   AssetNature = "TANGIBLE"
	Intangible AssetNature = "INTANGIBLE"
)

type AssetType string

const (
	AssetEquipment  AssetType = "EQUIPMENT"
	AssetFurniture  AssetType = "FURNITURE"
	AssetVehicle    AssetType = "VEHICLE"
	AssetElectronic AssetType = "ELECTRONIC"
	AssetOther      AssetType = "OTHER"
)
