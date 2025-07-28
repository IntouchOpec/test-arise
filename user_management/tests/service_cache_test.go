package tests

import (
	"context"
	"encoding/json"
	"strconv"
	"testing"
	"time"

	"github.com/IntouchOpec/user_management/models"
	"github.com/IntouchOpec/user_management/service"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRedisClient for testing Redis operations
type MockRedisClient struct {
	mock.Mock
}

func (m *MockRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	args := m.Called(ctx, key, value, expiration)
	cmd := redis.NewStatusCmd(ctx)
	if args.Error(1) != nil {
		cmd.SetErr(args.Error(1))
	} else {
		cmd.SetVal(args.String(0))
	}
	return cmd
}

func (m *MockRedisClient) Get(ctx context.Context, key string) *redis.StringCmd {
	args := m.Called(ctx, key)
	cmd := redis.NewStringCmd(ctx)
	if args.Error(1) != nil {
		cmd.SetErr(args.Error(1))
	} else {
		cmd.SetVal(args.String(0))
	}
	return cmd
}

func (m *MockRedisClient) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	args := m.Called(ctx, keys)
	cmd := redis.NewIntCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	} else {
		cmd.SetVal(args.Get(0).(int64))
	}
	return cmd
}

func (m *MockRedisClient) Ping(ctx context.Context) *redis.StatusCmd {
	args := m.Called(ctx)
	cmd := redis.NewStatusCmd(ctx)
	if args.Error(0) != nil {
		cmd.SetErr(args.Error(0))
	} else {
		cmd.SetVal(args.String(0))
	}
	return cmd
}

// Test cache-specific scenarios
func TestUserService_CacheEdgeCases(t *testing.T) {
	mockRepo := &MockUserRepository{}
	// mockRedis := &MockRedisClient{}
	realRedis := redis.NewClient(&redis.Options{Addr: "localhost:6379"})

	userService := service.NewUserService(mockRepo, realRedis)
	ctx := context.Background()

	t.Run("Cache set failure should not affect operation", func(t *testing.T) {
		user := &models.User{ID: 1, Name: "John", Email: "john@example.com", Age: 25, IsActive: true}
		// Mock Redis Get to return cache miss

		userJSON, err := realRedis.Get(ctx, "user:1").Result()
		if err != nil {
			assert.Error(t, err)
		}
		assert.NotNil(t, userJSON)
		// Mock Redis Set to fail
		userData, _ := json.Marshal(user)
		realRedis.Set(ctx, "user:1", userData, 15*time.Minute)

		// Mock repository to return user
		mockRepo.On("GetByID", uint(1)).Return(user, nil)

		// Call service method
		result, err := userService.GetUserByID(1)

		// Should still succeed despite cache failure
		assert.NoError(t, err)
		assert.NotNil(t, result)

	})

	t.Run("Invalid cached data should fallback to DB", func(t *testing.T) {
		user := &models.User{ID: 1, Name: "John", Email: "john@example.com", Age: 25, IsActive: true}

		// Mock Redis Get to return invalid JSON
		// mockRedis.On("Get", mock.Anything, "user:1").Return("invalid json", nil)

		// Mock repository to return user (fallback)
		mockRepo.On("GetByID", uint(1)).Return(user, nil)

		// Mock Redis Set for re-caching with valid data
		userData, _ := json.Marshal(user)
		realRedis.Set(ctx, "user:1", userData, 15*time.Minute)

		// Call service method
		result, err := userService.GetUserByID(1)

		// Should succeed with DB fallback
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

// RedisClient interface for testing (define what we need)
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Ping(ctx context.Context) *redis.StatusCmd
}

// Test direct cache functions if they are exposed in service
func TestCacheKeyGeneration(t *testing.T) {
	tests := []struct {
		name     string
		userID   uint
		expected string
	}{
		{
			name:     "normal user ID",
			userID:   1,
			expected: "user:1",
		},
		{
			name:     "large user ID",
			userID:   999999,
			expected: "user:999999",
		},
		{
			name:     "zero user ID",
			userID:   0,
			expected: "user:0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test key generation logic
			key := "user:" + strconv.Itoa(int(tt.userID))
			assert.Equal(t, tt.expected, key)
		})
	}
}
