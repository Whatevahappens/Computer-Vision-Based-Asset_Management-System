package handler

import (
	"backend/internal/dto"
	"backend/internal/middleware"
	"backend/internal/repository"
	"backend/internal/service"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartAudit(c *gin.Context) {
	var req dto.StartAuditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := middleware.GetUserID(c)
	session, err := service.StartAudit(req.LocationID, req.Notes, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, session)
}

func RunCVAudit(c *gin.Context) {
	sessionID := c.Param("id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file required"})
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read file"})
		return
	}

	userID := middleware.GetUserID(c)
	result, err := service.RunCVAudit(sessionID, imageData, header.Filename, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func GetAuditSession(c *gin.Context) {
	session, err := repository.FindAuditSessionByID(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "session not found"})
		return
	}
	c.JSON(http.StatusOK, session)
}

func ListAuditSessions(c *gin.Context) {
	var q dto.PaginationQuery
	c.ShouldBindQuery(&q)
	if q.Limit <= 0 {
		q.Limit = 20
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	sessions, total, err := repository.ListAuditSessions(q.Offset(), q.Limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	pages := int(total) / q.Limit
	if int(total)%q.Limit != 0 {
		pages++
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data: sessions, Total: total, Page: q.Page, Limit: q.Limit, TotalPages: pages,
	})
}

func ListNotifications(c *gin.Context) {
	userID := middleware.GetUserID(c)
	notifs, err := repository.ListNotifications(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, notifs)
}

func MarkNotificationRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := repository.MarkNotificationRead(c.Param("id"), userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "marked as read"})
}

func MarkAllNotificationsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	if err := repository.MarkAllNotificationsRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all mark as read"})
}
