package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gomodule/redigo/redis"
	uuid "github.com/satori/go.uuid"
)

func Refresh(db *sql.DB, cache redis.Conn) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
		c, err := r.Cookie("session_token")
		if err != nil {
			if err == http.ErrNoCookie {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		sessionToken := c.Value

		response, err := cache.Do("GET", sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if response == nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// (END) The code uptil this point is the same as the first part of the `Welcome` route

		// Now, create a new session token for the current user
		newSessionToken := uuid.NewV4().String()
		_, err = cache.Do("SETEX", newSessionToken, "120", fmt.Sprintf("%s", response))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Delete the older session token
		_, err = cache.Do("DEL", sessionToken)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Set the new token as the users `session_token` cookie
		http.SetCookie(w, &http.Cookie{
			Name:    "session_token",
			Value:   newSessionToken,
			Expires: time.Now().Add(120 * time.Second),
		})
	}
}
