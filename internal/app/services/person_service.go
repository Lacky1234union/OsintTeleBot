package services

import (
	"context"
	"strings"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/models"
	"github.com/Lacky1234union/OsintTeleBot/internal/app/repositories"
	"github.com/Lacky1234union/OsintTeleBot/internal/share/errs"
)

type PersonService struct {
	repo repositories.PersonRepository
}

func NewPersonService(repo repositories.PersonRepository) *PersonService {
	return &PersonService{repo: repo}
}

func (s *PersonService) RegisterUser(ctx context.Context, user models.Person) error {
	if err := user.Validate(); err != nil {
		return err
	}
	return s.repo.Create(ctx, user)
}

func (s *PersonService) FindUserByEmail(ctx context.Context, email string) (models.Person, error) {
	if ctx == nil {
		return models.Person{}, errs.ErrNilContext
	}
	if email = strings.TrimSpace(email); email == "" {
		return models.Person{}, errs.ErrBadData
	}
	return s.repo.FindByEmail(ctx, email)
}

func (s *PersonService) FindUserByName(ctx context.Context, name string) (models.Person, error) {
	if ctx == nil {
		return models.Person{}, errs.ErrNilContext
	}
	if name = strings.TrimSpace(name); name == "" {
		return models.Person{}, errs.ErrBadData
	}
	return s.repo.FindByName(ctx, name)
}

func (s *PersonService) FindUserByPhone(ctx context.Context, phone string) (models.Person, error) {
	if ctx == nil {
		return models.Person{}, errs.ErrNilContext
	}
	if phone = strings.TrimSpace(phone); phone == "" {
		return models.Person{}, errs.ErrBadData
	}
	return s.repo.FindByPhone(ctx, phone)
}

// Другие сервисные методы...
