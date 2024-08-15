package sqlite

import (
	"database/sql"
	"fmt"
	m "github.com/movax01h/kladovkin-telegram-bot/internal/models"
	"github.com/movax01h/kladovkin-telegram-bot/internal/repository"
)

var _ repository.UserRepository = (*SQLiteUserRepository)(nil)

// NewSQLiteUserRepository creates a new SQLiteUserRepository.
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// SQLiteUserRepository implements the UserRepository interface using SQLite.
type SQLiteUserRepository struct {
	db *sql.DB
}

// CreateUser inserts or updates a user in the database.
func (r *SQLiteUserRepository) CreateUser(user *m.User) error {
	query := `INSERT INTO users (id, name, email, telegram_id, created_at, updated_at)
              VALUES (?, ?, ?, ?, ?, ?, ?)
              ON CONFLICT(id) DO UPDATE SET
              name=excluded.name, email=excluded.email, telegram_id=excluded.telegram_id, updated_at=excluded.updated_at, last_notified=excluded.last_notified`

	_, err := r.db.Exec(query, user.ID, user.Name, user.Email, user.TelegramID, CurrentTimestamp(), CurrentTimestamp())
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID from the database.
func (r *SQLiteUserRepository) GetUserByID(id int64) (*m.User, error) {
	query := `SELECT id, name, email, telegram_id, created_at, updated_at FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var user m.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.TelegramID, &user.CreatedAt, &user.UpdatedAt, &user.LastNotified)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database.
func (r *SQLiteUserRepository) GetAllUsers() ([]*m.User, error) {
	query := `SELECT id, name, email, telegram_id, created_at, updated_at FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	var users []*m.User
	for rows.Next() {
		var user m.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.TelegramID, &user.CreatedAt, &user.UpdatedAt, &user.LastNotified); err != nil {
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}
		users = append(users, &user)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during rows iteration: %w", err)
	}
	return users, nil
}

// UpdateUser updates a user in the database.
func (r *SQLiteUserRepository) UpdateUser(user *m.User) error {
	query := `UPDATE users SET name = ?, email = ?, telegram_id = ?, updated_at = ?, last_notified = ? WHERE id = ?`
	_, err := r.db.Exec(query, user.Name, user.Email, user.TelegramID, CurrentTimestamp(), CurrentTimestamp(), user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user from the database.
func (r *SQLiteUserRepository) DeleteUser(id int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
