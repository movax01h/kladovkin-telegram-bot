package sqlite

import (
	"context"
	"database/sql"
	"github.com/movax01h/kladovkin-telegram-bot/internal/repository"
)

type sqliteUserRepository struct {
	db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) repository.UserRepository {
	return &sqliteUserRepository{db: db}
}

func (r *sqliteUserRepository) CreateUser(ctx context.Context, user *repository.User) error {
	// SQLite specific SQL query to insert a user
	_, err := r.db.ExecContext(ctx, "INSERT INTO users (...) VALUES (...)", ...)
	return err
}

func (r *sqliteUserRepository) GetUserByID(ctx context.Context, id int64) (*repository.User, error) {
	// SQLite specific SQL query to get a user by ID
	row := r.db.QueryRowContext(ctx, "SELECT ... FROM users WHERE id = ?", id)
	var user repository.User
	err := row.Scan(&user.ID, &user.Name, ...)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Other methods...