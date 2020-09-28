package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type User struct {
	ID          int
	Username    string
	Email       string
	Password    string
	CreatedDate time.Time
}
