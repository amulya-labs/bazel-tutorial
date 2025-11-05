package main

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// SQLiteStore manages the SQLite database for economic indicators
type SQLiteStore struct {
	db *sql.DB
}

// NewSQLiteStore creates or opens a SQLite database and initializes schema
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	// Ensure directory exists
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Open database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	store := &SQLiteStore{db: db}

	// Initialize schema
	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return store, nil
}

// Close closes the database connection
func (s *SQLiteStore) Close() error {
	return s.db.Close()
}

// initSchema creates the necessary tables if they don't exist
func (s *SQLiteStore) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS series (
		id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		unit TEXT,
		source TEXT DEFAULT 'FRED'
	);

	CREATE TABLE IF NOT EXISTS observations (
		series_id TEXT NOT NULL,
		date TEXT NOT NULL,
		value REAL NOT NULL,
		PRIMARY KEY (series_id, date),
		FOREIGN KEY (series_id) REFERENCES series(id)
	);

	CREATE INDEX IF NOT EXISTS idx_observations_date 
		ON observations(date);

	CREATE TABLE IF NOT EXISTS refresh_log (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		source TEXT NOT NULL,
		started_at TEXT NOT NULL,
		finished_at TEXT,
		ok INTEGER DEFAULT 0,
		message TEXT
	);
	`

	_, err := s.db.Exec(schema)
	return err
}

// StoreSeries stores or updates series metadata
func (s *SQLiteStore) StoreSeries(series *Series) error {
	query := `
	INSERT INTO series (id, name, unit, source)
	VALUES (?, ?, ?, 'FRED')
	ON CONFLICT(id) DO UPDATE SET
		name = excluded.name,
		unit = excluded.unit
	`

	_, err := s.db.Exec(query, series.ID, series.Title, series.Units)
	return err
}

// StoreObservations stores observations for a series, replacing existing ones
func (s *SQLiteStore) StoreObservations(seriesID string, observations []Observation) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Delete existing observations for this series
	if _, err := tx.Exec("DELETE FROM observations WHERE series_id = ?", seriesID); err != nil {
		return fmt.Errorf("failed to delete old observations: %w", err)
	}

	// Insert new observations
	stmt, err := tx.Prepare("INSERT INTO observations (series_id, date, value) VALUES (?, ?, ?)")
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	for _, obs := range observations {
		// Parse value as float
		var value float64
		if _, err := fmt.Sscanf(obs.Value, "%f", &value); err != nil {
			continue // Skip invalid values
		}

		if _, err := stmt.Exec(seriesID, obs.Date, value); err != nil {
			return fmt.Errorf("failed to insert observation: %w", err)
		}
	}

	return tx.Commit()
}

// LogRefreshStart logs the start of a refresh operation
func (s *SQLiteStore) LogRefreshStart(source string) (int64, error) {
	result, err := s.db.Exec(
		"INSERT INTO refresh_log (source, started_at) VALUES (?, ?)",
		source,
		time.Now().UTC().Format(time.RFC3339),
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}

// LogRefreshEnd logs the completion of a refresh operation
func (s *SQLiteStore) LogRefreshEnd(id int64, ok bool, message string) error {
	okInt := 0
	if ok {
		okInt = 1
	}

	_, err := s.db.Exec(
		"UPDATE refresh_log SET finished_at = ?, ok = ?, message = ? WHERE id = ?",
		time.Now().UTC().Format(time.RFC3339),
		okInt,
		message,
		id,
	)
	return err
}

// GetLatestObservation retrieves the most recent observation for a series
func (s *SQLiteStore) GetLatestObservation(seriesID string) (string, float64, error) {
	var date string
	var value float64

	err := s.db.QueryRow(`
		SELECT date, value 
		FROM observations 
		WHERE series_id = ?
		ORDER BY date DESC
		LIMIT 1
	`, seriesID).Scan(&date, &value)

	if err != nil {
		return "", 0, err
	}

	return date, value, nil
}
