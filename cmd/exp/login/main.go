package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/moroz/bcrypt-demo/config"
	"github.com/moroz/bcrypt-demo/models"
)

func main() {
	db := sqlx.MustConnect("postgres", config.DATABASE_URL)

	var email, password string
	fmt.Print("Email: ")
	fmt.Scanln(&email)
	fmt.Print("Password: ")
	fmt.Scanln(&password)

	user, err := models.AuthenticateUserByEmailPassword(db, email, password)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Signed in user with ID: %s\n", user.ID.String())
	}
}
