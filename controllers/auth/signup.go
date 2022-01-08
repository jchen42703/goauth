package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password" db:"password"`
	Username string `json:"username" db:"username"`
}

func Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse and decode the request body into a new `Credentials` instance
		creds := &Credentials{}
		err := json.NewDecoder(r.Body).Decode(creds)
		if err != nil {
			// If there is something wrong with the request body, return a 400 status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Salt and hash the password using the bcrypt algorithm
		// The second argument is the cost of hashing, which we arbitrarily set as 8 (this value can be more or less, depending on the computing power you wish to utilize)
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), 8)
		if err != nil {
			return
		}

		// Next, insert the username, along with the hashed password into the database
		if _, err = db.Query("insert into users values ($1, $2)", creds.Username, string(hashedPassword)); err != nil {
			// If there is any issue with inserting into the database, return a 500 error
			w.WriteHeader(http.StatusInternalServerError)
			log.Println("failed to write to server: ", err)
			return
		}
		// We reach this point if the credentials we correctly stored in the database, and the default status of 200 is sent back
	}
}
