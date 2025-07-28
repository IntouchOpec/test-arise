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
			}

			mockRepo.AssertExpectations(t)
		})
	}
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
		name       string
		page       int
		pageSize   int
		mockUsers  []models.User
		mockCount  int64
		mockError  error
		countError error
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockRepo := new(MockUserRepository)
			userService := service.NewUserService(mockRepo, nil)

			// Calculate expected offset
			offset := (tt.page - 1) * tt.pageSize

			// Mock setup
			mockRepo.On("GetAll", offset, tt.pageSize).Return(tt.mockUsers, tt.mockError)
			mockRepo.On("Count").Return(tt.mockCount, tt.countError)

			// Execute
			result, total, err := userService.GetAllUsers(tt.page, tt.pageSize)

			// Assertions
			assert.NoError(t, err)
			assert.Equal(t, len(tt.mockUsers), len(result))
			assert.Equal(t, tt.mockCount, total)

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
			if !tt.expectedError {
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
