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
