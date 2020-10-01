package models

import (
	"context"
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: no matching record found")

type UserRepository interface {
	Insert(context.Context, string, string, string) (int64, error)
	GetByEmail(context.Context, string) (string, error)
}

type AccessTokenRepository interface {
	Insert(context.Context, string, string, time.Time, time.Time) (int64, error)
}

type User struct {
	ID          int
	Username    string
	Email       string
	Password    string
	CreatedDate time.Time
}

type AccessToken struct {
	ID                  int
	Email               string
	Token               string
	RefreshToken        string
	CreatedDate         time.Time
	TokenExpired        time.Time
	RefreshTokenExpired time.Time
}
