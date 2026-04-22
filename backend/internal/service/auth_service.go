package service

import (
	"backend/internal/dto"
	"backend/internal/middleware"
	"backend/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func Login(username, password string, expiryHours int) (string, dto.UserResponse, error) {
	user, err := repository.FindUserByUsername(username)
	if err != nil {
		return "", dto.UserResponse{}, errors.New("invalid credentials")
	}
	if user.Status == "INACTIVE" || user.Status == "SUSPENDED" {
		return "", dto.UserResponse{}, errors.New("account is disabled")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", dto.UserResponse{}, errors.New("invalid credentials")
	}
	token, err := middleware.GenerateToken(user.ID, user.Username, string(user.Role), expiryHours)
	if err != nil {
		return "", dto.UserResponse{}, errors.New("failed to generate token")
	}
	resp := dto.UserResponse{
		ID: user.ID, FirstName: user.FirstName, LastName: user.LastName,
		Email: user.Email, Username: user.Username, Phone: user.Phone,
		Status: string(user.Status), Role: string(user.Role), DepartmentID: user.DepartmentID,
	}
	return token, resp, nil
}

func ChangePassword(userID, currentPwd, newPwd string) error {
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(currentPwd)); err != nil {
		return errors.New("current password is incorrect")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPwd), 12)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hash)
	return repository.UpdateUser(user)
}
