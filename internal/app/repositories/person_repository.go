package repositories

import (
	"context"
	"database/sql"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/models"
	"github.com/Lacky1234union/OsintTeleBot/internal/share/errs"
	"github.com/Lacky1234union/OsintTeleBot/internal/share/loger"
)

var dbLogger = loger.New("person_repository")

// DB is an interface that defines the database operations we need
type DB interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// PersonRepository defines the interface for person repository operations
type PersonRepository interface {
	Create(ctx context.Context, person models.Person) error
	FindByName(ctx context.Context, name string) (models.Person, error)
	FindByPhone(ctx context.Context, phone string) (models.Person, error)
	FindByEmail(ctx context.Context, email string) (models.Person, error)
}

// personRepository implements PersonRepository interface
type personRepository struct {
	db DB
}

// NewPersonRepository creates a new PersonRepository instance
func NewPersonRepository(db DB) PersonRepository {
	return &personRepository{
		db: db,
	}
}

func (r *personRepository) Create(ctx context.Context, person models.Person) error {
	if ctx == nil {
		dbLogger.WithContext(ctx).Warn("nil context in Create")
		return errs.ErrNilContext
	}
	_, err := r.db.ExecContext(ctx, "INSERT INTO persons (id, name, birthday, created, edited ) VALUES (?,?,?,?,?)",
		person.ID, person.Name, person.BirthDay, person.Created, person.Edited)
	if err != nil {
		dbLogger.WithContext(ctx).
			WithField("person_id", person.ID).
			WithError(err).
			Warn("error creating person")
		return errs.ErrPersonCreate.Err(err)
	}
	return nil
}

func (r *personRepository) FindByName(ctx context.Context, name string) (models.Person, error) {
	if ctx == nil {
		dbLogger.WithContext(ctx).Warn("nil context in FindByName")
		return models.Person{}, errs.ErrNilContext
	}
	var pers models.Person
	row := r.db.QueryRowContext(ctx, "SELECT * FROM persons WHERE name = ?", name)
	err := row.Scan(&pers.ID, &pers.Name, &pers.BirthDay, &pers.Created, &pers.Edited)
	if err != nil {
		logEntry := dbLogger.WithContext(ctx).WithField("name", name).WithError(err)
		if err == sql.ErrNoRows {
			logEntry.Warn("person not found by name")
			return models.Person{}, errs.ErrPersonNotFound
		}
		logEntry.Warn("error scanning person by name")
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}
	return pers, nil
}

func (r *personRepository) FindByPhone(ctx context.Context, phone string) (models.Person, error) {
	if ctx == nil {
		dbLogger.WithContext(ctx).Warn("nil context in FindByPhone")
		return models.Person{}, errs.ErrNilContext
	}
	var pers models.Person
	row := r.db.QueryRowContext(ctx, "SELECT id FROM phones WHERE phone = ?", phone)
	err := row.Scan(&pers.ID)
	if err != nil {
		logEntry := dbLogger.WithContext(ctx).WithField("phone", phone).WithError(err)
		if err == sql.ErrNoRows {
			logEntry.Warn("phone not found")
			return models.Person{}, errs.ErrPhoneNotFound
		}
		logEntry.Warn("error scanning phone")
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}

	row = r.db.QueryRowContext(ctx, "SELECT * FROM persons WHERE id = ?", pers.ID)
	err = row.Scan(&pers.ID, &pers.Name, &pers.BirthDay, &pers.Created, &pers.Edited)
	if err != nil {
		logEntry := dbLogger.WithContext(ctx).WithField("person_id", pers.ID).WithError(err)
		if err == sql.ErrNoRows {
			logEntry.Warn("person not found by phone")
			return models.Person{}, errs.ErrPersonNotFound
		}
		logEntry.Warn("error scanning person by phone")
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}
	return pers, nil
}

func (r *personRepository) FindByEmail(ctx context.Context, email string) (models.Person, error) {
	if ctx == nil {
		dbLogger.WithContext(ctx).Warn("nil context in FindByEmail")
		return models.Person{}, errs.ErrNilContext
	}
	var pers models.Person
	row := r.db.QueryRowContext(ctx, "SELECT id FROM emails WHERE email = ?", email)
	err := row.Scan(&pers.ID)
	if err != nil {
		logEntry := dbLogger.WithContext(ctx).WithField("email", email).WithError(err)
		if err == sql.ErrNoRows {
			logEntry.Warn("email not found")
			return models.Person{}, errs.ErrEmailNotFound
		}
		logEntry.Warn("error scanning email")
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}

	row = r.db.QueryRowContext(ctx, "SELECT * FROM persons WHERE id = ?", pers.ID)
	err = row.Scan(&pers.ID, &pers.Name, &pers.BirthDay, &pers.Created, &pers.Edited)
	if err != nil {
		logEntry := dbLogger.WithContext(ctx).WithField("person_id", pers.ID).WithError(err)
		if err == sql.ErrNoRows {
			logEntry.Warn("person not found by email")
			return models.Person{}, errs.ErrPersonNotFound
		}
		logEntry.Warn("error scanning person by email")
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}
	return pers, nil
}

func (r *personRepository) FindByNick(ctx context.Context, nick string) (models.Person, error) {
	if ctx == nil {
		dbLogger.WithContext(ctx).Warn("nil context in FindByNick")
		return models.Person{}, errs.ErrNilContext
	}
	var pers models.Person
	row := r.db.QueryRowContext(ctx, "SELECT id FROM nicks WHERE nick = ?", nick)
	err := row.Scan(&pers.ID)
	if err != nil {
		logEntry := dbLogger.WithContext(ctx).WithField("nick", nick).WithError(err)
		if err == sql.ErrNoRows {
			logEntry.Warn("nick not found")
			return models.Person{}, errs.ErrNickNotFound
		}
		logEntry.Warn("error scanning nick")
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}

	row = r.db.QueryRowContext(ctx, "SELECT * FROM persons WHERE id = ?", pers.ID)
	err = row.Scan(&pers.ID, &pers.Name, &pers.BirthDay, &pers.Created, &pers.Edited)
	if err != nil {
		logEntry := dbLogger.WithContext(ctx).WithField("person_id", pers.ID).WithError(err)
		if err == sql.ErrNoRows {
			logEntry.Warn("person not found by nick")
			return models.Person{}, errs.ErrPersonNotFound
		}
		logEntry.Warn("error scanning person by nick")
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}
	return pers, nil
}
