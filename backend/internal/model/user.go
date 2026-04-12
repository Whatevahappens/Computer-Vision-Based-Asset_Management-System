package model

type User struct {
	Firstname string
	Lastname  string
	Email     string
	Username  string
	Password  string
	Phone     string
	Status    UserStatus
	UserRole  Role
}
