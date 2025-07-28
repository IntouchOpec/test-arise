package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/IntouchOpec/user_management/controllers"
	"github.com/IntouchOpec/user_management/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockUserService is a mock implementation of UserService
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(req models.UserRequest) (*models.UserResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByID(id uint) (*models.UserResponse, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) GetAllUsers(page, pageSize int) ([]models.UserResponse, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0).([]models.UserResponse), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserService) UpdateUser(id uint, req models.UserRequest) (*models.UserResponse, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserResponse), args.Error(1)
}

func (m *MockUserService) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

func TestUserController_CreateUser(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    models.UserRequest
		mockReturn     *models.UserResponse
		mockError      error
		expectedStatus int
		expectedError  bool
	}{
		{
			name: "successful creation",
			requestBody: models.UserRequest{
				Name:    "John Doe",
				Email:   "john@example.com",
				Age:     30,
				Phone:   "1234567890",
				Address: "123 Main St",
			},
			mockReturn: &models.UserResponse{
				ID:       1,
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      30,
				Phone:    "1234567890",
				Address:  "123 Main St",
				IsActive: true,
			},
			mockError:      nil,
			expectedStatus: http.StatusCreated,
			expectedError:  false,
		},
		{
			name: "email already exists",
			requestBody: models.UserRequest{
				Name:  "Jane Doe",
				Email: "existing@example.com",
				Age:   25,
			},
			mockReturn:     nil,
			mockError:      errors.New("user with email existing@example.com already exists"),
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			controller := controllers.NewUserController(mockService)
			router := setupTestRouter()

			// Mock setup
			mockService.On("CreateUser", tt.requestBody).Return(tt.mockReturn, tt.mockError)

			// Route setup
			router.POST("/users", controller.CreateUser)

			// Create request
			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.Contains(t, response, "error")
			} else {
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "data")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserController_GetUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockReturn     *models.UserResponse
		mockError      error
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "successful retrieval",
			userID: "1",
			mockReturn: &models.UserResponse{
				ID:       1,
				Name:     "John Doe",
				Email:    "john@example.com",
				Age:      30,
				IsActive: true,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "user not found",
			userID:         "999",
			mockReturn:     nil,
			mockError:      errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
		{
			name:           "invalid user ID",
			userID:         "invalid",
			mockReturn:     nil,
			mockError:      nil,
			expectedStatus: http.StatusBadRequest,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			controller := controllers.NewUserController(mockService)
			router := setupTestRouter()

			// Mock setup (only for valid IDs)
			if tt.userID != "invalid" {
				userID, _ := strconv.ParseUint(tt.userID, 10, 32)
				mockService.On("GetUserByID", uint(userID)).Return(tt.mockReturn, tt.mockError)
			}

			// Route setup
			router.GET("/users/:id", controller.GetUser)

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/users/"+tt.userID, nil)

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.Contains(t, response, "error")
			} else {
				assert.Contains(t, response, "data")
			}

			if tt.userID != "invalid" {
				mockService.AssertExpectations(t)
			}
		})
	}
}

func TestUserController_GetUsers(t *testing.T) {
	tests := []struct {
		name           string
		queryParams    string
		mockUsers      []models.UserResponse
		mockTotal      int64
		mockError      error
		expectedStatus int
	}{
		{
			name:        "successful retrieval with default pagination",
			queryParams: "",
			mockUsers: []models.UserResponse{
				{ID: 1, Name: "John Doe", Email: "john@example.com", Age: 30, IsActive: true},
				{ID: 2, Name: "Jane Doe", Email: "jane@example.com", Age: 25, IsActive: true},
			},
			mockTotal:      2,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:        "successful retrieval with custom pagination",
			queryParams: "?page=2&page_size=5",
			mockUsers: []models.UserResponse{
				{ID: 6, Name: "User 6", Email: "user6@example.com", Age: 28, IsActive: true},
			},
			mockTotal:      10,
			mockError:      nil,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "service error",
			queryParams:    "",
			mockUsers:      nil,
			mockTotal:      0,
			mockError:      errors.New("database error"),
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			controller := controllers.NewUserController(mockService)
			router := setupTestRouter()

			// Parse expected parameters
			page, pageSize := 1, 10
			if tt.queryParams != "" {
				if tt.queryParams == "?page=2&page_size=5" {
					page, pageSize = 2, 5
				}
			}

			// Mock setup
			mockService.On("GetAllUsers", page, pageSize).Return(tt.mockUsers, tt.mockTotal, tt.mockError)

			// Route setup
			router.GET("/users", controller.GetUsers)

			// Create request
			req, _ := http.NewRequest(http.MethodGet, "/users"+tt.queryParams, nil)

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.mockError != nil {
				assert.Contains(t, response, "error")
				assert.Equal(t, tt.mockError.Error(), response["error"])
			} else {
				assert.Contains(t, response, "data")
				assert.Contains(t, response, "pagination")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserController_UpdateUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		requestBody    models.UserRequest
		mockReturn     *models.UserResponse
		mockError      error
		expectedStatus int
		expectedError  bool
	}{
		{
			name:   "successful update",
			userID: "1",
			requestBody: models.UserRequest{
				Name:    "John Updated",
				Email:   "john.updated@example.com",
				Age:     31,
				Phone:   "1234567890",
				Address: "456 New St",
			},
			mockReturn: &models.UserResponse{
				ID:       1,
				Name:     "John Updated",
				Email:    "john.updated@example.com",
				Age:      31,
				Phone:    "1234567890",
				Address:  "456 New St",
				IsActive: true,
			},
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:   "user not found",
			userID: "999",
			requestBody: models.UserRequest{
				Name:  "Non Existent",
				Email: "nonexistent@example.com",
				Age:   30,
			},
			mockReturn:     nil,
			mockError:      errors.New("user not found"),
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			controller := controllers.NewUserController(mockService)
			router := setupTestRouter()

			// Mock setup
			userID, _ := strconv.ParseUint(tt.userID, 10, 32)
			mockService.On("UpdateUser", uint(userID), tt.requestBody).Return(tt.mockReturn, tt.mockError)

			// Route setup
			router.PUT("/users/:id", controller.UpdateUser)

			// Create request
			requestBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest(http.MethodPut, "/users/"+tt.userID, bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.Contains(t, response, "error")
			} else {
				assert.Contains(t, response, "message")
				assert.Contains(t, response, "data")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserController_DeleteUser(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		mockError      error
		expectedStatus int
		expectedError  bool
	}{
		{
			name:           "successful deletion",
			userID:         "1",
			mockError:      nil,
			expectedStatus: http.StatusOK,
			expectedError:  false,
		},
		{
			name:           "user not found",
			userID:         "999",
			mockError:      errors.New("failed to delete user: user not found"),
			expectedStatus: http.StatusNotFound,
			expectedError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup
			mockService := new(MockUserService)
			controller := controllers.NewUserController(mockService)
			router := setupTestRouter()

			// Mock setup
			userID, _ := strconv.ParseUint(tt.userID, 10, 32)
			mockService.On("DeleteUser", uint(userID)).Return(tt.mockError)

			// Route setup
			router.DELETE("/users/:id", controller.DeleteUser)

			// Create request
			req, _ := http.NewRequest(http.MethodDelete, "/users/"+tt.userID, nil)

			// Perform request
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err := json.Unmarshal(w.Body.Bytes(), &response)
			assert.NoError(t, err)

			if tt.expectedError {
				assert.Contains(t, response, "error")
			} else {
				assert.Contains(t, response, "message")
			}

			mockService.AssertExpectations(t)
		})
	}
}

func TestUserController_HealthCheck(t *testing.T) {
	// Setup
	mockService := new(MockUserService)
	controller := controllers.NewUserController(mockService)
	router := setupTestRouter()

	// Route setup
	router.GET("/health", controller.HealthCheck)

	// Create request
	req, _ := http.NewRequest(http.MethodGet, "/health", nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, "healthy", response["status"])
	assert.Contains(t, response, "timestamp")
}
