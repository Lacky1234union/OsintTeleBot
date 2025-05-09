package repositories

import (
	"context"

	"github.com/Lacky1234union/OsintTeleBot/internal/app/models"
	"github.com/jmoiron/sqlx"
)

type PersonRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *PersonRepository {
	return &PersonRepository{db: db}
}

func (r *PersonRepository) Create(ctx context.Context, person *models.Person) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO persons (id, name, birthday, created, edited ) VALUES (	?,?,?,?,?)", person.ID, person.Name, person.BirthDay, person.Created, person.Edited)
	return err
}

func (r *PersonRepository) FindByName(ctx context.Context, name string) (models.Person, error) {
	var pers models.Person
	row := r.db.QueryRow("SELECT * FROM persons WHERE name ='?' ", name)
	error := row.Scan(&pers.ID, &pers.Name, &pers.BirthDay, &pers.Created, &pers.Edited)
	return pers, error
}

func (r *PersonRepository) FindByPhone(ctx context.Context, phone string) (models.Person, error) {
	var pers models.Person
	row := r.db.QueryRow("SELECT id FROM phones WHERE  phone='?' ", phone)
	error := row.Scan(&pers.ID)
	if error != nil {
		return pers, error
	}

	row = r.db.QueryRow("SELECT * FROM persons WHERE  id='?' ", pers.ID)
	error = row.Scan(&pers.ID, &pers.Name, &pers.BirthDay, &pers.Created, &pers.Edited)
	return pers, error
}
