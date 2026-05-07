package db

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations применяет SQL миграции из каталога migrationsDir.
func RunMigrations(ctx context.Context, db *sql.DB, migrationsDir string) error {
	return RunMigrationsFor(ctx, db, DBPostgres, migrationsDir)
}

// RunMigrationsFor применяет SQL миграции из каталога migrationsDir.
func RunMigrationsFor(ctx context.Context, db *sql.DB, dbType DBType, migrationsDir string) error {
	if err := ensureMigrationsTable(ctx, db, dbType); err != nil {
		return err
	}

	files, err := sqlFiles(migrationsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		applied, err := isMigrationApplied(ctx, db, dbType, file)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := os.ReadFile(filepath.Join(migrationsDir, file))
		if err != nil {
			return fmt.Errorf("чтение миграции %s: %w", file, err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if _, err = tx.ExecContext(ctx, string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("ошибка выполнения миграции %s: %w", file, err)
		}

		if _, err = tx.ExecContext(ctx, insertHistorySQL(dbType, "schema_migrations"), file); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("ошибка записи schema_migrations для %s: %w", file, err)
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func ensureMigrationsTable(ctx context.Context, db *sql.DB, dbType DBType) error {
	query := `
CREATE TABLE IF NOT EXISTS schema_migrations (
    id BIGSERIAL PRIMARY KEY,
    filename TEXT NOT NULL UNIQUE,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`
	if dbType == DBMySQL {
		query = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    filename VARCHAR(255) NOT NULL UNIQUE,
    applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	}
	_, err := db.ExecContext(ctx, query)
	return err
}

func isMigrationApplied(ctx context.Context, db *sql.DB, dbType DBType, filename string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, existsSQL("schema_migrations", dbType), filename).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

// RunSeeds применяет SQL seed-файлы из каталога seedsDir один раз.
func RunSeeds(ctx context.Context, db *sql.DB, seedsDir string) error {
	return RunSeedsFor(ctx, db, DBPostgres, seedsDir)
}

// RunSeedsFor применяет SQL seed-файлы из каталога seedsDir.
func RunSeedsFor(ctx context.Context, db *sql.DB, dbType DBType, seedsDir string) error {
	if err := ensureSeedHistoryTable(ctx, db, dbType); err != nil {
		return err
	}

	files, err := sqlFiles(seedsDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		applied, err := isSeedApplied(ctx, db, dbType, file)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		content, err := os.ReadFile(filepath.Join(seedsDir, file))
		if err != nil {
			return fmt.Errorf("чтение seed %s: %w", file, err)
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if _, err = tx.ExecContext(ctx, string(content)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("ошибка выполнения seed %s: %w", file, err)
		}

		if _, err = tx.ExecContext(ctx, insertHistorySQL(dbType, "seed_history"), file); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("ошибка записи seed_history для %s: %w", file, err)
		}

		if err = tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func ensureSeedHistoryTable(ctx context.Context, db *sql.DB, dbType DBType) error {
	query := `
CREATE TABLE IF NOT EXISTS seed_history (
    id BIGSERIAL PRIMARY KEY,
    filename TEXT NOT NULL UNIQUE,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
);`
	if dbType == DBMySQL {
		query = `
CREATE TABLE IF NOT EXISTS seed_history (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    filename VARCHAR(255) NOT NULL UNIQUE,
    applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`
	}
	_, err := db.ExecContext(ctx, query)
	return err
}

func isSeedApplied(ctx context.Context, db *sql.DB, dbType DBType, filename string) (bool, error) {
	var exists bool
	err := db.QueryRowContext(ctx, existsSQL("seed_history", dbType), filename).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func sqlFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("чтение каталога %s: %w", dir, err)
	}

	files := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasSuffix(strings.ToLower(name), ".sql") {
			files = append(files, name)
		}
	}

	sort.Strings(files)
	return files, nil
}

func existsSQL(table string, dbType DBType) string {
	if dbType == DBMySQL {
		return fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE filename = ?)", table)
	}
	return fmt.Sprintf("SELECT EXISTS(SELECT 1 FROM %s WHERE filename = $1)", table)
}

func insertHistorySQL(dbType DBType, table string) string {
	if dbType == DBMySQL {
		return fmt.Sprintf("INSERT INTO %s (filename) VALUES (?)", table)
	}
	return fmt.Sprintf("INSERT INTO %s (filename) VALUES ($1)", table)
}
