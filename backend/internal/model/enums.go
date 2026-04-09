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
	Other    DocumentType = "OTHER"
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
