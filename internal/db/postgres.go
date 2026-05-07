package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// OpenWithRetry подключается к БД с несколькими попытками.
func OpenWithRetry(ctx context.Context, driverName, dsn string, maxAttempts int, delay time.Duration) (*sql.DB, error) {
	var lastErr error

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		db, err := sql.Open(driverName, dsn)
		if err != nil {
			lastErr = err
		} else {
			db.SetMaxOpenConns(25)
			db.SetMaxIdleConns(25)
			db.SetConnMaxIdleTime(5 * time.Minute)
			db.SetConnMaxLifetime(30 * time.Minute)

			pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			pingErr := db.PingContext(pingCtx)
			cancel()
			if pingErr == nil {
				return db, nil
			}
			_ = db.Close()
			lastErr = pingErr
		}

		log.Printf("[db] попытка %d/%d неуспешна: %v", attempt, maxAttempts, lastErr)

		if attempt < maxAttempts {
			select {
			case <-time.After(delay):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}

	return nil, fmt.Errorf("не удалось подключиться к БД после %d попыток: %w", maxAttempts, lastErr)
}
