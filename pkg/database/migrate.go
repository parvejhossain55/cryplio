package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrator wraps golang-migrate for easier use
type Migrator struct {
	m             *migrate.Migrate
	migrationsDir string
}

// NewMigrator creates a new migrator instance from a config
func NewMigrator(cfg *Config, migrationsDir string) (*Migrator, error) {
	absPath, err := filepath.Abs(migrationsDir)
	if err != nil {
		absPath = migrationsDir
	}

	// Ensure migrations directory exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("migrations directory does not exist: %s", absPath)
	}

	// Build source URL (file://)
	sourceURL := fmt.Sprintf("file://%s", absPath)

	// Build database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode)

	// Create migrator instance
	m, err := migrate.New(sourceURL, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrator: %w", err)
	}

	return &Migrator{m: m, migrationsDir: absPath}, nil
}

// Apply runs all pending migrations
func (m *Migrator) Apply() error {
	if err := m.m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}

// Down rolls back the last migration
func (m *Migrator) Down() error {
	if err := m.m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("rollback failed: %w", err)
	}
	return nil
}

// Steps migrates N steps (positive = up, negative = down)
func (m *Migrator) Steps(n int) error {
	if err := m.m.Steps(n); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("step migration failed: %w", err)
	}
	return nil
}

// AppliedVersionsDB queries the database for applied migration versions.
// golang-migrate stores the current version as a single row, so all lower
// contiguous versions are considered applied.
func AppliedVersionsDB(db *sql.DB) ([]int, error) {
	var version int
	var dirty bool
	err := db.QueryRow("SELECT version, dirty FROM schema_migrations LIMIT 1").Scan(&version, &dirty)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("query applied version: %w", err)
	}
	if dirty {
		return nil, fmt.Errorf("database migration version %d is dirty", version)
	}

	versions := make([]int, 0, version+1)
	for v := 0; v <= version; v++ {
		versions = append(versions, v)
	}

	return versions, nil
}
