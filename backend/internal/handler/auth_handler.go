package handler

import (
	"backend/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	token := service.Login("test", "123")
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func RequestPasswordReset(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "reset email sent"})
}
