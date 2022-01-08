package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func Login(db *sql.DB, cache redis.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse and decode the request body into a new `Credentials` instance
		creds := &Credentials{}
		err := json.NewDecoder(r.Body).Decode(creds)
		if err != nil {
			// If there is something wrong with the request body, return a 400 status
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Get the existing entry present in the database for the given username
		result := db.QueryRow("select password from users where username=$1", creds.Username)
		if err != nil {
			// If there is an issue with the database, return a 500 error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// We create another instance of `Credentials` to store the credentials we get from the database
		storedCreds := &Credentials{}
		// Store the obtained password in `storedCreds`
		err = result.Scan(&storedCreds.Password)
		if err != nil {
			// If an entry with the username does not exist, send an "Unauthorized"(401) status
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			// If the error is of any other type, send a 500 status
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Compare the stored hashed password, with the hashed version of the password that was received
		if err = bcrypt.CompareHashAndPassword([]byte(storedCreds.Password), []byte(creds.Password)); err != nil {
			// If the two passwords don't match, return a 401 status
			w.WriteHeader(http.StatusUnauthorized)
		}

		// If we reach this point, that means the users password was correct, and that they are authorized
		// The default 200 status is sent
		// Create a new random session token
		sessionToken := uuid.NewV4().String()
		// Set the token in the cache, along with the user whom it represents
		// The token has an expiry time of 120 seconds
		_, err = cache.Do("SETEX", sessionToken, "120", creds.Username)
		if err != nil {
			// If there is an error in setting the cache, return an internal server error
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Finally, we set the client cookie for "session_token" as the session token we just generated
		// we also set an expiry time of 120 seconds, the same as the cache
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   sessionToken,
			Expires: time.Now().Add(120 * time.Second),
		})
	}
}
