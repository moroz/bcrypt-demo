package models

import (
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/moroz/uuidv7-go"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuidv7.UUID `db:"id"`
	Email        string      `db:"email"`
	PasswordHash string      `db:"password_hash"`
	InsertedAt   time.Time   `db:"inserted_at"`
	UpdatedAt    time.Time   `db:"updated_at"`
}

const USER_COLUMNS = "id, email, password_hash, inserted_at, updated_at"

func CreateUser(db *sqlx.DB, email, password, passwordConfirmation string) (*User, error) {
	if email == "" {
		return nil, errors.New("Email cannot be blank!")
	}
	if password == "" {
		return nil, errors.New("Password cannot be blank!")
	}
	if password != passwordConfirmation {
		return nil, errors.New("Passwords do not match!")
	}
	digest, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	passwordHash := string(digest)

	result := User{}
	id := uuidv7.Generate()
	err = db.Get(&result, "insert into users (id, email, password_hash) values ($1, $2, $3) returning "+USER_COLUMNS, id.String(), email, passwordHash)
	return &result, err
}

func AuthenticateUserByEmailPassword(db *sqlx.DB, email, password string) (*User, error) {
	result := User{}
	err := db.Get(&result, "select "+USER_COLUMNS+" from users where password_hash is not null and email = $1", email)
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.PasswordHash), []byte(password))
	if err != nil {
		return nil, err
	}

	return &result, nil
}
