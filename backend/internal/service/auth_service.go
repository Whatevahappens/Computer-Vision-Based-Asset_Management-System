package service

import (
	"backend/internal/repository"
)

func Login(email, password string) string {
	user := repository.FindUserByEmail(email)
	return "fake-token-for-" + user.Email
}
