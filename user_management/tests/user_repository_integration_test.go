package tests

import (
	"testing"

	"github.com/IntouchOpec/user_management/repository"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

// Test repository constructor and interface implementation
func TestUserRepository_Constructor(t *testing.T) {
	db := &gorm.DB{}
	repo := repository.NewUserRepository(db)

	assert.NotNil(t, repo)

	// Verify it implements the interface
	var _ repository.UserRepository = repo
}

// Test repository type assertion and interface compliance
func TestUserRepository_InterfaceCompliance(t *testing.T) {
	db := &gorm.DB{}
	repo := repository.NewUserRepository(db)

	// Verify interface compliance by checking all methods exist
	// This will cause compilation error if any interface method is missing

	// Test that the repository implements UserRepository interface
	_, ok := repo.(repository.UserRepository)
	assert.True(t, ok, "Repository should implement UserRepository interface")
}

// Test that all interface methods are properly defined
func TestUserRepository_MethodsExist(t *testing.T) {
	// This test ensures that the concrete implementation has all required methods
	// by using interface assignment

	db := &gorm.DB{}
	var repo repository.UserRepository = repository.NewUserRepository(db)

	assert.NotNil(t, repo)

	// These tests verify method signatures exist without executing them
	// (since we don't have a real database connection)

	// The mere fact that this compiles means all interface methods are implemented
	// with correct signatures
}
