package services

import (
	"context"
	"testing"
	"time"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/models"
	"github.com/Lacky1234union/OsintTeleBot/internal/share/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockPersonRepository is a mock implementation of PersonRepository
type MockPersonRepository struct {
	mock.Mock
}

func (m *MockPersonRepository) Create(ctx context.Context, person models.Person) error {
	args := m.Called(ctx, person)
	return args.Error(0)
}

func (m *MockPersonRepository) FindByName(ctx context.Context, name string) (models.Person, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(models.Person), args.Error(1)
}

func (m *MockPersonRepository) FindByPhone(ctx context.Context, phone string) (models.Person, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(models.Person), args.Error(1)
}

func (m *MockPersonRepository) FindByEmail(ctx context.Context, email string) (models.Person, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.Person), args.Error(1)
}

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name          string
		user          models.Person
		setupMock     func(*MockPersonRepository)
		expectedError error
	}{
		{
			name: "successful registration",
			user: models.Person{
				ID:       uuid.New(),
				Name:     "John Doe",
				BirthDay: time.Now(),
				Created:  time.Now(),
				Edited:   time.Now(),
			},
			setupMock: func(mockRepo *MockPersonRepository) {
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "invalid user data",
			user: models.Person{
				ID:       uuid.Nil,
				Name:     "",
				BirthDay: time.Now(),
				Created:  time.Now(),
				Edited:   time.Now(),
			},
			setupMock:     func(mockRepo *MockPersonRepository) {},
			expectedError: errs.ErrBadData,
		},
		{
			name: "repository error",
			user: models.Person{
				ID:       uuid.New(),
				Name:     "John Doe",
				BirthDay: time.Now(),
				Created:  time.Now(),
				Edited:   time.Now(),
			},
			setupMock: func(mockRepo *MockPersonRepository) {
				mockRepo.On("Create", mock.Anything, mock.Anything).Return(errs.ErrPersonCreate)
			},
			expectedError: errs.ErrPersonCreate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPersonRepository)
			tt.setupMock(mockRepo)
			service := NewPersonService(mockRepo)

			err := service.RegisterUser(context.Background(), tt.user)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFindUserByEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		setupMock     func(*MockPersonRepository)
		expectedError error
	}{
		{
			name:  "successful find",
			email: "john@example.com",
			setupMock: func(mockRepo *MockPersonRepository) {
				mockRepo.On("FindByEmail", mock.Anything, "john@example.com").Return(models.Person{
					ID:       uuid.New(),
					Name:     "John Doe",
					BirthDay: time.Now(),
					Created:  time.Now(),
					Edited:   time.Now(),
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:          "empty email",
			email:         "",
			setupMock:     func(mockRepo *MockPersonRepository) {},
			expectedError: errs.ErrBadData,
		},
		{
			name:          "whitespace email",
			email:         "   ",
			setupMock:     func(mockRepo *MockPersonRepository) {},
			expectedError: errs.ErrBadData,
		},
		{
			name:  "not found",
			email: "nonexistent@example.com",
			setupMock: func(mockRepo *MockPersonRepository) {
				mockRepo.On("FindByEmail", mock.Anything, "nonexistent@example.com").Return(models.Person{}, errs.ErrEmailNotFound)
			},
			expectedError: errs.ErrEmailNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPersonRepository)
			tt.setupMock(mockRepo)
			service := NewPersonService(mockRepo)

			_, err := service.FindUserByEmail(context.Background(), tt.email)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFindUserByName(t *testing.T) {
	tests := []struct {
		name          string
		searchName    string
		setupMock     func(*MockPersonRepository)
		expectedError error
	}{
		{
			name:       "successful find",
			searchName: "John Doe",
			setupMock: func(mockRepo *MockPersonRepository) {
				mockRepo.On("FindByName", mock.Anything, "John Doe").Return(models.Person{
					ID:       uuid.New(),
					Name:     "John Doe",
					BirthDay: time.Now(),
					Created:  time.Now(),
					Edited:   time.Now(),
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:          "empty name",
			searchName:    "",
			setupMock:     func(mockRepo *MockPersonRepository) {},
			expectedError: errs.ErrBadData,
		},
		{
			name:          "whitespace name",
			searchName:    "   ",
			setupMock:     func(mockRepo *MockPersonRepository) {},
			expectedError: errs.ErrBadData,
		},
		{
			name:       "not found",
			searchName: "Nonexistent",
			setupMock: func(mockRepo *MockPersonRepository) {
				mockRepo.On("FindByName", mock.Anything, "Nonexistent").Return(models.Person{}, errs.ErrPersonNotFound)
			},
			expectedError: errs.ErrPersonNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPersonRepository)
			tt.setupMock(mockRepo)
			service := NewPersonService(mockRepo)

			_, err := service.FindUserByName(context.Background(), tt.searchName)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}

func TestFindUserByPhone(t *testing.T) {
	tests := []struct {
		name          string
		phone         string
		setupMock     func(*MockPersonRepository)
		expectedError error
	}{
		{
			name:  "successful find",
			phone: "+1234567890",
			setupMock: func(mockRepo *MockPersonRepository) {
				mockRepo.On("FindByPhone", mock.Anything, "+1234567890").Return(models.Person{
					ID:       uuid.New(),
					Name:     "John Doe",
					BirthDay: time.Now(),
					Created:  time.Now(),
					Edited:   time.Now(),
				}, nil)
			},
			expectedError: nil,
		},
		{
			name:          "empty phone",
			phone:         "",
			setupMock:     func(mockRepo *MockPersonRepository) {},
			expectedError: errs.ErrBadData,
		},
		{
			name:          "whitespace phone",
			phone:         "   ",
			setupMock:     func(mockRepo *MockPersonRepository) {},
			expectedError: errs.ErrBadData,
		},
		{
			name:  "not found",
			phone: "+9999999999",
			setupMock: func(mockRepo *MockPersonRepository) {
				mockRepo.On("FindByPhone", mock.Anything, "+9999999999").Return(models.Person{}, errs.ErrPhoneNotFound)
			},
			expectedError: errs.ErrPhoneNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockPersonRepository)
			tt.setupMock(mockRepo)
			service := NewPersonService(mockRepo)

			_, err := service.FindUserByPhone(context.Background(), tt.phone)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mockRepo.AssertExpectations(t)
		})
	}
}
