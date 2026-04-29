package database

import (
	"database/sql"
	"fmt"
)

// Config holds database connection configuration
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// DefaultConfig returns default config for local development
func DefaultConfig() *Config {
	return &Config{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "",
		DBName:   "cryplio_db",
		SSLMode:  "disable",
	}
}

// ConnStr returns PostgreSQL connection string
func (c *Config) ConnStr(withDB bool) string {
	dbname := c.DBName
	if !withDB {
		dbname = "postgres"
	}
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, dbname, c.SSLMode)
}

// Open connects to the database (with DB name)
func Open(cfg *Config) (*sql.DB, error) {
	return sql.Open("postgres", cfg.ConnStr(true))
}

// OpenAdmin connects to PostgreSQL admin database (for creating DB)
func OpenAdmin(cfg *Config) (*sql.DB, error) {
	return sql.Open("postgres", cfg.ConnStr(false))
}

// DatabaseExists checks if the database exists
func DatabaseExists(cfg *Config) (bool, error) {
	db, err := OpenAdmin(cfg)
	if err != nil {
		return false, err
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)", cfg.DBName).Scan(&exists)
	return exists, err
}

// CreateDatabase creates the database if it doesn't exist
func CreateDatabase(cfg *Config) error {
	exists, err := DatabaseExists(cfg)
	if err != nil {
		return err
	}
	if exists {
		return nil
	}

	db, err := OpenAdmin(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	return err
}

// EnsureDatabase ensures database exists (creates if needed)
func EnsureDatabase(cfg *Config) error {
	return CreateDatabase(cfg)
}
