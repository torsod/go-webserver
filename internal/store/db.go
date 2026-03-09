package store

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB wraps a pgxpool connection pool
type DB struct {
	Pool *pgxpool.Pool
}

// New creates a new database connection pool
func New(databaseURL string) (*DB, error) {
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		return nil, fmt.Errorf("parse database URL: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		pool.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}

	slog.Info("connected to PostgreSQL")
	return &DB{Pool: pool}, nil
}

// Close closes the connection pool
func (db *DB) Close() {
	db.Pool.Close()
}

// RunMigrations applies all pending SQL migration files
func (db *DB) RunMigrations(migrationsDir string) error {
	// Create migrations tracking table
	_, err := db.Pool.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("create migrations table: %w", err)
	}

	// Get applied migrations
	rows, err := db.Pool.Query(context.Background(), "SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return fmt.Errorf("query migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var v int
		if err := rows.Scan(&v); err != nil {
			return fmt.Errorf("scan migration version: %w", err)
		}
		applied[v] = true
	}

	// Find migration files
	entries, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var upFiles []string
	for _, e := range entries {
		if strings.HasSuffix(e.Name(), ".up.sql") {
			upFiles = append(upFiles, e.Name())
		}
	}
	sort.Strings(upFiles)

	// Apply pending migrations
	for _, f := range upFiles {
		// Extract version number from filename (e.g., "001_create_users.up.sql" -> 1)
		var version int
		fmt.Sscanf(f, "%d_", &version)
		if version == 0 {
			continue
		}

		if applied[version] {
			continue
		}

		sql, err := os.ReadFile(filepath.Join(migrationsDir, f))
		if err != nil {
			return fmt.Errorf("read migration %s: %w", f, err)
		}

		tx, err := db.Pool.Begin(context.Background())
		if err != nil {
			return fmt.Errorf("begin transaction for %s: %w", f, err)
		}

		if _, err := tx.Exec(context.Background(), string(sql)); err != nil {
			tx.Rollback(context.Background())
			return fmt.Errorf("apply migration %s: %w", f, err)
		}

		if _, err := tx.Exec(context.Background(), "INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback(context.Background())
			return fmt.Errorf("record migration %s: %w", f, err)
		}

		if err := tx.Commit(context.Background()); err != nil {
			return fmt.Errorf("commit migration %s: %w", f, err)
		}

		slog.Info("applied migration", "file", f, "version", version)
	}

	return nil
}
