package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

func New(dbPath string) (*Storage, error) {
	const fn = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open the database: %w", err)
	}
	return &Storage{db: db}, nil
}
