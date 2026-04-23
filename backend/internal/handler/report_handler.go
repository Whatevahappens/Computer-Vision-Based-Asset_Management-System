package handler

import (
	"backend/internal/dto"
	"backend/internal/middleware"
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetDashboard(c *gin.Context) {
	stats, err := service.GetDashboardStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stats)
}

func GenerateReport(c *gin.Context) {
	var req dto.ReportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	filePath, err := service.GenerateReport(req.ReportType, req.Format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "report generated",
		"filePath": filePath,
		"download": "/api/v1/reports/download/" + filePath.Base(filePath),
	})
}

func DownloadReport(c *gin.Context) {
	filename := c.Param("filename")
	filePath := "/tmp/" + filename
	c.File(filePath)
}

func CalculateDepreciation(c *gin.Context) {
	var req dto.CalculateDepreciationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	result, err := service.CalculateDepreciation(req.AssetID, req.Method, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func RevalueAsset(c *gin.Context) {
	var req dto.RevalueAssetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := middleware.GetUserID(c)
	if err := service.RevalueAsset(req.AssetID, req.NewValue, req.Reason, userID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "asset revalued"})
}
