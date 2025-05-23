package repositories

import (
	"context"
	"database/sql"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/models"
	"github.com/Lacky1234union/OsintTeleBot/internal/share/errs"
	"github.com/Lacky1234union/OsintTeleBot/internal/share/loger"
	"github.com/jmoiron/sqlx"
)

var dbLogger = loger.New("person_repository")

type PersonRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(ctx context.Context, person models.Person) error {
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

func (r *PersonRepository) FindByName(ctx context.Context, name string) (models.Person, error) {
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

func (r *PersonRepository) FindByPhone(ctx context.Context, phone string) (models.Person, error) {
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

func (r *PersonRepository) FindByEmail(ctx context.Context, email string) (models.Person, error) {
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
