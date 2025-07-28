package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/IntouchOpec/user_management/models"
	"github.com/IntouchOpec/user_management/repository"
	"github.com/go-redis/redis/v8"
)

// UserService interface defines user business logic methods
type UserService interface {
	CreateUser(req models.UserRequest) (*models.UserResponse, error)
	GetUserByID(id uint) (*models.UserResponse, error)
	GetAllUsers(page, pageSize int) ([]models.UserResponse, int64, error)
	UpdateUser(id uint, req models.UserRequest) (*models.UserResponse, error)
	DeleteUser(id uint) error
}

// userService implements UserService interface
type userService struct {
	userRepo    repository.UserRepository
	redisClient *redis.Client
	ctx         context.Context
}

// NewUserService creates a new user service instance
func NewUserService(userRepo repository.UserRepository, redisClient *redis.Client) UserService {
	return &userService{
		userRepo:    userRepo,
		redisClient: redisClient,
		ctx:         context.Background(),
	}
}

// CreateUser creates a new user
func (s *userService) CreateUser(req models.UserRequest) (*models.UserResponse, error) {
	// Check if user with email already exists
	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, fmt.Errorf("user with email %s already exists", req.Email)
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Age:      req.Age,
		Phone:    req.Phone,
		Address:  req.Address,
		IsActive: true,
	}

	if req.IsActive != nil {
		user.IsActive = *req.IsActive
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}

	// Cache the user
	s.cacheUser(user)

	response := user.ToResponse()
	return &response, nil
}

// GetUserByID retrieves a user by ID
func (s *userService) GetUserByID(id uint) (*models.UserResponse, error) {
	// Try to get from cache first
	if cachedUser := s.getCachedUser(id); cachedUser != nil {
		response := cachedUser.ToResponse()
		return &response, nil
	}

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Cache the user
	s.cacheUser(user)

	response := user.ToResponse()
	return &response, nil
}

// GetAllUsers retrieves all users with pagination
func (s *userService) GetAllUsers(page, pageSize int) ([]models.UserResponse, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize

	users, err := s.userRepo.GetAll(offset, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get users: %v", err)
	}

	total, err := s.userRepo.Count()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %v", err)
	}

	var responses []models.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
		// Cache each user
		s.cacheUser(&user)
	}

	return responses, total, nil
}

// UpdateUser updates an existing user
func (s *userService) UpdateUser(id uint, req models.UserRequest) (*models.UserResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Check if email is being changed and if it already exists
	if user.Email != req.Email {
		existingUser, _ := s.userRepo.GetByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, fmt.Errorf("user with email %s already exists", req.Email)
		}
	}

	user.UpdateFromRequest(req)

	if err := s.userRepo.Update(user); err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}

	// Update cache
	s.cacheUser(user)

	response := user.ToResponse()
	return &response, nil
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(id uint) error {
	if err := s.userRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	// Remove from cache
	s.removeCachedUser(id)

	return nil
}

// cacheUser caches a user in Redis
func (s *userService) cacheUser(user *models.User) {
	if s.redisClient == nil {
		return
	}

	userJSON, err := json.Marshal(user)
	if err != nil {
		return
	}

	key := fmt.Sprintf("user:%d", user.ID)
	s.redisClient.Set(s.ctx, key, userJSON, 15*time.Minute)
}

// getCachedUser retrieves a user from Redis cache
func (s *userService) getCachedUser(id uint) *models.User {
	if s.redisClient == nil {
		return nil
	}

	key := fmt.Sprintf("user:%d", id)
	userJSON, err := s.redisClient.Get(s.ctx, key).Result()
	if err != nil {
		return nil
	}

	var user models.User
	if err := json.Unmarshal([]byte(userJSON), &user); err != nil {
		return nil
	}

	return &user
}

// removeCachedUser removes a user from Redis cache
func (s *userService) removeCachedUser(id uint) {
	if s.redisClient == nil {
		return
	}

	key := fmt.Sprintf("user:%d", id)
	s.redisClient.Del(s.ctx, key)
}
