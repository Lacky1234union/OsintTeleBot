package models

import (
	"time"

	"github.com/gofrs/uuid"
)

type User struct {
	ID         uuid.UUID
	TelegramID int64
	UserName   string
	Created    time.Time
	Edited     time.Time
}

type IPaddress struct {
	ID      uuid.UUID
	IP      string
	Created time.Time
}
