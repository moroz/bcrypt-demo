package models

import (
	"errors"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/jmoiron/sqlx"
	"github.com/moroz/uuidv7-go"
)

type User struct {
	ID           uuidv7.UUID `db:"id"`
	Email        string      `db:"email"`
	PasswordHash string      `db:"password_hash"`
	InsertedAt   time.Time   `db:"inserted_at"`
	UpdatedAt    time.Time   `db:"updated_at"`
}

const USER_COLUMNS = "id, email, password_hash, inserted_at, updated_at"

var ARGON2_PARAMS = argon2id.Params{
	Memory:      46 * 1024, // 46 MiB
	Iterations:  1,
	Parallelism: 1,
	SaltLength:  16,
	KeyLength:   16,
}

func CreateUser(db *sqlx.DB, email, password, passwordConfirmation string) (*User, error) {
	// 信箱不能為空
	if email == "" {
		return nil, errors.New("Email cannot be blank!")
	}

	// 密碼不能為空
	if password == "" {
		return nil, errors.New("Password cannot be blank!")
	}

	// 密碼確認必須與密碼相符
	if password != passwordConfirmation {
		return nil, errors.New("Passwords do not match!")
	}

	// 套用 argon2id 密碼雜湊函數
	digest, err := argon2id.CreateHash(password, &ARGON2_PARAMS)
	if err != nil {
		return nil, err
	}

	result := User{}

	// 產生 UUIDv7
	id := uuidv7.Generate().String()

	// SQL INSERT 插入資料，用 RETURNING 即可與插入同時取得新的一筆資料
	// 如果插入成功，result 變數裡就會是剛新增的使用者資料
	err = db.Get(
		&result,
		`insert into users (id, email, password_hash)
        values ($1, $2, $3) returning `+USER_COLUMNS,
		// SQL 語法中三個佔位符 $1, $2, $3 需提供三個參數
		id, email, digest,
	)
	return &result, err
}

func AuthenticateUserByEmailPassword(db *sqlx.DB, email, password string) (*User, error) {
	result := User{}
	err := db.Get(&result, "select "+USER_COLUMNS+" from users where password_hash is not null and email = $1", email)
	if err != nil {
		return nil, err
	}

	match, err := argon2id.ComparePasswordAndHash(password, result.PasswordHash)
	if err != nil {
		return nil, err
	}
	if !match {
		return nil, errors.New("Invalid password")
	}

	return &result, nil
}
