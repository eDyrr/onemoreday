package database

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

func New(dbPath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", dbPath)
	fmt.Println("connection to db successful")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
