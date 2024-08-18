package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	m "github.com/movax01h/kladovkin-telegram-bot/internal/models"
	"github.com/movax01h/kladovkin-telegram-bot/internal/repository"
	"time"
)

var _ repository.SubscriptionRepository = (*SQLiteSubscriptionRepository)(nil)

// SQLiteSubscriptionRepository implements the SubscriptionRepository interface using SQLite.
type SQLiteSubscriptionRepository struct {
	db *sql.DB
}

// NewSQLiteSubscriptionRepository creates a new instance of SQLiteSubscriptionRepository.
func NewSQLiteSubscriptionRepository(db *sql.DB) *SQLiteSubscriptionRepository {
	return &SQLiteSubscriptionRepository{db: db}
}

// CreateSubscription inserts or updates a subscription in the database.
func (r *SQLiteSubscriptionRepository) CreateSubscription(subscription *m.Subscription) error {
	query := `
        INSERT INTO subscriptions (user_id, city, storage, unit_size, status, created_at, updated_at)
        VALUES (?, ?, ?, ?, ?, ?, ?)
        ON CONFLICT(id) DO UPDATE SET
            user_id = excluded.user_id,
            city = excluded.city,
            storage = excluded.storage,
            unit_size = excluded.unit_size,
            status = excluded.status,
            updated_at = excluded.updated_at
    `
	_, err := r.db.Exec(
		query,
		subscription.UserID,
		subscription.City,
		subscription.Storage,
		subscription.UnitSize,
		subscription.Status,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("failed to save subscription: %w", err)
	}
	return nil
}

// GetSubscriptionByID retrieves a subscription by ID from the database.
func (r *SQLiteSubscriptionRepository) GetSubscriptionByID(id int64) (*m.Subscription, error) {
	query := `SELECT id, user_id, city, storage, unit_size, status, created_at, updated_at FROM subscriptions WHERE id = ?`
	row := r.db.QueryRow(query, id)

	var subscription m.Subscription
	err := row.Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.City,
		&subscription.Storage,
		&subscription.UnitSize,
		&subscription.Status,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get subscription by ID: %w", err)
	}
	return &subscription, nil
}

// GetSubscriptionsByUserID retrieves all subscriptions by user ID from the database.
func (r *SQLiteSubscriptionRepository) GetSubscriptionsByUserID(userID int64) ([]*m.Subscription, error) {
	query := `SELECT id, user_id, city, storage, unit_size, status, created_at, updated_at FROM subscriptions WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get subscriptions by user ID: %w", err)
	}
	defer rows.Close()

	var subscriptions []*m.Subscription
	for rows.Next() {
		var subscription m.Subscription
		if err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.City,
			&subscription.Storage,
			&subscription.UnitSize,
			&subscription.Status,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan subscription row: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during rows iteration: %w", err)
	}
	return subscriptions, nil
}

// GetAllSubscriptions retrieves all subscriptions from the database.
func (r *SQLiteSubscriptionRepository) GetAllSubscriptions() ([]*m.Subscription, error) {
	query := `SELECT id, user_id, city, storage, unit_size, status, created_at, updated_at FROM subscriptions`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*m.Subscription
	for rows.Next() {
		var subscription m.Subscription
		if err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.City,
			&subscription.Storage,
			&subscription.UnitSize,
			&subscription.Status,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan subscription row: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during rows iteration: %w", err)
	}
	return subscriptions, nil
}

// GetActiveSubscriptions retrieves all active subscriptions from the database.
func (r *SQLiteSubscriptionRepository) GetActiveSubscriptions() ([]*m.Subscription, error) {
	query := `SELECT id, user_id, city, storage, unit_size, status, created_at, updated_at FROM subscriptions WHERE status = 'active'`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get active subscriptions: %w", err)
	}
	defer rows.Close()

	var subscriptions []*m.Subscription
	for rows.Next() {
		var subscription m.Subscription
		if err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.City,
			&subscription.Storage,
			&subscription.UnitSize,
			&subscription.Status,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan active subscription row: %w", err)
		}
		subscriptions = append(subscriptions, &subscription)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during active subscriptions iteration: %w", err)
	}
	return subscriptions, nil
}

// UpdateSubscription updates a subscription in the database.
func (r *SQLiteSubscriptionRepository) UpdateSubscription(subscription *m.Subscription) error {
	query := `UPDATE subscriptions SET user_id = ?, city = ?, storage = ?, unit_size = ?, status = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(
		query,
		subscription.UserID,
		subscription.City,
		subscription.Storage,
		subscription.UnitSize,
		subscription.Status,
		time.Now(),
		subscription.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}
	return nil
}

// DeleteSubscription deletes a subscription from the database.
func (r *SQLiteSubscriptionRepository) DeleteSubscription(id int64) error {
	_, err := r.db.Exec("DELETE FROM subscriptions WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}
	return nil
}
