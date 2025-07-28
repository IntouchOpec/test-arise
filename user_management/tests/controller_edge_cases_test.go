package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/IntouchOpec/user_management/controllers"
	"github.com/IntouchOpec/user_management/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserController_CreateUser_InvalidJSON(t *testing.T) {
	// Setup
	mockService := new(MockUserService)
	controller := controllers.NewUserController(mockService)
	router := gin.New()
	router.POST("/users", controller.CreateUser)

	// Create request with invalid JSON
	req, _ := http.NewRequest(http.MethodPost, "/users", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestUserController_GetUser_InvalidID(t *testing.T) {
	// Setup
	mockService := new(MockUserService)
	controller := controllers.NewUserController(mockService)
	router := gin.New()
	router.GET("/users/:id", controller.GetUser)

	// Create request with invalid ID
	req, _ := http.NewRequest(http.MethodGet, "/users/invalid", nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid user ID", response["error"])
}

func TestUserController_UpdateUser_InvalidID(t *testing.T) {
	// Setup
	mockService := new(MockUserService)
	controller := controllers.NewUserController(mockService)
	router := gin.New()
	router.PUT("/users/:id", controller.UpdateUser)

	// Create request with invalid ID
	requestBody := map[string]interface{}{
		"name":  "Updated Name",
		"email": "updated@example.com",
		"age":   30,
	}
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPut, "/users/invalid", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid user ID", response["error"])
}

func TestUserController_UpdateUser_InvalidJSON(t *testing.T) {
	// Setup
	mockService := new(MockUserService)
	controller := controllers.NewUserController(mockService)
	router := gin.New()
	router.PUT("/users/:id", controller.UpdateUser)

	// Create request with invalid JSON
	req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBufferString("invalid json"))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid request body", response["error"])
}

func TestUserController_DeleteUser_InvalidID(t *testing.T) {
	// Setup
	mockService := new(MockUserService)
	controller := controllers.NewUserController(mockService)
	router := gin.New()
	router.DELETE("/users/:id", controller.DeleteUser)

	// Create request with invalid ID
	req, _ := http.NewRequest(http.MethodDelete, "/users/invalid", nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "Invalid user ID", response["error"])
}

func TestUserController_DeleteUser_InternalServerError(t *testing.T) {
	// Setup
	mockService := new(MockUserService)
	controller := controllers.NewUserController(mockService)
	router := gin.New()
	router.DELETE("/users/:id", controller.DeleteUser)

	// Mock setup - simulate internal server error
	mockService.On("DeleteUser", uint(1)).Return(errors.New("internal server error"))

	// Create request
	req, _ := http.NewRequest(http.MethodDelete, "/users/1", nil)

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")

	mockService.AssertExpectations(t)
}

func TestUserController_UpdateUser_UpdateError(t *testing.T) {
	// Setup
	mockService := new(MockUserService)
	controller := controllers.NewUserController(mockService)
	router := gin.New()
	router.PUT("/users/:id", controller.UpdateUser)

	requestBody := models.UserRequest{
		Name:  "Updated Name",
		Email: "updated@example.com",
		Age:   30,
	}

	// Mock setup - simulate update error (not user not found)
	mockService.On("UpdateUser", uint(1), requestBody).Return(nil, errors.New("database connection failed"))

	// Create request
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Contains(t, response, "error")
	assert.Equal(t, "database connection failed", response["error"])

	mockService.AssertExpectations(t)
}
