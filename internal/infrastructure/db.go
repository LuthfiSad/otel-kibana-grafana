package infrastructure

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/lib/pq"
)

func NewDB() *sql.DB {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:***@host.docker.internal:5432/testing-otel?sslmode=disable"
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	return db
}
