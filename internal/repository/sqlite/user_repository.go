package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// User represents a user in the system.
type User struct {
	ID        int64
	Username  string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SQLiteUserRepository implements the UserRepository interface using SQLite.
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new instance of SQLiteUserRepository.
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// Create inserts a new user into the database.
func (r *SQLiteUserRepository) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, email, created_at, updated_at)
		VALUES (?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, user.Username, user.Email, now, now)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	user.ID = id
	user.CreatedAt = now
	user.UpdatedAt = now
	return nil
}

// GetByID retrieves a user by its ID.
func (r *SQLiteUserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE id = ?
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to retrieve user by id: %w", err)
	}

	return user, nil
}

// GetByEmail retrieves a user by their email address.
func (r *SQLiteUserRepository) GetByEmail(ctx context.Context, email string) (*User, error) {
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
		WHERE email = ?
	`

	user := &User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, fmt.Errorf("failed to retrieve user by email: %w", err)
	}

	return user, nil
}

// Update updates an existing user in the database.
func (r *SQLiteUserRepository) Update(ctx context.Context, user *User) error {
	query := `
		UPDATE users
		SET username = ?, email = ?, updated_at = ?
		WHERE id = ?
	`
	user.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query, user.Username, user.Email, user.UpdatedAt, user.ID)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

// Delete deletes a user from the database.
func (r *SQLiteUserRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM users
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

// GetAll retrieves all users from the database.
func (r *SQLiteUserRepository) GetAll(ctx context.Context) ([]*User, error) {
	query := `
		SELECT id, username, email, created_at, updated_at
		FROM users
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all users: %w", err)
	}
	defer rows.Close()

	var users []*User
	for rows.Next() {
		user := &User{}
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	return users, nil
}
