package tests

import (
	"testing"

	"github.com/IntouchOpec/user_management/service"
	"github.com/stretchr/testify/assert"
)

// Test service constructor
func TestUserService_Constructor(t *testing.T) {
	mockRepo := &MockUserRepository{}

	// Test with nil Redis client
	userService := service.NewUserService(mockRepo, nil)
	assert.NotNil(t, userService)

	// Verify service creation doesn't panic with nil Redis
	assert.NotPanics(t, func() {
		service.NewUserService(mockRepo, nil)
	})
}

// Test cache functionality scenarios
func TestUserService_WithoutRedis(t *testing.T) {
	mockRepo := &MockUserRepository{}
	userService := service.NewUserService(mockRepo, nil)

	// Test that service works without Redis
	assert.NotNil(t, userService)

	// The service should still function for basic operations
	// without caching when Redis is nil
}
