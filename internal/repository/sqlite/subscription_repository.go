package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Subscription represents a user subscription to a unit.
type Subscription struct {
	ID          int64
	UserID      int64
	UnitID      int64
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SQLiteSubscriptionRepository implements the SubscriptionRepository interface using SQLite.
type SQLiteSubscriptionRepository struct {
	db *sql.DB
}

// NewSQLiteSubscriptionRepository creates a new instance of SQLiteSubscriptionRepository.
func NewSQLiteSubscriptionRepository(db *sql.DB) *SQLiteSubscriptionRepository {
	return &SQLiteSubscriptionRepository{db: db}
}

// Create inserts a new subscription into the database.
func (r *SQLiteSubscriptionRepository) Create(ctx context.Context, subscription *Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, unit_id, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, subscription.UserID, subscription.UnitID, subscription.Status, now, now)
	if err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	subscription.ID = id
	subscription.CreatedAt = now
	subscription.UpdatedAt = now
	return nil
}

// GetByID retrieves a subscription by its ID.
func (r *SQLiteSubscriptionRepository) GetByID(ctx context.Context, id int64) (*Subscription, error) {
	query := `
		SELECT id, user_id, unit_id, status, created_at, updated_at
		FROM subscriptions
		WHERE id = ?
	`

	subscription := &Subscription{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.UnitID,
		&subscription.Status,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Subscription not found
		}
		return nil, fmt.Errorf("failed to retrieve subscription by id: %w", err)
	}

	return subscription, nil
}

// Update updates an existing subscription in the database.
func (r *SQLiteSubscriptionRepository) Update(ctx context.Context, subscription *Subscription) error {
	query := `
		UPDATE subscriptions
		SET user_id = ?, unit_id = ?, status = ?, updated_at = ?
		WHERE id = ?
	`
	subscription.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query, subscription.UserID, subscription.UnitID, subscription.Status, subscription.UpdatedAt, subscription.ID)
	if err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	return nil
}

// Delete deletes a subscription from the database.
func (r *SQLiteSubscriptionRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	return nil
}

// GetByUserID retrieves all subscriptions for a specific user.
func (r *SQLiteSubscriptionRepository) GetByUserID(ctx context.Context, userID int64) ([]*Subscription, error) {
	query := `
		SELECT id, user_id, unit_id, status, created_at, updated_at
		FROM subscriptions
		WHERE user_id = ?
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve subscriptions by user id: %w", err)
	}
	defer rows.Close()

	var subscriptions []*Subscription
	for rows.Next() {
		subscription := &Subscription{}
		if err := rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.UnitID,
			&subscription.Status,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan subscription: %w", err)
		}
		subscriptions = append(subscriptions, subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	return subscriptions, nil
}