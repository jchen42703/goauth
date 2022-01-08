package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber"
	"github.com/gomodule/redigo/redis"
	_ "github.com/lib/pq" // <------------ here

	"github.com/jchen42703/goauth/controllers"
)

// The "db" package level variable will hold the reference to our database instance
var db *sql.DB
var cache redis.Conn

func initDB(dbInfo string) error {
	var err error
	// Connect to the postgres db
	//you might have to change the connection string to add your database credentials
	db, err = sql.Open("postgres", dbInfo)
	if err != nil {
		return err
	}

	err = db.Ping()
	return err
}

func initCache() error {
	// Initialize the redis connection to a redis instance running on your local machine
	conn, err := redis.DialURL("redis://localhost")
	if err != nil {
		return err
	}
	// Assign the connection to the package level `cache` variable
	cache = conn
	return nil
}

func main() {
	// initialize our database connection
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"))
	err := initDB(dbinfo)
	if err != nil {
		log.Fatal("failed to init db: ", err)
	}
	defer db.Close()

	err = initCache()
	if err != nil {
		log.Fatal("failed to init cache: ", err)
	}

	app := fiber.New()
	// GET /api/list
	app.Get("/api/list", func(c *fiber.Ctx) error {
		fmt.Println("🥉 Last handler")
		return c.SendString("Hello, World 👋!")
	})

	log.Fatal(app.Listen(":3000"))

	// "Signin" and "Signup" are handler that we will implement
	http.HandleFunc("/login", controllers.Login(db, cache))
	http.HandleFunc("/signup", controllers.Signup(db))
	http.HandleFunc("/welcome", controllers.Welcome(db, cache))
	http.HandleFunc("/refresh", controllers.Refresh(db, cache))
	// http.HandleFunc("/welcome", Welcome)

	// start the server on port 8000
	log.Println("Listening on port 8000...")
	log.Fatal(http.ListenAndServe(":8000", nil))
}
