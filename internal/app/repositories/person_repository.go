package repositories

import (
	"context"
	"database/sql"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/models"
	"github.com/Lacky1234union/OsintTeleBot/internal/share/errs"
	"github.com/jmoiron/sqlx"
)

type PersonRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(ctx context.Context, person models.Person) error {
	if ctx == nil {
		return errs.ErrNilContext
	}
	_, err := r.db.ExecContext(ctx, "INSERT INTO persons (id, name, birthday, created, edited ) VALUES (?,?,?,?,?)",
		person.ID, person.Name, person.BirthDay, person.Created, person.Edited)
	if err != nil {
		return errs.ErrPersonCreate.Err(err)
	}
	return nil
}

func (r *PersonRepository) FindByName(ctx context.Context, name string) (models.Person, error) {
	if ctx == nil {
		return models.Person{}, errs.ErrNilContext
	}
	var pers models.Person
	row := r.db.QueryRowContext(ctx, "SELECT * FROM persons WHERE name = ?", name)
	err := row.Scan(&pers.ID, &pers.Name, &pers.BirthDay, &pers.Created, &pers.Edited)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Person{}, errs.ErrPersonNotFound
		}
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}
	return pers, nil
}

func (r *PersonRepository) FindByPhone(ctx context.Context, phone string) (models.Person, error) {
	if ctx == nil {
		return models.Person{}, errs.ErrNilContext
	}
	var pers models.Person
	row := r.db.QueryRowContext(ctx, "SELECT id FROM phones WHERE phone = ?", phone)
	err := row.Scan(&pers.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Person{}, errs.ErrPhoneNotFound
		}
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}

	row = r.db.QueryRowContext(ctx, "SELECT * FROM persons WHERE id = ?", pers.ID)
	err = row.Scan(&pers.ID, &pers.Name, &pers.BirthDay, &pers.Created, &pers.Edited)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Person{}, errs.ErrPersonNotFound
		}
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}
	return pers, nil
}

func (r *PersonRepository) FindByEmail(ctx context.Context, email string) (models.Person, error) {
	if ctx == nil {
		return models.Person{}, errs.ErrNilContext
	}
	var pers models.Person
	row := r.db.QueryRowContext(ctx, "SELECT id FROM emails WHERE email = ?", email)
	err := row.Scan(&pers.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Person{}, errs.ErrEmailNotFound
		}
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}

	row = r.db.QueryRowContext(ctx, "SELECT * FROM persons WHERE id = ?", pers.ID)
	err = row.Scan(&pers.ID, &pers.Name, &pers.BirthDay, &pers.Created, &pers.Edited)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Person{}, errs.ErrPersonNotFound
		}
		return models.Person{}, errs.ErrDatabaseScan.Err(err)
	}
	return pers, nil
}
