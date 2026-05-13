package repository

import (
	"backend/internal/database"
	"backend/internal/model"
)

func CreateAuditSession(s *model.AuditSession) error {
	return database.DB.Create(s).Error
}

func UpdateAuditSession(s *model.AuditSession) error {
	return database.DB.Save(s).Error
}

func FindAuditSessionByID(id string) (*model.AuditSession, error) {
	var s model.AuditSession
	err := database.DB.Preload("Location").Preload("Performer").
		Preload("Findings").Preload("Findings.Evidence").Preload("Summaries").
		First(&s, "id = ?", id).Error
	return &s, err
}

func ListAuditSessions(offset, limit int) ([]model.AuditSession, int64, error) {
	var sessions []model.AuditSession
	var total int64
	database.DB.Model(&model.AuditSession{}).Count(&total)
	err := database.DB.Preload("Location").Preload("Performer").Offset(offset).Limit(limit).Order("started_at DESC").Find(&sessions).Error
	return sessions, total, err
}

func CreateAuditFinding(f *model.AuditFinding) error {
	return database.DB.Create(f).Error
}

func CreateAuditEvidence(e *model.AuditEvidence) error {
	return database.DB.Create(e).Error
}

func CreateAuditSummary(s *model.AuditSummary) error {
	return database.DB.Create(s).Error
}

func CountAuditSessions() (int64, error) {
	var count int64
	err := database.DB.Model(&model.AuditSession{}).Count(&count).Error
	return count, err
}

func CreateNotification(n *model.Notification) error {
	return database.DB.Create(n).Error
}

func ListNotifications(userID string) ([]model.Notification, error) {
	var notifs []model.Notification
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&notifs).Error
	return notifs, err
}

func MarkNotificationRead(id, userID string) error {
	return database.DB.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Update("is_read", true).Error
}

func MarkAllNotificationsRead(userID string) error {
	return database.DB.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error
}
