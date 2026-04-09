package router

import (
	"backend/internal/handler"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", handler.Login)
			auth.POST("/password-reset", handler.RequestPasswordReset)
		}
	}
}
