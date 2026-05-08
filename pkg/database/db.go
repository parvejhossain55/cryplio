package database

import (
	"database/sql"
)

// DB wraps sql.DB for easier use
type DB struct {
	*sql.DB
}

// NewDB creates a new DB wrapper
func NewDB(sqlDB *sql.DB) *DB {
	return &DB{DB: sqlDB}
}
