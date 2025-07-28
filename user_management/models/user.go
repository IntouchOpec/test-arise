package models

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user in the system
type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"not null;size:100" validate:"required,min=2,max=100"`
	Email     string         `json:"email" gorm:"uniqueIndex;not null;size:100" validate:"required,email"`
	Age       int            `json:"age" gorm:"not null" validate:"required,min=0,max=150"`
	Phone     string         `json:"phone" gorm:"size:20" validate:"omitempty,min=10,max=20"`
	Address   string         `json:"address" gorm:"size:255" validate:"omitempty,max=255"`
	IsActive  bool           `json:"is_active" gorm:"default:true"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// UserRequest represents the request payload for creating/updating users
type UserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Age      int    `json:"age" validate:"required,min=0,max=150"`
	Phone    string `json:"phone" validate:"omitempty,min=10,max=20"`
	Address  string `json:"address" validate:"omitempty,max=255"`
	IsActive *bool  `json:"is_active,omitempty"`
}

// UserResponse represents the response payload for user operations
type UserResponse struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts User model to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Age:       u.Age,
		Phone:     u.Phone,
		Address:   u.Address,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// UpdateFromRequest updates user fields from request
func (u *User) UpdateFromRequest(req UserRequest) {
	u.Name = req.Name
	u.Email = req.Email
	u.Age = req.Age
	u.Phone = req.Phone
	u.Address = req.Address
	if req.IsActive != nil {
		u.IsActive = *req.IsActive
	}
}

// TableName specifies the table name for GORM
func (User) TableName() string {
	return "users"
}
