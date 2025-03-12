package models

import (
	"database/sql"
	"time"
)

type Users struct {
	ID             string
	Name           string
	Email          string
	HashedPassword []byte
	CreatedAt      time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (mdl *UserModel) Insert(name, email, password string) error {
	return nil
}

func (mdl *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

func (mdl *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}
