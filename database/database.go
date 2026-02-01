package database

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func InitDB(connectionString string) (*sql.DB, error) {
	// Debug: log connection string length (not the actual string for security)
	log.Printf("DB connection string length: %d", len(connectionString))

	if connectionString == "" {
		log.Fatal("DB_CONN environment variable is empty!")
	}

	// Open database menggunakan pgx driver
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Printf("sql.Open error: %v", err)
		return nil, err
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test connection
	log.Println("Testing database connection...")
	err = db.Ping()
	if err != nil {
		log.Printf("db.Ping error: %v", err)
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}