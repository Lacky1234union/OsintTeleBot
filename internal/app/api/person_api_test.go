package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/models"
	"github.com/Lacky1234union/OsintTeleBot/internal/share/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPersonService struct {
	mock.Mock
}

func (m *MockPersonService) RegisterUser(ctx context.Context, user models.Person) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockPersonService) FindUserByEmail(ctx context.Context, email string) (models.Person, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.Person), args.Error(1)
}
func (m *MockPersonService) FindUserByName(ctx context.Context, name string) (models.Person, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(models.Person), args.Error(1)
}
func (m *MockPersonService) FindUserByPhone(ctx context.Context, phone string) (models.Person, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(models.Person), args.Error(1)
}

func TestRegisterUserHandler(t *testing.T) {
	mockService := new(MockPersonService)
	api := NewPersonAPI(mockService)

	person := models.Person{
		ID:       uuid.New(),
		Name:     "John Doe",
		BirthDay: time.Now(),
		Created:  time.Now(),
		Edited:   time.Now(),
	}
	body, _ := json.Marshal(person)

	t.Run("success", func(t *testing.T) {
		mockService.On("RegisterUser", mock.Anything, person).Return(nil)
		req := httptest.NewRequest(http.MethodPost, "/api/person", bytes.NewReader(body))
		rw := httptest.NewRecorder()
		api.RegisterUserHandler(rw, req)
		assert.Equal(t, http.StatusCreated, rw.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/person", bytes.NewReader([]byte("not-json")))
		rw := httptest.NewRecorder()
		api.RegisterUserHandler(rw, req)
		assert.Equal(t, http.StatusBadRequest, rw.Code)
	})

	t.Run("service error", func(t *testing.T) {
		mockService.On("RegisterUser", mock.Anything, person).Return(errs.ErrAlreadyExists)
		req := httptest.NewRequest(http.MethodPost, "/api/person", bytes.NewReader(body))
		rw := httptest.NewRecorder()
		api.RegisterUserHandler(rw, req)
		assert.Equal(t, http.StatusConflict, rw.Code)
		mockService.AssertExpectations(t)
	})
}

func TestFindUserByEmailHandler(t *testing.T) {
	mockService := new(MockPersonService)
	api := NewPersonAPI(mockService)
	person := models.Person{ID: uuid.New(), Name: "John Doe"}

	t.Run("success", func(t *testing.T) {
		mockService.On("FindUserByEmail", mock.Anything, "john@example.com").Return(person, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/person/email?email=john@example.com", nil)
		rw := httptest.NewRecorder()
		api.FindUserByEmailHandler(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("FindUserByEmail", mock.Anything, "notfound@example.com").Return(models.Person{}, errs.ErrEmailNotFound)
		req := httptest.NewRequest(http.MethodGet, "/api/person/email?email=notfound@example.com", nil)
		rw := httptest.NewRecorder()
		api.FindUserByEmailHandler(rw, req)
		assert.Equal(t, http.StatusNotFound, rw.Code)
		mockService.AssertExpectations(t)
	})
}

func TestFindUserByNameHandler(t *testing.T) {
	mockService := new(MockPersonService)
	api := NewPersonAPI(mockService)
	person := models.Person{ID: uuid.New(), Name: "John Doe"}

	t.Run("success", func(t *testing.T) {
		mockService.On("FindUserByName", mock.Anything, "John Doe").Return(person, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/person/name?name=John+Doe", nil)
		rw := httptest.NewRecorder()
		api.FindUserByNameHandler(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("FindUserByName", mock.Anything, "NotExist").Return(models.Person{}, errs.ErrPersonNotFound)
		req := httptest.NewRequest(http.MethodGet, "/api/person/name?name=NotExist", nil)
		rw := httptest.NewRecorder()
		api.FindUserByNameHandler(rw, req)
		assert.Equal(t, http.StatusNotFound, rw.Code)
		mockService.AssertExpectations(t)
	})
}

func TestFindUserByPhoneHandler(t *testing.T) {
	mockService := new(MockPersonService)
	api := NewPersonAPI(mockService)
	person := models.Person{ID: uuid.New(), Name: "John Doe"}

	t.Run("success", func(t *testing.T) {
		mockService.On("FindUserByPhone", mock.Anything, "+1234567890").Return(person, nil)
		req := httptest.NewRequest(http.MethodGet, "/api/person/phone?phone=+1234567890", nil)
		rw := httptest.NewRecorder()
		api.FindUserByPhoneHandler(rw, req)
		assert.Equal(t, http.StatusOK, rw.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("FindUserByPhone", mock.Anything, "+0000000000").Return(models.Person{}, errs.ErrPhoneNotFound)
		req := httptest.NewRequest(http.MethodGet, "/api/person/phone?phone=+0000000000", nil)
		rw := httptest.NewRecorder()
		api.FindUserByPhoneHandler(rw, req)
		assert.Equal(t, http.StatusNotFound, rw.Code)
		mockService.AssertExpectations(t)
	})
}
