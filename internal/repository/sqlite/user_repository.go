package sqlite

import (
	"database/sql"
	"errors"
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
	query := `
		INSERT INTO users (id, telegram_id, username, first_name, last_name, last_notified, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			telegram_id = excluded.telegram_id,
			username = excluded.username,
			first_name = excluded.first_name,
			last_name = excluded.last_name,
			last_notified = excluded.last_notified,
			updated_at = excluded.updated_at
	`
	_, err := r.db.Exec(
		query,
		user.ID,
		user.TelegramID,
		user.UserName,
		user.FirstName,
		user.LastName,
		user.LastNotified,
		user.CreatedAt,
		user.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save user: %w", err)
	}
	return nil
}

// GetUserByID retrieves a user by ID from the database.
func (r *SQLiteUserRepository) GetUserByID(id int64) (*m.User, error) {
	query := `SELECT id, telegram_id, username, first_name, last_name, last_notified, created_at, updated_at FROM users WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var user m.User
	err := row.Scan(&user.ID, &user.TelegramID, &user.UserName, &user.FirstName, &user.LastName, &user.LastNotified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}
	return &user, nil
}

// GetByTelegramID retrieves a user by Telegram ID from the database.
func (r *SQLiteUserRepository) GetByTelegramID(id int64) (*m.User, error) {
	query := `SELECT id, telegram_id, username, first_name, last_name, last_notified, created_at, updated_at FROM users WHERE telegram_id = ?`
	row := r.db.QueryRow(query, id)

	var user m.User
	err := row.Scan(&user.ID, &user.TelegramID, &user.UserName, &user.FirstName, &user.LastName, &user.LastNotified, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by Telegram ID: %w", err)
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database.
func (r *SQLiteUserRepository) GetAllUsers() ([]*m.User, error) {
	query := `SELECT id, telegram_id, username, first_name, last_name, last_notified, created_at, updated_at FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer rows.Close()

	users := make([]*m.User, 0)
	for rows.Next() {
		var user m.User
		err := rows.Scan(&user.ID, &user.TelegramID, &user.UserName, &user.FirstName, &user.LastName, &user.LastNotified, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, &user)
	}
	return users, nil
}

// UpdateUser updates a user in the database.
func (r *SQLiteUserRepository) UpdateUser(user *m.User) error {
	query := `
		UPDATE users
		SET telegram_id = ?, username = ?, first_name = ?, last_name = ?, last_notified = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(
		query,
		user.TelegramID,
		user.UserName,
		user.FirstName,
		user.LastName,
		user.LastNotified,
		user.UpdatedAt,
		user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// DeleteUser deletes a user from the database.
func (r *SQLiteUserRepository) DeleteUser(user *m.User) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.Exec(query, user.ID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}
