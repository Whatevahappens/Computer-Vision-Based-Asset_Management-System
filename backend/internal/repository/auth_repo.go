package repository

import (
	"backend/internal/model"
)

func FindUserByEmail(email string) model.User {
	return model.User{Email: email}
}
