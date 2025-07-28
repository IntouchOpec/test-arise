package tests

import (
	"testing"

	"github.com/IntouchOpec/user_management/controllers"
	"github.com/IntouchOpec/user_management/routes"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSetupRoutes(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(MockUserService)
	userController := controllers.NewUserController(mockService)

	// Execute
	routes.SetupRoutes(router, userController)

	// Get the registered routes
	routesList := router.Routes()

	// Expected routes
	expectedGetRoutes := []string{"/health", "/swagger/*any", "/api/v1/users", "/api/v1/users/:id"}
	expectedPostRoutes := []string{"/api/v1/users"}
	expectedPutRoutes := []string{"/api/v1/users/:id"}
	expectedDeleteRoutes := []string{"/api/v1/users/:id"}

	// Verify routes are registered
	routeMap := make(map[string][]string)
	for _, route := range routesList {
		routeMap[route.Method] = append(routeMap[route.Method], route.Path)
	}

	// Check health endpoint
	assert.Contains(t, routeMap["GET"], "/health")

	// Check swagger endpoint
	assert.Contains(t, routeMap["GET"], "/swagger/*any")

	// Check user endpoints
	assert.Contains(t, routeMap["POST"], "/api/v1/users")
	assert.Contains(t, routeMap["GET"], "/api/v1/users")
	assert.Contains(t, routeMap["GET"], "/api/v1/users/:id")
	assert.Contains(t, routeMap["PUT"], "/api/v1/users/:id")
	assert.Contains(t, routeMap["DELETE"], "/api/v1/users/:id")

	// Verify expected routes exist
	for _, route := range expectedGetRoutes {
		assert.Contains(t, routeMap["GET"], route)
	}
	for _, route := range expectedPostRoutes {
		assert.Contains(t, routeMap["POST"], route)
	}
	for _, route := range expectedPutRoutes {
		assert.Contains(t, routeMap["PUT"], route)
	}
	for _, route := range expectedDeleteRoutes {
		assert.Contains(t, routeMap["DELETE"], route)
	}

	// Check that routes are properly grouped
	assert.True(t, len(routesList) >= 7) // At least 7 routes should be registered
}

func TestSetupRoutes_Integration(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(MockUserService)
	userController := controllers.NewUserController(mockService)

	// Execute
	routes.SetupRoutes(router, userController)

	// Verify router is not nil and has routes
	assert.NotNil(t, router)

	routesList := router.Routes()
	assert.True(t, len(routesList) > 0)

	// Verify that all expected handler types are registered
	methodCounts := make(map[string]int)
	for _, route := range routesList {
		methodCounts[route.Method]++
	}

	// Should have at least one GET, POST, PUT, DELETE
	assert.True(t, methodCounts["GET"] >= 3)    // health, swagger, users, users/:id
	assert.True(t, methodCounts["POST"] >= 1)   // create user
	assert.True(t, methodCounts["PUT"] >= 1)    // update user
	assert.True(t, methodCounts["DELETE"] >= 1) // delete user
}

func TestSetupRoutes_APIVersioning(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(MockUserService)
	userController := controllers.NewUserController(mockService)

	// Execute
	routes.SetupRoutes(router, userController)

	// Get routes and verify API versioning
	routesList := router.Routes()

	// Check that all user routes are under /api/v1
	userRoutes := []string{}
	for _, route := range routesList {
		if route.Path != "/health" && route.Path != "/swagger/*any" {
			userRoutes = append(userRoutes, route.Path)
		}
	}

	// All user routes should be under /api/v1
	for _, route := range userRoutes {
		assert.Contains(t, route, "/api/v1/")
	}
}

func TestSetupRoutes_ControllerBinding(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(MockUserService)
	userController := controllers.NewUserController(mockService)

	// Execute - this should not panic
	assert.NotPanics(t, func() {
		routes.SetupRoutes(router, userController)
	})

	// Verify routes are bound to handlers
	routesList := router.Routes()

	// Each route should have a handler
	for _, route := range routesList {
		assert.NotNil(t, route.HandlerFunc)
		assert.NotEmpty(t, route.Method)
		assert.NotEmpty(t, route.Path)
	}
}

func TestSetupRoutes_GroupStructure(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	router := gin.New()
	mockService := new(MockUserService)
	userController := controllers.NewUserController(mockService)

	// Execute
	routes.SetupRoutes(router, userController)

	// Get routes and analyze structure
	routesList := router.Routes()

	// Group routes by their path structure
	topLevelRoutes := []string{}
	apiV1Routes := []string{}
	userRoutes := []string{}

	for _, route := range routesList {
		if route.Path == "/health" || route.Path == "/swagger/*any" {
			topLevelRoutes = append(topLevelRoutes, route.Path)
		} else if len(route.Path) > 7 && route.Path[:7] == "/api/v1" {
			apiV1Routes = append(apiV1Routes, route.Path)
			if len(route.Path) > 12 && route.Path[7:13] == "/users" {
				userRoutes = append(userRoutes, route.Path)
			}
		}
	}

	// Verify group structure
	assert.True(t, len(topLevelRoutes) >= 2) // health and swagger
	assert.True(t, len(apiV1Routes) >= 5)    // all user endpoints
	assert.True(t, len(userRoutes) >= 5)     // all user CRUD operations
}

// MockUserService for routes testing
type MockUserServiceForRoutes struct {
	mock.Mock
}

func (m *MockUserServiceForRoutes) CreateUser(req interface{}) (interface{}, error) {
	args := m.Called(req)
	return args.Get(0), args.Error(1)
}

func (m *MockUserServiceForRoutes) GetUserByID(id uint) (interface{}, error) {
	args := m.Called(id)
	return args.Get(0), args.Error(1)
}

func (m *MockUserServiceForRoutes) GetAllUsers(page, pageSize int) (interface{}, int64, error) {
	args := m.Called(page, pageSize)
	return args.Get(0), args.Get(1).(int64), args.Error(2)
}

func (m *MockUserServiceForRoutes) UpdateUser(id uint, req interface{}) (interface{}, error) {
	args := m.Called(id, req)
	return args.Get(0), args.Error(1)
}

func (m *MockUserServiceForRoutes) DeleteUser(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}
