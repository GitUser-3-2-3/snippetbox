package models

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
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
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}
	stmt := `INSERT INTO users (name, email, hashed_password, created)
               VALUES (?, ?, ?, UTC_TIMESTAMP())`

	_, err = mdl.DB.Exec(stmt, name, email, string(hashedPassword))
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && strings.Contains(mysqlErr.Message, "users_uc_email") {
			if mysqlErr.Number == 1062 {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	return nil
}

func (mdl *UserModel) Exists(_ int) (bool, error) {
	return false, nil
}

func (mdl *UserModel) Authenticate(email, password string) (int, error) {
	var hashedPassword []byte
	var id int
	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`

	err := mdl.DB.QueryRow(stmt, email).Scan(&id, &hashedPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}
	return id, nil
}
