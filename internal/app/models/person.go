package models

import (
	"strings"
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
	Phone   string
	Created time.Time
}

type NickName struct {
	ID       uuid.UUID
	NickName string
	Created  time.Time
}

type IPadress struct {
	ID     uuid.UUID
	IPtype string
	IP     string
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
	if p.ID == uuid.Nil {
		return errs.ErrBadData.Msg("invalid phone ID")
	}

	phone := strings.TrimSpace(p.Phone)
	if phone == "" {
		return errs.ErrBadData.Msg("phone number cannot be empty")
	}

	// Remove any non-digit characters for length validation
	digitsOnly := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}, phone)

	if len(digitsOnly) < 10 || len(digitsOnly) > 15 {
		return errs.ErrBadData.Msg("invalid phone number format: must be between 10 and 15 digits")
	}

	return nil
}

func NewPhone(id uuid.UUID, phone string) (Phone, error) {
	number := Phone{
		ID:      id,
		Phone:   phone,
		Created: time.Now(),
	}
	return number, number.Validate()
}
