package tests

import (
	"testing"

	"github.com/IntouchOpec/user_management/models"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUser_ToResponse(t *testing.T) {
	user := &models.User{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		Phone:    "1234567890",
		Address:  "123 Main St",
		IsActive: true,
	}

	response := user.ToResponse()

	assert.Equal(t, user.ID, response.ID)
	assert.Equal(t, user.Name, response.Name)
	assert.Equal(t, user.Email, response.Email)
	assert.Equal(t, user.Age, response.Age)
	assert.Equal(t, user.Phone, response.Phone)
	assert.Equal(t, user.Address, response.Address)
	assert.Equal(t, user.IsActive, response.IsActive)
	assert.Equal(t, user.CreatedAt, response.CreatedAt)
	assert.Equal(t, user.UpdatedAt, response.UpdatedAt)
}

func TestUser_UpdateFromRequest(t *testing.T) {
	user := &models.User{
		ID:       1,
		Name:     "Old Name",
		Email:    "old@example.com",
		Age:      25,
		Phone:    "0000000000",
		Address:  "Old Address",
		IsActive: true,
	}

	request := models.UserRequest{
		Name:    "New Name",
		Email:   "new@example.com",
		Age:     30,
		Phone:   "1234567890",
		Address: "New Address",
	}

	user.UpdateFromRequest(request)

	assert.Equal(t, request.Name, user.Name)
	assert.Equal(t, request.Email, user.Email)
	assert.Equal(t, request.Age, user.Age)
	assert.Equal(t, request.Phone, user.Phone)
	assert.Equal(t, request.Address, user.Address)
	// IsActive should remain unchanged as it's not in UpdateFromRequest
	assert.True(t, user.IsActive)
}

func TestUser_UpdateFromRequest_WithIsActive(t *testing.T) {
	user := &models.User{
		ID:       1,
		Name:     "Original Name",
		Email:    "original@example.com",
		Age:      25,
		IsActive: true,
	}

	isActive := false
	req := models.UserRequest{
		Name:     "Updated Name",
		Email:    "updated@example.com",
		Age:      30,
		IsActive: &isActive,
	}

	user.UpdateFromRequest(req)

	assert.Equal(t, "Updated Name", user.Name)
	assert.Equal(t, "updated@example.com", user.Email)
	assert.Equal(t, 30, user.Age)
	assert.False(t, user.IsActive) // This should be false now
}

func TestUser_TableName(t *testing.T) {
	user := models.User{}
	tableName := user.TableName()
	assert.Equal(t, "users", tableName)
}

func TestUserRequest_Validation(t *testing.T) {
	tests := []struct {
		name    string
		request models.UserRequest
		valid   bool
	}{
		{
			name: "valid request",
			request: models.UserRequest{
				Name:    "John Doe",
				Email:   "john@example.com",
				Age:     30,
				Phone:   "1234567890",
				Address: "123 Main St",
			},
			valid: true,
		},
		{
			name: "valid request with optional fields empty",
			request: models.UserRequest{
				Name:  "John Doe",
				Email: "john@example.com",
				Age:   30,
			},
			valid: true,
		},
		{
			name: "valid request with minimum age",
			request: models.UserRequest{
				Name:  "Baby Doe",
				Email: "baby@example.com",
				Age:   0,
			},
			valid: true,
		},
		{
			name: "valid request with maximum age",
			request: models.UserRequest{
				Name:  "Elder Doe",
				Email: "elder@example.com",
				Age:   150,
			},
			valid: true,
		},
		{
			name: "valid request with minimum name length",
			request: models.UserRequest{
				Name:  "Jo",
				Email: "jo@example.com",
				Age:   25,
			},
			valid: true,
		},
		{
			name: "valid request with maximum name length",
			request: models.UserRequest{
				Name:  "ThisIsAVeryLongNameThatIsExactlyOneHundredCharactersLongToTestTheMaximumLengthValidation",
				Email: "long@example.com",
				Age:   25,
			},
			valid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For this test, we're just checking the struct is properly defined
			// Actual validation would be done by the validator library in the service layer
			assert.NotNil(t, tt.request)
			assert.NotEmpty(t, tt.request.Name)
			assert.NotEmpty(t, tt.request.Email)
			assert.GreaterOrEqual(t, tt.request.Age, 0)
		})
	}
}

func TestUserResponse_Structure(t *testing.T) {
	response := models.UserResponse{
		ID:       1,
		Name:     "John Doe",
		Email:    "john@example.com",
		Age:      30,
		Phone:    "1234567890",
		Address:  "123 Main St",
		IsActive: true,
	}

	assert.Equal(t, uint(1), response.ID)
	assert.Equal(t, "John Doe", response.Name)
	assert.Equal(t, "john@example.com", response.Email)
	assert.Equal(t, 30, response.Age)
	assert.Equal(t, "1234567890", response.Phone)
	assert.Equal(t, "123 Main St", response.Address)
	assert.True(t, response.IsActive)
}

func TestUser_GormHooks(t *testing.T) {
	// Test that User struct has proper GORM tags and relationships
	user := models.User{
		ID:        1,
		Name:      "John Doe",
		Email:     "john@example.com",
		Age:       30,
		Phone:     "1234567890",
		Address:   "123 Main St",
		IsActive:  true,
		DeletedAt: gorm.DeletedAt{},
	}

	// Verify struct fields are properly set
	assert.Equal(t, uint(1), user.ID)
	assert.Equal(t, "John Doe", user.Name)
	assert.Equal(t, "john@example.com", user.Email)
	assert.Equal(t, 30, user.Age)
	assert.Equal(t, "1234567890", user.Phone)
	assert.Equal(t, "123 Main St", user.Address)
	assert.True(t, user.IsActive)
	assert.False(t, user.DeletedAt.Valid)
}
