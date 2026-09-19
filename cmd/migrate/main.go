package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"

	"github.com/KAZI-CODE-HIJABI/Backend/internal/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var migrationName = regexp.MustCompile(`^(\d+)_.+\.(up|down)\.sql$`)

type migration struct {
	version int64
	path    string
}

func main() {
	if len(os.Args) != 2 || (os.Args[1] != "up" && os.Args[1] != "down") {
		fmt.Fprintln(os.Stderr, "usage: go run ./cmd/migrate [up|down]")
		os.Exit(2)
	}
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	dsn := cfg.MigrationDatabaseURL
	if dsn == "" {
		dsn = cfg.DatabaseURL
	}
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "connect to database:", err)
		os.Exit(1)
	}
	defer pool.Close()
	if err := pool.Ping(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, "connect to database:", err)
		os.Exit(1)
	}
	if err := run(context.Background(), pool, os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, "run migrations:", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, pool *pgxpool.Pool, direction string) error {
	if _, err := pool.Exec(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version BIGINT PRIMARY KEY, applied_at TIMESTAMPTZ NOT NULL DEFAULT now())`); err != nil {
		return err
	}
	if direction == "up" {
		return up(ctx, pool)
	}
	return down(ctx, pool)
}

func up(ctx context.Context, pool *pgxpool.Pool) error {
	migrations, err := find("up")
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		var applied bool
		if err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM schema_migrations WHERE version = $1)`, migration.version).Scan(&applied); err != nil {
			return err
		}
		if applied {
			continue
		}
		sql, err := os.ReadFile(migration.path)
		if err != nil {
			return err
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, string(sql)); err == nil {
			_, err = tx.Exec(ctx, `INSERT INTO schema_migrations (version) VALUES ($1)`, migration.version)
		}
		if err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("apply %s: %w", filepath.Base(migration.path), err)
		}
		if err := tx.Commit(ctx); err != nil {
			return err
		}
	}
	return nil
}

func down(ctx context.Context, pool *pgxpool.Pool) error {
	var version int64
	if err := pool.QueryRow(ctx, `SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 1`).Scan(&version); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	}
	migrations, err := find("down")
	if err != nil {
		return err
	}
	for _, migration := range migrations {
		if migration.version != version {
			continue
		}
		sql, err := os.ReadFile(migration.path)
		if err != nil {
			return err
		}
		tx, err := pool.Begin(ctx)
		if err != nil {
			return err
		}
		if _, err = tx.Exec(ctx, string(sql)); err == nil {
			_, err = tx.Exec(ctx, `DELETE FROM schema_migrations WHERE version = $1`, version)
		}
		if err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("reverse %s: %w", filepath.Base(migration.path), err)
		}
		return tx.Commit(ctx)
	}
	return fmt.Errorf("missing down migration for version %d", version)
}

func find(direction string) ([]migration, error) {
	entries, err := os.ReadDir("migrations")
	if err != nil {
		return nil, err
	}
	var migrations []migration
	for _, entry := range entries {
		matches := migrationName.FindStringSubmatch(entry.Name())
		if entry.IsDir() || matches == nil || matches[2] != direction {
			continue
		}
		version, err := strconv.ParseInt(matches[1], 10, 64)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, migration{version: version, path: filepath.Join("migrations", entry.Name())})
	}
	sort.Slice(migrations, func(i, j int) bool { return migrations[i].version < migrations[j].version })
	return migrations, nil
}
