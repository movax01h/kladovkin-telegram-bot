package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

// Unit represents a unit that users can subscribe to.
type Unit struct {
	ID        int64
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// SQLiteUnitRepository implements the UnitRepository interface using SQLite.
type SQLiteUnitRepository struct {
	db *sql.DB
}

// NewSQLiteUnitRepository creates a new instance of SQLiteUnitRepository.
func NewSQLiteUnitRepository(db *sql.DB) *SQLiteUnitRepository {
	return &SQLiteUnitRepository{db: db}
}

// Create inserts a new unit into the database.
func (r *SQLiteUnitRepository) Create(ctx context.Context, unit *Unit) error {
	query := `
		INSERT INTO units (name, created_at, updated_at)
		VALUES (?, ?, ?)
	`
	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, unit.Name, now, now)
	if err != nil {
		return fmt.Errorf("failed to create unit: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	unit.ID = id
	unit.CreatedAt = now
	unit.UpdatedAt = now
	return nil
}

// GetByID retrieves a unit by its ID.
func (r *SQLiteUnitRepository) GetByID(ctx context.Context, id int64) (*Unit, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM units
		WHERE id = ?
	`

	unit := &Unit{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&unit.ID,
		&unit.Name,
		&unit.CreatedAt,
		&unit.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Unit not found
		}
		return nil, fmt.Errorf("failed to retrieve unit by id: %w", err)
	}

	return unit, nil
}

// Update updates an existing unit in the database.
func (r *SQLiteUnitRepository) Update(ctx context.Context, unit *Unit) error {
	query := `
		UPDATE units
		SET name = ?, updated_at = ?
		WHERE id = ?
	`
	unit.UpdatedAt = time.Now()

	_, err := r.db.ExecContext(ctx, query, unit.Name, unit.UpdatedAt, unit.ID)
	if err != nil {
		return fmt.Errorf("failed to update unit: %w", err)
	}

	return nil
}

// Delete deletes a unit from the database.
func (r *SQLiteUnitRepository) Delete(ctx context.Context, id int64) error {
	query := `
		DELETE FROM units
		WHERE id = ?
	`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete unit: %w", err)
	}

	return nil
}

// GetAll retrieves all units from the database.
func (r *SQLiteUnitRepository) GetAll(ctx context.Context) ([]*Unit, error) {
	query := `
		SELECT id, name, created_at, updated_at
		FROM units
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve all units: %w", err)
	}
	defer rows.Close()

	var units []*Unit
	for rows.Next() {
		unit := &Unit{}
		if err := rows.Scan(
			&unit.ID,
			&unit.Name,
			&unit.CreatedAt,
			&unit.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan unit: %w", err)
		}
		units = append(units, unit)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error occurred during rows iteration: %w", err)
	}

	return units, nil
}
