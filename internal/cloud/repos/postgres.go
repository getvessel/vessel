package repos

import (
	"database/sql"
	"log"
	// Postgres driver expected later
	// _ "github.com/lib/pq"
)

type CloudDB struct {
	db *sql.DB
}

func NewCloudDB(dsn string) (*CloudDB, error) {
	// Stub setup for PG connection
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		log.Println("Warning: Failed to ping Postgres (if DSN is empty, this is expected)")
	}

	return &CloudDB{db: db}, nil
}
