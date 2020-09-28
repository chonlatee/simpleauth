package mariadb

import (
	"context"
	"database/sql"
	"log"
)

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(ctx context.Context, username, email, password string) (int64, error) {

	var id int64
	err := u.DB.QueryRowContext(ctx, "insert into users (username, email, password) values (?, ?, ?) returning id", username, email, password).Scan(&id)
	if err != nil {
		log.Println(err)
	}
	return id, err
}

func (u *UserModel) GetByEmail(ctx context.Context, email string) (string, error) {
	var password string
	err := u.DB.QueryRowContext(ctx, "select password from users where email = ?", email).Scan(&password)

	if err != nil {
		return "", err
	}

	// no rows
	if len(password) == 0 {
		return "", sql.ErrNoRows
	}

	return password, nil

}
