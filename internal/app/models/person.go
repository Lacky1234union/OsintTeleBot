package models

import (
	"time"

	"github.com/Lacky1234union/OsintTeleBot/internal/share/errs"
	"github.com/google/uuid"
)

type Person struct {
	ID       uuid.UUID
	Name     string
	BirthDay time.Time
	Created  time.Time
	Edited   time.Time
}

type Email struct {
	ID       uuid.UUID
	Email    string
	Password string
	Created  time.Time
}

type Phone struct {
	ID      uuid.UUID
	Phone   int
	Created time.Time
}

type NickName struct {
	ID       uuid.UUID
	NickName string
	Created  time.Time
}

func (p *Person) Validate() error {
	if p.ID == uuid.Nil || p.Name == "" || p.BirthDay.Unix() > time.Now().Unix() {
		return errs.ErrBadData
	}
	return nil
}

func NewPerson(name string, birthday time.Time) (Person, error) {
	pers := Person{
		uuid.New(),
		name,
		birthday,
		time.Now(),
		time.Now(),
	}
	return pers, pers.Validate()
}

func (e *Email) Validate() error {
	if e.ID == uuid.Nil || e.Email == "" || e.Password == "" {
		return errs.ErrBadData
	}
	return nil
}

func NewEmail(id uuid.UUID, email string, passwd string) (Email, error) {
	temp := Email{
		id,
		email,
		passwd,
		time.Now(),
	}
	if temp.Validate() == nil {
		return temp, nil
	}
	return temp, temp.Validate()
}

func (n *NickName) Validate() error {
	if n.ID == uuid.Nil || n.NickName == "" {
		return errs.ErrBadData
	}
	return nil
}

func NewNickName(id uuid.UUID, nick string) (NickName, error) {
	name := NickName{
		id,
		nick,
		time.Now(),
	}
	return name, name.Validate()
}

func (p *Phone) Validate() error {
	if p.ID == uuid.Nil || p.Phone <= 0 {
		return errs.ErrBadData
	}
	return nil
}

func NewPhone(id uuid.UUID, phone int) (Phone, error) {
	number := Phone{
		id,
		phone,
		time.Now(),
	}
	return number, number.Validate()
}
