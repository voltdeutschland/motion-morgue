package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

func DataDir() (string, error) {
	assoc := os.Getenv("MOTION_MORGUE_ASSOC")
	if assoc == "" {
		return "", fmt.Errorf("MOTION_MORGUE_ASSOC environment variable not set")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not get home directory: %w", err)
	}

	return filepath.Join(homeDir, ".motion-morgue", assoc), nil
}

func Init() error {
	dataDir, err := DataDir()
	if err != nil {
		return err
	}

	dbPath := filepath.Join(dataDir, "data.db")
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return fmt.Errorf("could not create data directory: %w", err)
	}

	DB, err = sql.Open("sqlite", dbPath)
	if err != nil {
		return fmt.Errorf("could not open database: %w", err)
	}

	if err := migrate(); err != nil {
		return fmt.Errorf("could not run migrations: %w", err)
	}

	return nil
}

func migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS assemblies (
		id INTEGER PRIMARY KEY,
		title TEXT NOT NULL,
		start_date TEXT,
		end_date TEXT,
		protocol_pdf TEXT
	);

	CREATE TABLE IF NOT EXISTS motions (
		id INTEGER PRIMARY KEY,
		assembly_id INTEGER NOT NULL REFERENCES assemblies(id),
		title TEXT NOT NULL,
		sort_number TEXT NOT NULL,
		pdf_path TEXT
	);

	CREATE TABLE IF NOT EXISTS amendments (
		id INTEGER PRIMARY KEY,
		motion_id INTEGER NOT NULL REFERENCES motions(id),
		title TEXT,
		sort_number TEXT NOT NULL,
		pdf_path TEXT
	);
	`

	_, err := DB.Exec(schema)
	if err != nil {
		return err
	}

	// Migration: add title column if it doesn't exist
	DB.Exec("ALTER TABLE amendments ADD COLUMN title TEXT")

	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}
