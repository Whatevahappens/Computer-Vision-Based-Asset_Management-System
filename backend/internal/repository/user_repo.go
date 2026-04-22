package repository

import (
	"backend/internal/database"
	"backend/internal/model"
)

func CreateUser(user *model.User) error {
	return database.DB.Create(user).Error
}

func FindUserByID(id string) (*model.User, error) {
	var user model.User
	err := database.DB.Preload("Department").First(&user, "id = ?", id).Error
	return &user, err
}

func FindUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, "username = ?", username).Error
	return &user, err
}

func FindUserByEmail(email string) (*model.User, error) {
	var user model.User
	err := database.DB.First(&user, "email = ?", email).Error
	return &user, err
}

func ListUsers(offset, limit int, search string) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	q := database.DB.Model(&model.User{})
	if search != "" {
		q = q.Where("first_name ILIKE ? OR last_name ILIKE ? OR email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	q.Count(&total)
	err := q.Preload("Department").Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func UpdateUser(user *model.User) error {
	return database.DB.Save(user).Error
}

func DeactivateUser(id string) error {
	return database.DB.Model(&model.User{}).Where("id = ?", id).Update("status", model.UserInactive).Error
}
