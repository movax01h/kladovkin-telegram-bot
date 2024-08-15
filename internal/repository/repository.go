package repository

import (
	"database/sql"
)

// User represents the data structure for a user in the system.
type User struct {
	ID         int64
	Name       string
	Email      string
	TelegramID int64
	CreatedAt  string
	UpdatedAt  string
}

// Unit represents the data structure for a unit in the system.
type Unit struct {
	ID          int64
	Name        string
	City        string
	Size        string
	Dimension   string
	Price       float64
	Available   bool
	Description string
	CreatedAt   string
	UpdatedAt   string
}

// Subscription represents the data structure for a subscription in the system.
type Subscription struct {
	ID        int64
	UserID    int64
	UnitID    int64
	Status    string
	CreatedAt string
	UpdatedAt string
}

// UserRepository defines the methods to interact with the user data.
type UserRepository interface {
	SaveUser(user User) error
	GetUserByID(id int64) (*User, error)
	GetAllUsers() ([]User, error)
	DeleteUser(id int64) error
}

// UnitRepository defines the methods to interact with the unit data.
type UnitRepository interface {
	SaveUnit(unit Unit) error
	GetUnitByID(id int64) (*Unit, error)
	GetAllUnits() ([]Unit, error)
	DeleteUnit(id int64) error
}

// SubscriptionRepository defines the methods to interact with the subscription data.
type SubscriptionRepository interface {
	SaveSubscription(subscription Subscription) error
	GetSubscriptionByID(id int64) (*Subscription, error)
	GetAllSubscriptions() ([]Subscription, error)
	DeleteSubscription(id int64) error
}

// SQLiteUserRepository is the SQLite implementation of the UserRepository.
type SQLiteUserRepository struct {
	db *sql.DB
}

// NewSQLiteUserRepository creates a new SQLiteUserRepository.
func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

// SaveUser saves a user in the database.
func (r *SQLiteUserRepository) SaveUser(user User) error {
	_, err := r.db.Exec("INSERT INTO users (name, email, created_at, updated_at) VALUES (?, ?, ?, ?)",
		user.Name, user.Email, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByID retrieves a user by ID from the database.
func (r *SQLiteUserRepository) GetUserByID(id int64) (*User, error) {
	row := r.db.QueryRow("SELECT id, name, email, created_at, updated_at FROM users WHERE id = ?", id)
	var user User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetAllUsers retrieves all users from the database.
func (r *SQLiteUserRepository) GetAllUsers() ([]User, error) {
	rows, err := r.db.Query("SELECT id, name, email, created_at, updated_at FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// DeleteUser deletes a user from the database.
func (r *SQLiteUserRepository) DeleteUser(id int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	return err
}

// SQLiteUnitRepository is the SQLite implementation of the UnitRepository.
type SQLiteUnitRepository struct {
	db *sql.DB
}

// NewSQLiteUnitRepository creates a new SQLiteUnitRepository.
func NewSQLiteUnitRepository(db *sql.DB) *SQLiteUnitRepository {
	return &SQLiteUnitRepository{db: db}
}

// SaveUnit saves a unit in the database.
func (r *SQLiteUnitRepository) SaveUnit(unit Unit) error {
	_, err := r.db.Exec("INSERT INTO units (name, description, created_at, updated_at) VALUES (?, ?, ?, ?)",
		unit.Name, unit.Description, unit.CreatedAt, unit.UpdatedAt)
	return err
}

// GetUnitByID retrieves a unit by ID from the database.
func (r *SQLiteUnitRepository) GetUnitByID(id int64) (*Unit, error) {
	row := r.db.QueryRow("SELECT id, name, description, created_at, updated_at FROM units WHERE id = ?", id)
	var unit Unit
	err := row.Scan(&unit.ID, &unit.Name, &unit.Description, &unit.CreatedAt, &unit.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &unit, nil
}

// GetAllUnits retrieves all units from the database.
func (r *SQLiteUnitRepository) GetAllUnits() ([]Unit, error) {
	rows, err := r.db.Query("SELECT id, name, description, created_at, updated_at FROM units")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var units []Unit
	for rows.Next() {
		var unit Unit
		if err := rows.Scan(&unit.ID, &unit.Name, &unit.Description, &unit.CreatedAt, &unit.UpdatedAt); err != nil {
			return nil, err
		}
		units = append(units, unit)
	}
	return units, nil
}

// DeleteUnit deletes a unit from the database.
func (r *SQLiteUnitRepository) DeleteUnit(id int64) error {
	_, err := r.db.Exec("DELETE FROM units WHERE id = ?", id)
	return err
}

// SQLiteSubscriptionRepository is the SQLite implementation of the SubscriptionRepository.
type SQLiteSubscriptionRepository struct {
	db *sql.DB
}

// NewSQLiteSubscriptionRepository creates a new SQLiteSubscriptionRepository.
func NewSQLiteSubscriptionRepository(db *sql.DB) *SQLiteSubscriptionRepository {
	return &SQLiteSubscriptionRepository{db: db}
}

// SaveSubscription saves a subscription in the database.
func (r *SQLiteSubscriptionRepository) SaveSubscription(subscription Subscription) error {
	_, err := r.db.Exec("INSERT INTO subscriptions (user_id, unit_id, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		subscription.UserID, subscription.UnitID, subscription.Status, subscription.CreatedAt, subscription.UpdatedAt)
	return err
}

// GetSubscriptionByID retrieves a subscription by ID from the database.
func (r *SQLiteSubscriptionRepository) GetSubscriptionByID(id int64) (*Subscription, error) {
	row := r.db.QueryRow("SELECT id, user_id, unit_id, status, created_at, updated_at FROM subscriptions WHERE id = ?", id)
	var subscription Subscription
	err := row.Scan(&subscription.ID, &subscription.UserID, &subscription.UnitID, &subscription.Status, &subscription.CreatedAt, &subscription.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &subscription, nil
}

// GetAllSubscriptions retrieves all subscriptions from the database.
func (r *SQLiteSubscriptionRepository) GetAllSubscriptions() ([]Subscription, error) {
	rows, err := r.db.Query("SELECT id, user_id, unit_id, status, created_at, updated_at FROM subscriptions")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []Subscription
	for rows.Next() {
		var subscription Subscription
		if err := rows.Scan(&subscription.ID, &subscription.UserID, &subscription.UnitID, &subscription.Status, &subscription.CreatedAt, &subscription.UpdatedAt); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}
	return subscriptions, nil
}

// DeleteSubscription deletes a subscription from the database.
func (r *SQLiteSubscriptionRepository) DeleteSubscription(id int64) error {
	_, err := r.db.Exec("DELETE FROM subscriptions WHERE id = ?", id)
	return err
}
