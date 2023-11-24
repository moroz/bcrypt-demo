package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/moroz/bcrypt-demo/config"
)

func main() {
	db := sqlx.MustConnect("postgres", config.DATABASE_URL)
}
