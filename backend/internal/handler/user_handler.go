package handler

import (
	"backend/internal/dto"
	"backend/internal/model"
	"backend/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check duplicate email
	if existing, _ := repository.FindUserByEmail(req.Email); existing != nil && existing.ID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}
	// Check duplicate username
	if existing, _ := repository.FindUserByUsername(req.Username); existing != nil && existing.ID != "" {
		c.JSON(http.StatusConflict, gin.H{"error": "username already exists"})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &model.User{
		ID:           uuid.New().String(),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		Username:     req.Username,
		PasswordHash: string(hash),
		Phone:        req.Phone,
		Status:       model.UserActive,
		Role:         model.Role(req.Role),
	}
	if req.DepartmentID != "" {
		user.DepartmentID = &req.DepartmentID
	}

	if err := repository.CreateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, toUserResponse(user))
}

func ListUsers(c *gin.Context) {
	var q dto.PaginationQuery
	c.ShouldBindQuery(&q)
	if q.Limit <= 0 {
		q.Limit = 20
	}
	if q.Page <= 0 {
		q.Page = 1
	}

	users, total, err := repository.ListUsers(q.Offset(), q.Limit, q.Search)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var resp []dto.UserResponse
	for _, u := range users {
		resp = append(resp, toUserResponse(&u))
	}

	pages := int(total) / q.Limit
	if int(total)%q.Limit != 0 {
		pages++
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Data: resp, Total: total, Page: q.Page, Limit: q.Limit, TotalPages: pages,
	})
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	user, err := repository.FindUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	user, err := repository.FindUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Status != "" {
		user.Status = model.UserStatus(req.Status)
	}
	if req.DepartmentID != "" {
		user.DepartmentID = &req.DepartmentID
	}

	if err := repository.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}
	c.JSON(http.StatusOK, toUserResponse(user))
}

func DeactivateUser(c *gin.Context) {
	id := c.Param("id")
	if err := repository.DeactivateUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to deactivate user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deactivated"})
}

// helpers
func toUserResponse(u *model.User) dto.UserResponse {
	return dto.UserResponse{
		ID: u.ID, FirstName: u.FirstName, LastName: u.LastName,
		Email: u.Email, Username: u.Username, Phone: u.Phone,
		Status: string(u.Status), Role: string(u.Role), DepartmentID: u.DepartmentID,
	}
}

func getUserResponse(id string) (*dto.UserResponse, error) {
	u, err := repository.FindUserByID(id)
	if err != nil {
		return nil, err
	}
	r := toUserResponse(u)
	return &r, nil
}
