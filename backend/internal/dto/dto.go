package dto

import "time"

// Auth
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required"`
	NewPassword     string `json:"newPassword" binding:"required,min=6"`
}

// User
type CreateUserRequest struct {
	FirstName    string `json:"firstName" binding:"required"`
	LastName     string `json:"lastName" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Username     string `json:"username" binding:"required"`
	Password     string `json:"password" binding:"required,min=6"`
	Phone        string `json:"phone"`
	Role         string `json:"role" binding:"required"`
	DepartmentID string `json:"departmentId"`
}

type UpdateUserRequest struct {
	FirstName    string `json:"firstName"`
	LastName     string `json:"lastName"`
	Phone        string `json:"phone"`
	Status       string `json:"status"`
	DepartmentID string `json:"departmentId"`
}

type UserResponse struct {
	ID           string  `json:"id"`
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	Email        string  `json:"email"`
	Username     string  `json:"username"`
	Phone        string  `json:"phone"`
	Status       string  `json:"status"`
	Role         string  `json:"role"`
	DepartmentID *string `json:"departmentId"`
}

// Asset
type CreateAssetRequest struct {
	AssetName        string `json:"assetName" binding:"required"`
	SerialNumber     string `json:"serialNumber"`
	AcquisitionPrice int    `json:"acquisitionPrice" binding:"required"`
	AcquisitionDate  string `json:"acquisitionDate" binding:"required"`
	UsefulLifeMonths int    `json:"usefulLifeMonths" binding:"required"`
	Nature           string `json:"nature"`
	Description      string `json:"description"`
	AssetModelID     string `json:"assetModelId"`
	DepartmentID     string `json:"departmentId"`
	LocationID       string `json:"locationId"`
}

type UpdateAssetRequest struct {
	AssetName    string `json:"assetName"`
	Description  string `json:"description"`
	Status       string `json:"status"`
	LocationID   string `json:"locationId"`
	DepartmentID string `json:"departmentId"`
}

type AssignAssetRequest struct {
	UserID     string `json:"userId" binding:"required"`
	LocationID string `json:"locationId"`
	Notes      string `json:"notes"`
}

type TransferAssetRequest struct {
	ToUserID   string `json:"toUserId" binding:"required"`
	LocationID string `json:"locationId"`
	Notes      string `json:"notes"`
}

type DisposeAssetRequest struct {
	Reason        string `json:"reason" binding:"required"`
	ResidualValue int    `json:"residualValue"`
	Notes         string `json:"notes"`
}

// Asset Model
type CreateAssetModelRequest struct {
	Brand                   string `json:"brand" binding:"required"`
	ModelName               string `json:"modelName" binding:"required"`
	AssetType               string `json:"assetType" binding:"required"`
	Category                string `json:"category" binding:"required"`
	DefaultUnitPrice        int    `json:"defaultUnitPrice"`
	DefaultUsefulLifeMonths int    `json:"defaultUsefulLifeMonths"`
	DepreciationMethod      string `json:"depreciationMethod"`
}

// Location
type CreateLocationRequest struct {
	Name     string `json:"name" binding:"required"`
	Building string `json:"building"`
	Floor    string `json:"floor"`
	Room     string `json:"room"`
	Capacity int    `json:"capacity"`
}

// Department
type CreateDepartmentRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// Audit
type StartAuditRequest struct {
	LocationID string `json:"locationId" binding:"required"`
	Notes      string `json:"notes"`
}

type CVDetection struct {
	ClassName  string  `json:"class_name"`
	Confidence float64 `json:"confidence"`
	Box        []int   `json:"box"`
}

type CVDetectionResponse struct {
	Detections []CVDetection `json:"detections"`
	ImagePath  string        `json:"image_path"`
	ModelName  string        `json:"model_name"`
	ModelVer   string        `json:"model_version"`
}

// Depreciation
type CalculateDepreciationRequest struct {
	AssetID string `json:"assetId" binding:"required"`
	Method  string `json:"method" binding:"required"`
}

type DepreciationResult struct {
	AssetID          string  `json:"assetId"`
	AssetName        string  `json:"assetName"`
	AcquisitionPrice int     `json:"acquisitionPrice"`
	CurrentValue     int     `json:"currentValue"`
	MonthlyAmount    float64 `json:"monthlyAmount"`
	Method           string  `json:"method"`
	UsefulLifeMonths int     `json:"usefulLifeMonths"`
}

// Revaluation
type RevalueAssetRequest struct {
	AssetID  string `json:"assetId" binding:"required"`
	NewValue int    `json:"newValue" binding:"required"`
	Reason   string `json:"reason" binding:"required"`
}

// Report
type ReportRequest struct {
	ReportType string `json:"reportType" binding:"required"`
	Format     string `json:"format" binding:"required"`
	StartDate  string `json:"startDate"`
	EndDate    string `json:"endDate"`
	LocationID string `json:"locationId"`
}

// Notification
type NotificationResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Type      string    `json:"type"`
	IsRead    bool      `json:"isRead"`
	CreatedAt time.Time `json:"createdAt"`
}

// Pagination
type PaginationQuery struct {
	Page   int    `form:"page,default=1"`
	Limit  int    `form:"limit,default=20"`
	Search string `form:"search"`
	Status string `form:"status"`
	Sort   string `form:"sort,default=created_at"`
	Order  string `form:"order,default=desc"`
}

func (p *PaginationQuery) Offset() int {
	return (p.Page - 1) * p.Limit
}

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	TotalPages int         `json:"totalPages"`
}
