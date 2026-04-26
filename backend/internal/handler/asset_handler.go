package handler

import (
	"backend/internal/dto"
	"backend/internal/middleware"
	"backend/internal/model"
	"backend/internal/repository"
	"backend/internal/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func CreateAsset(c *gin.Context) {
	var req dto.CreateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	acqDate, err := time.Parse("2006-01-02", req.AcquisitionDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM_DD"})
		return
	}

	userID := middleware.GetUserID(c)
	nature := req.Nature
	if nature == "" {
		nature = string(model.Tangible)
	}

	asset, err := service.CreateAsset(
		req.AssetName, req.SerialNumber, req.AcquisitionPrice, acqDate,
		req.UsefulLifeMonths, nature, req.Description,
		req.AssetModelID, req.DepartmentID, req.LocationID, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// re-fetch with preloads
	full, _ := repository.FindAssetByID(asset.ID)
	c.JSON(http.StatusCreated, full)
}

func GetAsset(c *gin.Context) {
	asset, err := repository.FindAssetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}
	c.JSON(http.StatusOK, asset)
}

func ListAssets(c *gin.Context) {
	var q dto.PaginationQuery
	c.ShouldBindQuery(&q)
	if q.Limit <= 0 {
		q.Limit = 20
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	assets, total, err := repository.ListAssets(q.Offset(), q.Limit, q.Search, q.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pages := int(total) / q.Limit
	if int(total)%q.Limit != 0 {
		pages++
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data: assets, Total: total, Page: q.Page, Limit: q.Limit, TotalPages: pages,
	})
}

func UpdateAsset(c *gin.Context) {
	asset, err := repository.FindAssetByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "asset not found"})
		return
	}

	var req dto.UpdateAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.AssetName != "" {
		asset.AssetName = req.AssetName
	}
	if req.Description != "" {
		asset.Description = req.Description
	}
	if req.Status != "" {
		asset.Status = model.AssetStatus(req.Status)
	}
	if req.LocationID != "" {
		asset.LocationID = &req.LocationID
	}
	if req.DepartmentID != "" {
		asset.DepartmentID = &req.DepartmentID
	}

	if err := repository.UpdateAsset(asset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}

	userID := middleware.GetUserID(c)
	repository.CreateAssetHistory(&model.AssetHistory{
		ID: uuid.New().String(), ChangeType: model.Updated,
		ChangedAt: time.Now(), Description: "Asset updated",
		AssetID: asset.ID, UserID: userID,
	})

	full, _ := repository.FindAssetByID(asset.ID)
	c.JSON(http.StatusOK, full)
}

func AssignAsset(c *gin.Context) {
	var req dto.AssignAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := middleware.GetUserID(c)
	if err := service.AssignAsset(c.Param("id"), req.UserID, req.LocationID, req.Notes, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	asset, _ := repository.FindAssetByID(c.Param("id"))
	c.JSON(http.StatusOK, asset)
}

func TransferAsset(c *gin.Context) {
	var req dto.TransferAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := middleware.GetUserID(c)
	if err := service.TransferAsset(c.Param("id"), req.ToUserID, req.LocationID, req.Notes, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	asset, _ := repository.FindAssetByID(c.Param("id"))
	c.JSON(http.StatusOK, asset)
}

func DisposeAsset(c *gin.Context) {
	var req dto.DisposeAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := middleware.GetUserID(c)
	if err := service.DisposeAsset(c.Param("id"), req.Reason, req.ResidualValue, req.Notes, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "asset disposed"})
}

func GetAssetHistory(c *gin.Context) {
	history, err := repository.GetAssetHistory(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

func GetMyAssets(c *gin.Context) {
	userID := middleware.GetUserID(c)
	assets, err := repository.ListAssetsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, assets)
}

func CreateAssetModel(c *gin.Context) {
	var req dto.CreateAssetModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m := &model.AssetModel{
		ID:                      uuid.New().String(),
		Brand:                   req.Brand,
		ModelName:               req.ModelName,
		AssetModelType:          model.AssetType(req.AssetType),
		Category:                model.AssetCategory(req.Category),
		DefaultUnitPrice:        req.DefaultUnitPrice,
		DefaultUsefulLifeMonths: req.DefaultUsefulLifeMonths,
		DepreciationMethod:      model.DepreciationMethod(req.DepreciationMethod),
	}
	if m.DepreciationMethod == "" {
		m.DepreciationMethod = model.StraightLine
	}
	if err := repository.CreateAssetModel(m); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

func ListAssetModels(c *gin.Context) {
	models, err := repository.ListAssetModels()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models)
}

func GetAssetModel(c *gin.Context) {
	m, err := repository.FindAssetModelByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, m)
}

func CreateLocation(c *gin.Context) {
	var req dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loc := &model.Location{
		ID: uuid.New().String(), Name: req.Name,
		Building: req.Building, Floor: req.Floor, Room: req.Room, Capacity: req.Capacity,
	}
	if err := repository.CreateLocation(loc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, loc)
}

func ListLocations(c *gin.Context) {
	locs, err := repository.ListLocations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, locs)
}

func GetLocation(c *gin.Context) {
	l, err := repository.FindLocationByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, l)
}

func UpdateLocation(c *gin.Context) {
	l, err := repository.FindLocationByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var req dto.CreateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name != "" {
		l.Name = req.Name
	}
	if req.Building != "" {
		l.Building = req.Building
	}
	if req.Floor != "" {
		l.Floor = req.Floor
	}
	if req.Room != "" {
		l.Room = req.Room
	}
	if req.Capacity > 0 {
		l.Capacity = req.Capacity
	}
	repository.UpdateLocation(l)
	c.JSON(http.StatusOK, l)
}

func DeleteLocation(c *gin.Context) {
	if err := repository.DeleteLocation(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func CreateDepartment(c *gin.Context) {
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	d := &model.Department{
		ID: uuid.New().String(), Name: req.Name, Description: req.Description,
	}
	if err := repository.CreateDepartment(d); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, d)
}

func ListDepartments(c *gin.Context) {
	deps, err := repository.ListDepartments()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, deps)
}

func GetDepartment(c *gin.Context) {
	d, err := repository.FindDepartmentByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, d)
}

func UpdateDepartment(c *gin.Context) {
	d, err := repository.FindDepartmentByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	var req dto.CreateDepartmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Name != "" {
		d.Name = req.Name
	}
	if req.Description != "" {
		d.Description = req.Description
	}
	repository.UpdateDepartment(d)
	c.JSON(http.StatusOK, d)
}

func DeleteDepartment(c *gin.Context) {
	if err := repository.DeleteDepartment(c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}
