package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/moroz/password-demo/config"
	"github.com/moroz/password-demo/models"
)

func main() {
	db := sqlx.MustConnect("postgres", config.DATABASE_URL)

	var email, password, passwordConfirmation string
	fmt.Print("Email: ")
	fmt.Scanln(&email)
	fmt.Print("Password: ")
	fmt.Scanln(&password)
	fmt.Print("Confirm password: ")
	fmt.Scanln(&passwordConfirmation)

	user, err := models.CreateUser(db, email, password, passwordConfirmation)
	if err != nil {
		log.Println(err)
	} else {
		log.Printf("Created user with ID: %s\n", user.ID.String())
	}
}
