package database

import (
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// DB is the global database connection
var DB *sqlx.DB

// Init initializes the database connection
func Init() error {
	// Get data directory path
	dataDir := getDataDir()
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(dataDir, "postme.db")
	var err error
	DB, err = sqlx.Open("sqlite", dbPath+"?_pragma=foreign_keys(1)")
	if err != nil {
		return err
	}

	// Run migrations
	if err := RunMigrations(DB); err != nil {
		return err
	}

	return nil
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

// getDataDir returns the data directory path
func getDataDir() string {
	// Get executable directory
	exe, err := os.Executable()
	if err != nil {
		// Fallback to current directory
		return "data"
	}
	return filepath.Join(filepath.Dir(exe), "data")
}

// GetDB returns the database connection
func GetDB() *sqlx.DB {
	return DB
}
