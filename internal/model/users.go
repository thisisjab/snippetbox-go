package model

import (
	"database/sql"
	"errors"
	"github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

type UserModelInterface interface {
	Insert(name, email, password string) error
	Authenticate(email, password string) (int, error)
	Exists(id int) (bool, error)
}

type User struct {
	ID             int
	FullName       string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

type UserModel struct {
	DB *sql.DB
}

func (m *UserModel) Insert(fullName, email, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := `INSERT INTO users (full_name, email, hashed_password, created)
	VALUES (?, ?, ?, strftime('%Y-%m-%d %H:%M:%S', 'now'))`

	_, err = m.DB.Exec(stmt, fullName, email, string(hashedPassword))
	if err != nil {
		var sqlite3Error sqlite3.Error

		if errors.As(err, &sqlite3Error) {
			if sqlite3Error.Code == sqlite3.ErrConstraint && strings.Contains(sqlite3Error.Error(), "users.email") {
				return ErrDuplicateEmail
			}
		}

		return err
	}
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
	var userID int
	var hashedPassword string

	stmt := `SELECT id, hashed_password FROM users WHERE email = ?`

	err := m.DB.QueryRow(stmt, email).Scan(&userID, &hashedPassword)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrInvalidCredentials
		}

		return 0, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}

		return 0, err
	}

	return userID, nil
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool
	stmt := "SELECT EXISTS(SELECT true FROM users WHERE id = ?)"
	err := m.DB.QueryRow(stmt, id).Scan(&exists)

	return exists, err
}
