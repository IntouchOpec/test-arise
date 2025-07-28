package tests

import (
	"errors"
	"testing"

	"github.com/IntouchOpec/user_management/models"
	"github.com/IntouchOpec/user_management/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetAll(offset, limit int) ([]models.User, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}

func TestUserService_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		request        models.UserRequest
		existingUser   *models.User
		existingErr    error
		createErr      error
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "successful creation",
			request: models.UserRequest{
				Name:    "John Doe",
				Email:   "john@example.com",
				Age:     30,
				Phone:   "1234567890",
				Address: "123 Main St",
			},
			existingUser:  nil,
			existingErr:   errors.New("user not found"),
			createErr:     nil,
			expectedError: false,
		},
		{
			name: "successful creation with IsActive set to false",
			request: models.UserRequest{
				Name:     "Jane Doe",
				Email:    "jane@example.com",
				Age:      25,
				IsActive: boolPtr(false),
			},
			existingUser:  nil,
			existingErr:   errors.New("user not found"),
			createErr:     nil,
			expectedError: false,
		},
		{
			name: "email already exists",
			request: models.UserRequest{
				Name:  "Jane Doe",
				Email: "existing@example.com",
				Age:   25,
			},
			existingUser: &models.User{
				ID:    1,
				Email: "existing@example.com",
			},
			existingErr:    nil,
			createErr:      nil,
			expectedError:  true,
			expectedErrMsg: "user with email existing@example.com already exists",
		},
		{
			name: "create error",
			request: models.UserRequest{
				Name:  "Bob Smith",
				Email: "bob@example.com",
				Age:   35,
			},
			existingUser:   nil,
			existingErr:    errors.New("user not found"),
			createErr:      errors.New("database error"),
			expectedError:  true,
			expectedErrMsg: "failed to create user: database error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			userService := service.NewUserService(mockRepo, nil)

			// Mock setup
			mockRepo.On("GetByEmail", tt.request.Email).Return(tt.existingUser, tt.existingErr)
			if tt.existingUser == nil {
				mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(tt.createErr)
			}

			// Execute
			result, err := userService.CreateUser(tt.request)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Name, result.Name)
				assert.Equal(t, tt.request.Email, result.Email)
				assert.Equal(t, tt.request.Age, result.Age)

				// Check IsActive handling
				if tt.request.IsActive != nil {
					assert.Equal(t, *tt.request.IsActive, result.IsActive)
				} else {
					assert.True(t, result.IsActive) // Default is true
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

// Helper function to create bool pointer
func boolPtr(b bool) *bool {
	return &b
}

func TestUserService_GetUserByID(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		mockUser      *models.User
		mockError     error
		expectedError bool
	}{
		{
			name:   "successful retrieval",
			userID: 1,
			mockUser: &models.User{
				ID:       1,
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      30,
				IsActive: true,
			},
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "user not found",
			userID:        999,
			mockUser:      nil,
			mockError:     errors.New("user not found"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			userService := service.NewUserService(mockRepo, nil)

			// Mock setup
			mockRepo.On("GetByID", tt.userID).Return(tt.mockUser, tt.mockError)

			// Execute
			result, err := userService.GetUserByID(tt.userID)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.mockUser.ID, result.ID)
				assert.Equal(t, tt.mockUser.Name, result.Name)
				assert.Equal(t, tt.mockUser.Email, result.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_GetAllUsers(t *testing.T) {
	tests := []struct {
		name              string
		page              int
		pageSize          int
		mockUsers         []models.User
		mockCount         int64
		mockError         error
		countError        error
		expectGetAllError bool
		expectCountError  bool
	}{
		{
			name:     "successful retrieval",
			page:     1,
			pageSize: 10,
			mockUsers: []models.User{
				{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30, IsActive: true},
				{ID: 2, Name: "Jane Doe", Email: "jane@example.com", Age: 25, IsActive: true},
			},
			mockCount:  2,
			mockError:  nil,
			countError: nil,
		},
		{
			name:       "empty result",
			page:       1,
			pageSize:   10,
			mockUsers:  []models.User{},
			mockCount:  0,
			mockError:  nil,
			countError: nil,
		},
		{
			name:     "page less than 1",
			page:     0,
			pageSize: 10,
			mockUsers: []models.User{
				{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30, IsActive: true},
			},
			mockCount:  1,
			mockError:  nil,
			countError: nil,
		},
		{
			name:     "pageSize less than 1",
			page:     1,
			pageSize: 0,
			mockUsers: []models.User{
				{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30, IsActive: true},
			},
			mockCount:  1,
			mockError:  nil,
			countError: nil,
		},
		{
			name:     "pageSize greater than 100",
			page:     1,
			pageSize: 200,
			mockUsers: []models.User{
				{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30, IsActive: true},
			},
			mockCount:  1,
			mockError:  nil,
			countError: nil,
		},
		{
			name:              "GetAll error",
			page:              1,
			pageSize:          10,
			mockUsers:         []models.User{},
			mockCount:         0,
			mockError:         errors.New("database connection failed"),
			countError:        nil,
			expectGetAllError: true,
		},
		{
			name:             "Count error",
			page:             1,
			pageSize:         10,
			mockUsers:        []models.User{{ID: 1, Name: "John", Email: "john@example.com", Age: 30}},
			mockCount:        0,
			mockError:        nil,
			countError:       errors.New("count failed"),
			expectCountError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			userService := service.NewUserService(mockRepo, nil)

			// Calculate expected offset and page size
			expectedPage := tt.page
			expectedPageSize := tt.pageSize

			if expectedPage < 1 {
				expectedPage = 1
			}
			if expectedPageSize < 1 || expectedPageSize > 100 {
				expectedPageSize = 10
			}

			offset := (expectedPage - 1) * expectedPageSize

			// Mock setup
			mockRepo.On("GetAll", offset, expectedPageSize).Return(tt.mockUsers, tt.mockError)
			if !tt.expectGetAllError {
				mockRepo.On("Count").Return(tt.mockCount, tt.countError)
			}

			// Execute
			result, total, err := userService.GetAllUsers(tt.page, tt.pageSize)

			// Assertions
			if tt.expectGetAllError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, int64(0), total)
				assert.Contains(t, err.Error(), "failed to get users:")
			} else if tt.expectCountError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Equal(t, int64(0), total)
				assert.Contains(t, err.Error(), "failed to count users:")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, len(tt.mockUsers), len(result))
				assert.Equal(t, tt.mockCount, total)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         uint
		request        models.UserRequest
		existingUser   *models.User
		existingErr    error
		emailUser      *models.User
		emailErr       error
		updateErr      error
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name:   "successful update",
			userID: 1,
			request: models.UserRequest{
				Name:    "John Updated",
				Email:   "john.updated@example.com",
				Age:     31,
				Phone:   "1234567890",
				Address: "456 New St",
			},
			existingUser: &models.User{
				ID:       1,
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      30,
				IsActive: true,
			},
			existingErr:   nil,
			emailUser:     nil,
			emailErr:      errors.New("user not found"),
			updateErr:     nil,
			expectedError: false,
		},
		{
			name:   "update with same email (no check needed)",
			userID: 1,
			request: models.UserRequest{
				Name:    "John Updated",
				Email:   "john@example.com", // Same email
				Age:     31,
				Phone:   "1234567890",
				Address: "456 New St",
			},
			existingUser: &models.User{
				ID:       1,
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      30,
				IsActive: true,
			},
			existingErr:   nil,
			updateErr:     nil,
			expectedError: false,
		},
		{
			name:   "user not found for update",
			userID: 999,
			request: models.UserRequest{
				Name:  "Non Existent",
				Email: "nonexistent@example.com",
				Age:   30,
			},
			existingUser:   nil,
			existingErr:    errors.New("user not found"),
			expectedError:  true,
			expectedErrMsg: "user not found",
		},
		{
			name:   "email already exists for different user",
			userID: 1,
			request: models.UserRequest{
				Name:  "John Doe",
				Email: "existing@example.com",
				Age:   30,
			},
			existingUser: &models.User{
				ID:    1,
				Email: "john@example.com",
			},
			existingErr: nil,
			emailUser: &models.User{
				ID:    2,
				Email: "existing@example.com",
			},
			emailErr:       nil,
			updateErr:      nil,
			expectedError:  true,
			expectedErrMsg: "user with email existing@example.com already exists",
		},
		{
			name:   "update database error",
			userID: 1,
			request: models.UserRequest{
				Name:  "John Updated",
				Email: "john.updated@example.com",
				Age:   31,
			},
			existingUser: &models.User{
				ID:       1,
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      30,
				IsActive: true,
			},
			existingErr:    nil,
			emailUser:      nil,
			emailErr:       errors.New("user not found"),
			updateErr:      errors.New("database update failed"),
			expectedError:  true,
			expectedErrMsg: "failed to update user:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			userService := service.NewUserService(mockRepo, nil)

			// Mock setup
			mockRepo.On("GetByID", tt.userID).Return(tt.existingUser, tt.existingErr)
			if tt.existingUser != nil && tt.existingUser.Email != tt.request.Email {
				mockRepo.On("GetByEmail", tt.request.Email).Return(tt.emailUser, tt.emailErr)
			}
			// Add Update expectation for cases where we reach the update step
			if tt.existingUser != nil && (tt.emailUser == nil || tt.emailErr != nil) {
				mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(tt.updateErr)
			}

			// Execute
			result, err := userService.UpdateUser(tt.userID, tt.request)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
				assert.Contains(t, err.Error(), tt.expectedErrMsg)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)
				assert.Equal(t, tt.request.Name, result.Name)
				assert.Equal(t, tt.request.Email, result.Email)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUserService_DeleteUser(t *testing.T) {
	tests := []struct {
		name          string
		userID        uint
		mockError     error
		expectedError bool
	}{
		{
			name:          "successful deletion",
			userID:        1,
			mockError:     nil,
			expectedError: false,
		},
		{
			name:          "delete error",
			userID:        999,
			mockError:     errors.New("user not found"),
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			userService := service.NewUserService(mockRepo, nil)

			// Mock setup
			mockRepo.On("Delete", tt.userID).Return(tt.mockError)

			// Execute
			err := userService.DeleteUser(tt.userID)

			// Assertions
			if tt.expectedError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "failed to delete user:")
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
