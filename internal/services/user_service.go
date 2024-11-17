package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/yourusername/go-production-level/config"
	"github.com/yourusername/go-production-level/internal/models"
	"github.com/yourusername/go-production-level/internal/repository"
	"github.com/yourusername/go-production-level/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrEmailExists        = errors.New("email already exists")
)

type UserService interface {
	Create(ctx context.Context, user *models.User) error
	GetByID(ctx context.Context, id uint) (*models.UserResponse, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
	Delete(ctx context.Context, id uint) error
	List(ctx context.Context, offset, limit int) ([]models.UserResponse, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type UserServiceImpl struct {
	repo   repository.UserRepository
	redis  *redis.Client
	config *config.Config
}

func NewUserService(repo repository.UserRepository, redis *redis.Client, config *config.Config) UserService {
	return &UserServiceImpl{
		repo:   repo,
		redis:  redis,
		config: config,
	}
}

func (s *UserServiceImpl) Create(ctx context.Context, user *models.User) error {
	// Check if email already exists
	existingUser, err := s.repo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return ErrEmailExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return s.repo.Create(ctx, user)
}

func (s *UserServiceImpl) GetByID(ctx context.Context, id uint) (*models.UserResponse, error) {
	// Try to get from cache first
	cacheKey := fmt.Sprintf("user:%d", id)
	var userResp models.UserResponse

	// Check cache
	if err := s.redis.Get(ctx, cacheKey).Scan(&userResp); err == nil {
		return &userResp, nil
	}

	// If not in cache, get from database
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}

	userResp = models.UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
	}

	// Cache the result
	s.redis.Set(ctx, cacheKey, userResp, time.Hour)

	return &userResp, nil
}

func (s *UserServiceImpl) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *UserServiceImpl) Update(ctx context.Context, user *models.User) error {
	if user.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPassword)
	}

	err := s.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	// Invalidate cache
	s.redis.Del(ctx, "user:"+strconv.FormatUint(uint64(user.ID), 10))
	return nil
}

func (s *UserServiceImpl) Delete(ctx context.Context, id uint) error {
	err := s.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache
	s.redis.Del(ctx, "user:"+string(id))
	return nil
}

func (s *UserServiceImpl) List(ctx context.Context, offset, limit int) ([]models.UserResponse, error) {
	users, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	userResponses := make([]models.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = models.UserResponse{
			ID:        user.ID,
			CreatedAt: user.CreatedAt,
			Email:     user.Email,
			Name:      user.Name,
			Role:      user.Role,
		}
	}

	return userResponses, nil
}

func (s *UserServiceImpl) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user, s.config)
	if err != nil {
		return "", err
	}

	return token, nil
}
