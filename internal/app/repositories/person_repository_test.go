package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/models"
	"github.com/Lacky1234union/OsintTeleBot/internal/share/errs"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockDB is a mock implementation of DB interface
type MockDB struct {
	mock.Mock
}

func (m *MockDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	arguments := m.Called(ctx, query, args)
	return arguments.Get(0).(sql.Result), arguments.Error(1)
}

func (m *MockDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	arguments := m.Called(ctx, query, args)
	return arguments.Get(0).(*sql.Row)
}

func (m *MockDB) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	arguments := m.Called(ctx, dest, query, args)
	return arguments.Error(0)
}

func (m *MockDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	arguments := m.Called(ctx, dest, query, args)
	return arguments.Error(0)
}

// MockResult is a mock implementation of sql.Result
type MockResult struct {
	mock.Mock
}

func (m *MockResult) LastInsertId() (int64, error) {
	arguments := m.Called()
	return arguments.Get(0).(int64), arguments.Error(1)
}

func (m *MockResult) RowsAffected() (int64, error) {
	arguments := m.Called()
	return arguments.Get(0).(int64), arguments.Error(1)
}

func TestCreate(t *testing.T) {
	tests := []struct {
		name          string
		person        models.Person
		setupMock     func(*MockDB)
		expectedError error
	}{
		{
			name: "successful creation",
			person: models.Person{
				ID:       uuid.New(),
				Name:     "John Doe",
				BirthDay: time.Now(),
				Created:  time.Now(),
				Edited:   time.Now(),
			},
			setupMock: func(mockDB *MockDB) {
				mockResult := new(MockResult)
				mockResult.On("LastInsertId").Return(int64(1), nil)
				mockResult.On("RowsAffected").Return(int64(1), nil)
				mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).Return(mockResult, nil)
			},
			expectedError: nil,
		},
		{
			name: "database error",
			person: models.Person{
				ID:       uuid.New(),
				Name:     "John Doe",
				BirthDay: time.Now(),
				Created:  time.Now(),
				Edited:   time.Now(),
			},
			setupMock: func(mockDB *MockDB) {
				mockResult := new(MockResult)
				mockResult.On("LastInsertId").Return(int64(0), nil)
				mockResult.On("RowsAffected").Return(int64(0), nil)
				mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).Return(mockResult, errors.New("db error"))
			},
			expectedError: errs.ErrPersonCreate,
		},
		{
			name: "duplicate entry",
			person: models.Person{
				ID:       uuid.New(),
				Name:     "John Doe",
				BirthDay: time.Now(),
				Created:  time.Now(),
				Edited:   time.Now(),
			},
			setupMock: func(mockDB *MockDB) {
				mockResult := new(MockResult)
				mockResult.On("LastInsertId").Return(int64(0), nil)
				mockResult.On("RowsAffected").Return(int64(0), nil)
				mockDB.On("ExecContext", mock.Anything, mock.Anything, mock.Anything).Return(mockResult, errors.New("duplicate entry"))
			},
			expectedError: errs.ErrAlreadyExists,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)
			repo := NewPersonRepository(mockDB)

			err := repo.Create(context.Background(), tt.person)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mockDB.AssertExpectations(t)
		})
	}
}

func TestFindByName(t *testing.T) {
	tests := []struct {
		name          string
		searchName    string
		setupMock     func(*MockDB)
		expectedError error
	}{
		{
			name:       "successful find",
			searchName: "John",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("SelectContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "not found",
			searchName: "Nonexistent",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("SelectContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
			},
			expectedError: errs.ErrPersonNotFound,
		},
		{
			name:       "database error",
			searchName: "John",
			setupMock: func(mockDB *MockDB) {
				mockDB.On("SelectContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedError: errs.ErrDatabaseQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)
			repo := NewPersonRepository(mockDB)

			_, err := repo.FindByName(context.Background(), tt.searchName)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mockDB.AssertExpectations(t)
		})
	}
}

func TestFindByPhone(t *testing.T) {
	tests := []struct {
		name          string
		phone         string
		setupMock     func(*MockDB)
		expectedError error
	}{
		{
			name:  "successful find",
			phone: "+1234567890",
			setupMock: func(mockDB *MockDB) {
				row := &sql.Row{}
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(row)
				mockDB.On("GetContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "phone not found",
			phone: "+1234567890",
			setupMock: func(mockDB *MockDB) {
				row := &sql.Row{}
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(row)
				mockDB.On("GetContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
			},
			expectedError: errs.ErrPhoneNotFound,
		},
		{
			name:  "person not found",
			phone: "+1234567890",
			setupMock: func(mockDB *MockDB) {
				row := &sql.Row{}
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(row)
				mockDB.On("GetContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
			},
			expectedError: errs.ErrPersonNotFound,
		},
		{
			name:  "database error",
			phone: "+1234567890",
			setupMock: func(mockDB *MockDB) {
				row := &sql.Row{}
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(row)
				mockDB.On("GetContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedError: errs.ErrDatabaseQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)
			repo := NewPersonRepository(mockDB)

			_, err := repo.FindByPhone(context.Background(), tt.phone)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mockDB.AssertExpectations(t)
		})
	}
}

func TestFindByEmail(t *testing.T) {
	tests := []struct {
		name          string
		email         string
		setupMock     func(*MockDB)
		expectedError error
	}{
		{
			name:  "successful find",
			email: "john@example.com",
			setupMock: func(mockDB *MockDB) {
				row := &sql.Row{}
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(row)
				mockDB.On("GetContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:  "email not found",
			email: "nonexistent@example.com",
			setupMock: func(mockDB *MockDB) {
				row := &sql.Row{}
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(row)
				mockDB.On("GetContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
			},
			expectedError: errs.ErrEmailNotFound,
		},
		{
			name:  "person not found",
			email: "john@example.com",
			setupMock: func(mockDB *MockDB) {
				row := &sql.Row{}
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(row)
				mockDB.On("GetContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(sql.ErrNoRows)
			},
			expectedError: errs.ErrPersonNotFound,
		},
		{
			name:  "database error",
			email: "john@example.com",
			setupMock: func(mockDB *MockDB) {
				row := &sql.Row{}
				mockDB.On("QueryRowContext", mock.Anything, mock.Anything, mock.Anything).Return(row)
				mockDB.On("GetContext", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(errors.New("db error"))
			},
			expectedError: errs.ErrDatabaseQuery,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB := new(MockDB)
			tt.setupMock(mockDB)
			repo := NewPersonRepository(mockDB)

			_, err := repo.FindByEmail(context.Background(), tt.email)
			if tt.expectedError != nil {
				assert.ErrorIs(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
			mockDB.AssertExpectations(t)
		})
	}
}
