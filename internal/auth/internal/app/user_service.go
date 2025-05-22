package app

import (
	"errors"
	"time"

	"github.com/russunion/OsintTeleBot/internal/auth/internal/domain"
)

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExists         = errors.New("user already exists")
)

type userService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) domain.UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(username, email, password string) (*domain.User, error) {
	// Check if user exists
	if _, err := s.repo.GetByUsername(username); err == nil {
		return nil, ErrUserExists
	}
	if _, err := s.repo.GetByEmail(email); err == nil {
		return nil, ErrUserExists
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
		return nil, err
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) Login(username, password string) (*domain.User, error) {
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if !user.CheckPassword(password) {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *userService) GetByID(id int64) (*domain.User, error) {
	// This would typically be implemented in the repository
	// For now, we'll return an error
	return nil, ErrUserNotFound
}

func (s *userService) Update(user *domain.User) error {
	if user.Password != "" {
		if err := user.HashPassword(); err != nil {
			return err
		}
	}
	user.UpdatedAt = time.Now()
	return s.repo.Update(user)
}

func (s *userService) Delete(id int64) error {
	return s.repo.Delete(id)
}
