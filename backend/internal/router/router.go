package router

import (
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/model"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// CORS
	r.Use(corsMiddleware())

	// Health
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	api.POST("/auth/login", handler.Login)

	auth := api.Group("")
	auth.Use(middleware.AuthRequired())
	{
		// Auth
		auth.GET("/auth/me", handler.GetMe)
		auth.PUT("/auth/password", handler.ChangePassword)

		auth.GET("/dashboard", handler.GetDashboard)

		auth.GET("/my-assets", handler.GetMyAssets)

		auth.GET("/notifications", handler.ListNotifications)
		auth.GET("/notifications/:id/read", handler.MarkNotificationRead)
		auth.GET("/notifications/read-all", handler.MarkAllNotificationsRead)

		assets := auth.Group("/assets")
		assets.Use(middleware.RequireRole(string(model.AssetCustodian)))
		{
			assets.POST("", handler.CreateAsset)
			assets.GET("", handler.ListAssets)
			assets.GET("/:id", handler.GetAsset)
			assets.PUT("/:id", handler.UpdateAsset)
			assets.DELETE("/:id", handler.DeleteAsset)
			assets.POST("/:id/assign", handler.AssignAsset)
			assets.POST("/:id/transfer", handler.TransferAsset)
			assets.POST("/:id/dispose", handler.DisposeAsset)
			assets.GET("/:id/history", handler.GetAssetHistory)
		}

		// Asset Models
		assetModels := auth.Group("/asset-models")
		assetModels.Use(middleware.RequireRole(string(model.AssetCustodian)))
		{
			assetModels.POST("", handler.CreateAssetModel)
			assetModels.GET("", handler.ListAssetModels)
			assetModels.GET("/:id", handler.GetAssetModel)
		}

		// Audit (Custodian, Admin)
		audits := auth.Group("/audits")
		audits.Use(middleware.RequireRole(string(model.Admin), string(model.AssetCustodian)))
		{
			audits.POST("", handler.StartAudit)
			audits.GET("", handler.ListAuditSessions)
			audits.GET("/:id", handler.GetAuditSession)
			audits.POST("/:id/cv", handler.RunCVAudit)
		}

		// Depreciation (Accountant, Admin)
		depreciation := auth.Group("/depreciation")
		depreciation.Use(middleware.RequireRole(string(model.Accountant), string(model.Admin)))
		{
			depreciation.POST("/calculate", handler.CalculateDepreciation)
			depreciation.POST("/revalue", handler.RevalueAsset)
		}

		// Reports (Custodian, Accountant, Admin)
		reports := auth.Group("/reports")
		reports.Use(middleware.RequireRole(string(model.AssetCustodian), string(model.Admin), string(model.Accountant)))
		{
			reports.POST("/generate", handler.GenerateReport)
			reports.GET("/download/:filename", handler.DownloadReport)
		}

		// Locations (Custodian, Admin)
		locations := auth.Group("/locations")
		locations.Use(middleware.RequireRole(string(model.AssetCustodian), string(model.Admin)))
		{
			locations.POST("", handler.CreateLocation)
			locations.GET("", handler.ListLocations)
			locations.GET("/:id", handler.GetLocation)
			locations.PUT("/:id", handler.UpdateLocation)
			locations.DELETE("/:id", handler.DeleteLocation)
		}

		// Departments (Admin)
		departments := auth.Group("/departments")
		departments.Use(middleware.RequireRole(string(model.Admin)))
		{
			departments.POST("", handler.CreateDepartment)
			departments.GET("", handler.ListDepartments)
			departments.GET("/:id", handler.GetDepartment)
			departments.PUT("/:id", handler.UpdateDepartment)
			departments.DELETE("/:id", handler.DeleteDepartment)
		}

		// Users (Admin)
		users := auth.Group("/users")
		users.Use(middleware.RequireRole(string(model.Admin)))
		{
			users.POST("", handler.CreateUser)
			users.GET("", handler.ListUsers)
			users.GET("/:id", handler.GetUser)
			users.PUT("/:id", handler.UpdateUser)
			users.PUT("/:id/deactivate", handler.DeactivateUser)
		}
	}

	// Allow read-only list endpoints for all authenticated users
	readOnly := api.Group("")
	readOnly.Use(middleware.AuthRequired())
	{
		readOnly.GET("/locations-list", handler.ListLocations)
		readOnly.GET("/departments-list", handler.ListDepartments)
		readOnly.GET("/asset-models-list", handler.ListAssetModels)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
