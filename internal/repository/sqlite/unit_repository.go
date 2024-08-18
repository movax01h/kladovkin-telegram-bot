package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	m "github.com/movax01h/kladovkin-telegram-bot/internal/models"
	"github.com/movax01h/kladovkin-telegram-bot/internal/repository"
)

var _ repository.UnitRepository = (*SQLiteUnitRepository)(nil)

// SQLiteUnitRepository implements the UnitRepository interface using SQLite.
type SQLiteUnitRepository struct {
	db *sql.DB
}

// NewSQLiteUnitRepository creates a new instance of SQLiteUnitRepository.
func NewSQLiteUnitRepository(db *sql.DB) *SQLiteUnitRepository {
	return &SQLiteUnitRepository{db: db}
}

// CreateUnit inserts or updates a unit in the database.
func (r *SQLiteUnitRepository) CreateUnit(unit *m.Unit) error {
	query := `
		INSERT INTO units (id, name, city, size, dimension, price, available, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(id) DO UPDATE SET
			name = excluded.name,
			city = excluded.city,
			size = excluded.size,
			dimension = excluded.dimension,
			price = excluded.price,
			available = excluded.available,
			description = excluded.description,
			updated_at = excluded.updated_at
	`
	_, err := r.db.Exec(
		query,
		unit.ID,
		unit.Name,
		unit.City,
		unit.Size,
		unit.Dimension,
		unit.Price,
		unit.Available,
		unit.Description,
		unit.CreatedAt,
		unit.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to save unit: %w", err)
	}
	return nil
}

// GetUnitByID retrieves a unit by ID from the database.
func (r *SQLiteUnitRepository) GetUnitByID(id int64) (*m.Unit, error) {
	query := `
		SELECT id, name, city, size, dimension, price, available, description, created_at, updated_at 
		FROM units 
		WHERE id = ?
	`
	row := r.db.QueryRow(query, id)

	var unit m.Unit
	err := row.Scan(
		&unit.ID,
		&unit.Name,
		&unit.City,
		&unit.Size,
		&unit.Dimension,
		&unit.Price,
		&unit.Available,
		&unit.Description,
		&unit.CreatedAt,
		&unit.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get unit by ID: %w", err)
	}
	return &unit, nil
}

// GetAllUnits retrieves all units from the database.
func (r *SQLiteUnitRepository) GetAllUnits() ([]*m.Unit, error) {
	query := `
		SELECT id, name, city, size, dimension, price, available, description, created_at, updated_at 
		FROM units
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all units: %w", err)
	}
	defer rows.Close()

	var units []*m.Unit
	for rows.Next() {
		var unit m.Unit
		if err := rows.Scan(
			&unit.ID,
			&unit.Name,
			&unit.City,
			&unit.Size,
			&unit.Dimension,
			&unit.Price,
			&unit.Available,
			&unit.Description,
			&unit.CreatedAt,
			&unit.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan unit row: %w", err)
		}
		units = append(units, &unit)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during rows iteration: %w", err)
	}
	return units, nil
}

// GetCities retrieves the list of cities from the database.
func (r *SQLiteUnitRepository) GetCities() ([]string, error) {
	query := `
		SELECT DISTINCT city 
		FROM units
	`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get cities: %w", err)
	}
	defer rows.Close()

	var cities []string
	for rows.Next() {
		var city string
		if err := rows.Scan(&city); err != nil {
			return nil, fmt.Errorf("failed to scan city row: %w", err)
		}
		cities = append(cities, city)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during cities rows iteration: %w", err)
	}
	return cities, nil
}

// GetStoragesByCity retrieves storage names based on the city from the database.
func (r *SQLiteUnitRepository) GetStoragesByCity(city string) ([]string, error) {
	query := `
		SELECT DISTINCT name 
		FROM units 
		WHERE city = ?
	`
	rows, err := r.db.Query(query, city)
	if err != nil {
		return nil, fmt.Errorf("failed to get storages by city: %w", err)
	}
	defer rows.Close()

	var storages []string
	for rows.Next() {
		var storage string
		if err := rows.Scan(&storage); err != nil {
			return nil, fmt.Errorf("failed to scan storage row: %w", err)
		}
		storages = append(storages, storage)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during storages rows iteration: %w", err)
	}
	return storages, nil
}

// GetUnitSizesByStorage retrieves unit sizes based on the storage name from the database.
func (r *SQLiteUnitRepository) GetUnitSizesByStorage(storage string) ([]string, error) {
	query := `
		SELECT DISTINCT size 
		FROM units 
		WHERE name = ?
	`
	rows, err := r.db.Query(query, storage)
	if err != nil {
		return nil, fmt.Errorf("failed to get unit sizes by storage: %w", err)
	}
	defer rows.Close()

	var sizes []string
	for rows.Next() {
		var size string
		if err := rows.Scan(&size); err != nil {
			return nil, fmt.Errorf("failed to scan size row: %w", err)
		}
		sizes = append(sizes, size)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during unit sizes rows iteration: %w", err)
	}
	return sizes, nil
}

// UpdateUnit updates a unit in the database.
func (r *SQLiteUnitRepository) UpdateUnit(unit *m.Unit) error {
	query := `
		UPDATE units 
		SET name = ?, city = ?, size = ?, dimension = ?, price = ?, available = ?, description = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(
		query,
		unit.Name,
		unit.City,
		unit.Size,
		unit.Dimension,
		unit.Price,
		unit.Available,
		unit.Description,
		time.Now(),
		unit.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update unit: %w", err)
	}
	return nil
}

// DeleteUnit deletes a unit from the database.
func (r *SQLiteUnitRepository) DeleteUnit(id int64) error {
	_, err := r.db.Exec("DELETE FROM units WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete unit: %w", err)
	}
	return nil
}
