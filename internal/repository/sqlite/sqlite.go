package sqlite

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

// InitializeDatabase sets up the SQLite database with the required tables.
func InitializeDatabase(db *sql.DB) error {
	createTablesSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		telegram_id INTEGER NOT NULL,
		username TEXT,
		first_name TEXT,
		last_name TEXT,
		last_notified DATETIME,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS units (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		city TEXT NOT NULL,
		storage TEXT NOT NULL,
		size TEXT NOT NULL,
		price REAL NOT NULL,
		available BOOLEAN NOT NULL,
		description TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS subscriptions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		city TEXT NOT NULL,
		storage TEXT NOT NULL,
		unit_size TEXT NOT NULL,
		status TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	`

	_, err := db.Exec(createTablesSQL)
	if err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	log.Println("Database tables created or already exist.")
	return nil
}

// NewSQLiteDB initializes a new SQLite database connection.
func NewSQLiteDB(dbFilePath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %v", err)
	}

	// Set reasonable connection limits
	db.SetMaxOpenConns(1)    // SQLite uses a single writer, so only one open connection is allowed.
	db.SetConnMaxLifetime(0) // Connections are not closed automatically.

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping SQLite database: %v", err)
	}

	log.Println("SQLite database connected successfully.")
	return db, nil
}

// CurrentTimestamp returns the current time formatted as a string for use in SQLite.
func CurrentTimestamp() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
