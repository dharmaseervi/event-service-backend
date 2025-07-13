package migrations

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

func RunMigrations(db *sql.DB) error {
	// Check if migrations table exists
	if err := createMigrationsTable(db); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get all migration files
	migrationFiles, err := getMigrationFiles()
	if err != nil {
		return fmt.Errorf("failed to get migration files: %w", err)
	}

	// Get completed migrations
	completedMigrations, err := getCompletedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get completed migrations: %w", err)
	}

	// Run new migrations
	for _, m := range migrationFiles {
		if _, exists := completedMigrations[m.Version]; !exists {
			if err := runMigration(db, m); err != nil {
				return fmt.Errorf("failed to run migration %d: %w", m.Version, err)
			}
		}
	}

	return nil
}

type migration struct {
	Version  int
	UpFile   string
	DownFile string
}

func getMigrationFiles() ([]migration, error) {
	var migrations []migration

	err := filepath.WalkDir("migrations", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		filename := d.Name()
		parts := strings.Split(filename, "_")
		if len(parts) < 2 {
			return nil
		}

		version, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil
		}

		if strings.HasSuffix(filename, ".up.sql") {
			migrations = append(migrations, migration{
				Version: version,
				UpFile:  path,
			})

		} else if strings.HasSuffix(filename, ".down.sql") {
			// Find corresponding up migration and add down file
			for i, m := range migrations {
				if m.Version == version && m.DownFile == "" {
					migrations[i].DownFile = path
					break
				}
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func createMigrationsTable(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	return err
}

func getCompletedMigrations(db *sql.DB) (map[int]bool, error) {
	rows, err := db.Query("SELECT version FROM migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	completed := make(map[int]bool)
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, err
		}
		completed[version] = true
	}

	return completed, nil
}

func runMigration(db *sql.DB, m migration) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Read migration file
	content, err := os.ReadFile(m.UpFile)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}

	// Execute migration
	if _, err := tx.Exec(string(content)); err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Record migration
	if _, err := tx.Exec("INSERT INTO migrations (version) VALUES ($1)", m.Version); err != nil {
		return fmt.Errorf("failed to record migration: %w", err)
	}
	log.Println("Running migration:", m.UpFile)
	return tx.Commit()
}
