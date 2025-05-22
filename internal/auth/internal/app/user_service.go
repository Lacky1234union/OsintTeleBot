package app

import (
	"context"
	"errors"
	"time"

	"github.com/russunion/OsintTeleBot/internal/auth/internal/domain"
	"github.com/russunion/OsintTeleBot/internal/auth/internal/share/errs"
	"github.com/russunion/OsintTeleBot/internal/auth/internal/share/logger"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type userService struct {
	repo   domain.UserRepository
	logger *logger.Logger
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{
		repo:   repo,
		logger: logger.New(logger.LevelInfo),
	}
}

func (s *userService) Register(username, email, password string) (*domain.User, error) {
	ctx := context.Background()
	s.logger.Info(ctx, "Registering new user: %s", username)

	// Check if user exists
	if _, err := s.repo.GetByUsername(username); err == nil {
		s.logger.Warn(ctx, "User already exists: %s", username)
		return nil, errs.New("userService.Register", errs.ErrAlreadyExists, nil, "user already exists")
	}
	if _, err := s.repo.GetByEmail(email); err == nil {
		s.logger.Warn(ctx, "Email already exists: %s", email)
		return nil, errs.New("userService.Register", errs.ErrAlreadyExists, nil, "email already exists")
	}

	user := &domain.User{
		Username:  username,
		Email:     email,
		Password:  password,
		Role:      "user", // Default role
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.HashPassword(); err != nil {
		s.logger.Error(ctx, "Failed to hash password: %v", err)
		return nil, errs.New("userService.Register", errs.ErrInternal, err, "failed to hash password")
	}

	if err := s.repo.Create(user); err != nil {
		s.logger.Error(ctx, "Failed to create user: %v", err)
		return nil, errs.New("userService.Register", errs.ErrInternal, err, "failed to create user")
	}

	s.logger.Info(ctx, "Successfully registered user: %s", username)
	return user, nil
}

func (s *userService) Login(username, password string) (*domain.User, error) {
	ctx := context.Background()
	s.logger.Info(ctx, "Attempting login for user: %s", username)

	user, err := s.repo.GetByUsername(username)
	if err != nil {
		s.logger.Warn(ctx, "User not found: %s", username)
		return nil, errs.New("userService.Login", errs.ErrUnauthorized, err, "invalid credentials")
	}

	if !user.CheckPassword(password) {
		s.logger.Warn(ctx, "Invalid password for user: %s", username)
		return nil, errs.New("userService.Login", errs.ErrUnauthorized, nil, "invalid credentials")
	}

	s.logger.Info(ctx, "Successfully logged in user: %s", username)
	return user, nil
}

func (s *userService) GetByID(id int64) (*domain.User, error) {
	ctx := context.Background()
	s.logger.Info(ctx, "Getting user by ID: %d", id)

	// This would typically be implemented in the repository
	// For now, we'll return an error
	s.logger.Warn(ctx, "User not found: %d", id)
	return nil, errs.New("userService.GetByID", errs.ErrNotFound, nil, "user not found")
}

func (s *userService) Update(user *domain.User) error {
	ctx := context.Background()
	s.logger.Info(ctx, "Updating user: %s", user.Username)

	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			s.logger.Error(ctx, "Failed to hash password: %v", err)
			return errs.New("userService.Update", errs.ErrInternal, err, "failed to hash password")
		}
	}
	user.UpdatedAt = time.Now()

	if err := s.repo.Update(user); err != nil {
		s.logger.Error(ctx, "Failed to update user: %v", err)
		return errs.New("userService.Update", errs.ErrInternal, err, "failed to update user")
	}

	s.logger.Info(ctx, "Successfully updated user: %s", user.Username)
	return nil
}

func (s *userService) Delete(id int64) error {
	ctx := context.Background()
	s.logger.Info(ctx, "Deleting user with ID: %d", id)

	if err := s.repo.Delete(id); err != nil {
		s.logger.Error(ctx, "Failed to delete user: %v", err)
		return errs.New("userService.Delete", errs.ErrInternal, err, "failed to delete user")
	}

	s.logger.Info(ctx, "Successfully deleted user: %d", id)
	return nil
}
