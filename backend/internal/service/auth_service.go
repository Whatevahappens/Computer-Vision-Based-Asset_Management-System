package service

import (
	"backend/internal/dto"
	"backend/internal/repository"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	RequestPasswordReset(req dto.ResetPasswordRequest) error
	ConfirmPasswordReset(req dto.ResetPasswordConfirm) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token := "jwt-token-for-" + user.Email

	return &dto.LoginResponse{
		Token: token,
		Email: user.Email,
	}, nil
}

func (s *authService) RequestPasswordReset(req dto.ResetPasswordRequest) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil
	}

	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return errors.New("failed to generate token")
	}
	token := hex.EncodeToString(bytes)

	expiry := time.Now().Add(1 * time.Hour)
	user.ResetToken = token
	user.ResetExpiresAt = &expiry

	if err := s.userRepo.Update(user); err != nil {
		return errors.New("failed to save reset token")
	}

	return nil
}

func (s *authService) ConfirmPasswordReset(req model.ResetPasswordConfirm) {
	user, err := s.userRepo.FindByResetToken(req.Token)
	if err != nil {
		return errors.New("invalid or expired token")
	}
	if user.ResetExpiresAt == nil || user.ResetExpiresAt.Before((time.Now())) {
		return errors.New("invalid or expired token")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashed)
	user.ResetToken = ""
	user.ResetExpiresAt = nil

	return s.userRepo.Update(user)
}
