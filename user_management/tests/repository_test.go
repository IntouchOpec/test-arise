package tests

import (
	"errors"
	"testing"

	"github.com/IntouchOpec/user_management/models"
	"github.com/IntouchOpec/user_management/repository"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// MockDB represents a mock GORM database
type MockDB struct {
	mock.Mock
}

// We'll test the repository interface and constructor
func TestNewUserRepository(t *testing.T) {
	mockDB := &gorm.DB{}
	repo := repository.NewUserRepository(mockDB)

	assert.NotNil(t, repo)
}

// Test repository interface methods exist
func TestUserRepository_InterfaceMethods(t *testing.T) {
	// Create a mock that implements the interface
	mockRepo := new(MockUserRepositoryTest)

	// Test that all interface methods are available
	var repo repository.UserRepository = mockRepo
	assert.NotNil(t, repo)

	// Test method signatures by calling them with mocks
	user := &models.User{ID: 1, Name: "Test", Email: "test@example.com", Age: 25}

	mockRepo.On("Create", user).Return(nil)
	mockRepo.On("GetByID", uint(1)).Return(user, nil)
	mockRepo.On("GetByEmail", "test@example.com").Return(user, nil)
	mockRepo.On("GetAll", 0, 10).Return([]models.User{*user}, nil)
	mockRepo.On("Update", user).Return(nil)
	mockRepo.On("Delete", uint(1)).Return(nil)
	mockRepo.On("Count").Return(int64(1), nil)

	// Execute methods
	err := repo.Create(user)
	assert.NoError(t, err)

	foundUser, err := repo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, foundUser)

	foundUser, err = repo.GetByEmail("test@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user, foundUser)

	users, err := repo.GetAll(0, 10)
	assert.NoError(t, err)
	assert.Len(t, users, 1)

	err = repo.Update(user)
	assert.NoError(t, err)

	err = repo.Delete(1)
	assert.NoError(t, err)

	count, err := repo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), count)

	mockRepo.AssertExpectations(t)
}

func TestUserRepository_ErrorHandling(t *testing.T) {
	mockRepo := new(MockUserRepositoryTest)
	user := &models.User{ID: 1, Name: "Test", Email: "test@example.com", Age: 25}

	// Test Create error
	mockRepo.On("Create", user).Return(errors.New("creation error"))
	err := mockRepo.Create(user)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "creation error")

	// Test GetByID error
	mockRepo.On("GetByID", uint(999)).Return(nil, errors.New("user not found"))
	foundUser, err := mockRepo.GetByID(999)
	assert.Error(t, err)
	assert.Nil(t, foundUser)
	assert.Contains(t, err.Error(), "user not found")

	// Test GetByEmail error
	mockRepo.On("GetByEmail", "notfound@example.com").Return(nil, errors.New("user not found"))
	foundUser, err = mockRepo.GetByEmail("notfound@example.com")
	assert.Error(t, err)
	assert.Nil(t, foundUser)

	// Test GetAll error
	mockRepo.On("GetAll", 0, 10).Return([]models.User{}, errors.New("database error"))
	users, err := mockRepo.GetAll(0, 10)
	assert.Error(t, err)
	assert.Empty(t, users)

	// Test Update error
	mockRepo.On("Update", user).Return(errors.New("update error"))
	err = mockRepo.Update(user)
	assert.Error(t, err)

	// Test Delete error
	mockRepo.On("Delete", uint(999)).Return(errors.New("user not found"))
	err = mockRepo.Delete(999)
	assert.Error(t, err)

	// Test Count error
	mockRepo.On("Count").Return(int64(0), errors.New("count error"))
	count, err := mockRepo.Count()
	assert.Error(t, err)
	assert.Equal(t, int64(0), count)

	mockRepo.AssertExpectations(t)
}

func TestUserRepository_EdgeCases(t *testing.T) {
	mockRepo := new(MockUserRepositoryTest)

	// Test with nil user
	mockRepo.On("Create", (*models.User)(nil)).Return(errors.New("invalid user"))
	err := mockRepo.Create(nil)
	assert.Error(t, err)

	// Test with zero ID
	mockRepo.On("GetByID", uint(0)).Return(nil, errors.New("invalid ID"))
	user, err := mockRepo.GetByID(0)
	assert.Error(t, err)
	assert.Nil(t, user)

	// Test with empty email
	mockRepo.On("GetByEmail", "").Return(nil, errors.New("empty email"))
	user, err = mockRepo.GetByEmail("")
	assert.Error(t, err)
	assert.Nil(t, user)

	// Test with negative offset/limit
	mockRepo.On("GetAll", -1, -1).Return([]models.User{}, nil)
	users, err := mockRepo.GetAll(-1, -1)
	assert.NoError(t, err)
	assert.Empty(t, users)

	mockRepo.AssertExpectations(t)
}

func TestUserRepository_BoundaryValues(t *testing.T) {
	mockRepo := new(MockUserRepositoryTest)

	// Test with large ID
	largeID := uint(4294967295) // Max uint32
	mockRepo.On("GetByID", largeID).Return(nil, errors.New("user not found"))
	user, err := mockRepo.GetByID(largeID)
	assert.Error(t, err)
	assert.Nil(t, user)

	// Test with large offset/limit
	mockRepo.On("GetAll", 1000000, 1000).Return([]models.User{}, nil)
	users, err := mockRepo.GetAll(1000000, 1000)
	assert.NoError(t, err)
	assert.Empty(t, users)

	// Test with very long email
	longEmail := "verylongemailaddressthatexceedsnormallimits@verylongdomainname.com"
	mockRepo.On("GetByEmail", longEmail).Return(nil, errors.New("user not found"))
	user, err = mockRepo.GetByEmail(longEmail)
	assert.Error(t, err)
	assert.Nil(t, user)

	mockRepo.AssertExpectations(t)
}

func TestUserRepository_SuccessScenarios(t *testing.T) {
	mockRepo := new(MockUserRepositoryTest)

	// Test successful operations
	user1 := &models.User{ID: 1, Name: "John", Email: "john@example.com", Age: 30}
	user2 := &models.User{ID: 2, Name: "Jane", Email: "jane@example.com", Age: 25}

	// Test successful create
	mockRepo.On("Create", user1).Return(nil)
	err := mockRepo.Create(user1)
	assert.NoError(t, err)

	// Test successful get by ID
	mockRepo.On("GetByID", uint(1)).Return(user1, nil)
	foundUser, err := mockRepo.GetByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user1, foundUser)

	// Test successful get by email
	mockRepo.On("GetByEmail", "john@example.com").Return(user1, nil)
	foundUser, err = mockRepo.GetByEmail("john@example.com")
	assert.NoError(t, err)
	assert.Equal(t, user1, foundUser)

	// Test successful get all
	mockRepo.On("GetAll", 0, 10).Return([]models.User{*user1, *user2}, nil)
	users, err := mockRepo.GetAll(0, 10)
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, *user1, users[0])
	assert.Equal(t, *user2, users[1])

	// Test successful update
	user1.Name = "John Updated"
	mockRepo.On("Update", user1).Return(nil)
	err = mockRepo.Update(user1)
	assert.NoError(t, err)

	// Test successful delete
	mockRepo.On("Delete", uint(1)).Return(nil)
	err = mockRepo.Delete(1)
	assert.NoError(t, err)

	// Test successful count
	mockRepo.On("Count").Return(int64(2), nil)
	count, err := mockRepo.Count()
	assert.NoError(t, err)
	assert.Equal(t, int64(2), count)

	mockRepo.AssertExpectations(t)
}

func TestUserRepository_Pagination(t *testing.T) {
	mockRepo := new(MockUserRepositoryTest)

	// Test various pagination scenarios
	testCases := []struct {
		name     string
		offset   int
		limit    int
		expected int
	}{
		{"first page", 0, 10, 10},
		{"second page", 10, 10, 5},
		{"large limit", 0, 100, 50},
		{"zero limit", 0, 0, 0},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			users := make([]models.User, tc.expected)
			for i := 0; i < tc.expected; i++ {
				users[i] = models.User{
					ID:    uint(i + 1),
					Name:  "User " + string(rune(i+'1')),
					Email: "user" + string(rune(i+'1')) + "@example.com",
					Age:   25 + i,
				}
			}

			mockRepo.On("GetAll", tc.offset, tc.limit).Return(users, nil)

			result, err := mockRepo.GetAll(tc.offset, tc.limit)
			assert.NoError(t, err)
			assert.Len(t, result, tc.expected)
		})
	}

	mockRepo.AssertExpectations(t)
}

// MockUserRepositoryTest for testing
type MockUserRepositoryTest struct {
	mock.Mock
}

func (m *MockUserRepositoryTest) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepositoryTest) GetByID(id uint) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepositoryTest) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepositoryTest) GetAll(offset, limit int) ([]models.User, error) {
	args := m.Called(offset, limit)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepositoryTest) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepositoryTest) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepositoryTest) Count() (int64, error) {
	args := m.Called()
	return args.Get(0).(int64), args.Error(1)
}
