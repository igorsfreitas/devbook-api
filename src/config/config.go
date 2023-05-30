package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	// ConnectionString is the connection string to the database
	ConnectionString = ""

	// Port is the port to listen on
	Port = 0
)

// Load config from file
func Load() {
	var err error

	if err = godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	Port, err = strconv.Atoi(os.Getenv("API_PORT"))
	if err != nil {
		Port = 9000
	}

	ConnectionString = fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)
}
