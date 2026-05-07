package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"bd_lab_3/internal/config"
	"bd_lab_3/internal/db"
	"bd_lab_3/internal/handlers"
)

func main() {
	logger := log.New(os.Stdout, "[fleet] ", log.LstdFlags|log.Lshortfile)

	cfg, err := config.Load()
	if err != nil {
		logger.Fatalf("ошибка конфигурации: %v", err)
	}

	cmd := "serve"
	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	postgresDB, err := db.OpenWithRetry(ctx, "pgx", cfg.PostgresDSN(), 25, 2*time.Second)
	if err != nil {
		logger.Fatalf("ошибка подключения к PostgreSQL: %v", err)
	}
	if err := ensureMySQLDatabase(ctx, cfg, logger); err != nil {
		logger.Fatalf("ошибка подготовки базы MySQL: %v", err)
	}
	mysqlDB, err := db.OpenWithRetry(ctx, "mysql", cfg.MySQLDSN(), 25, 2*time.Second)
	if err != nil {
		logger.Fatalf("ошибка подключения к MySQL: %v", err)
	}
	dbManager := db.NewManager(db.ParseDBType(cfg.DefaultDB, db.DBPostgres), postgresDB, mysqlDB)
	defer func() {
		if closeErr := dbManager.Close(); closeErr != nil {
			logger.Printf("ошибка закрытия БД: %v", closeErr)
		}
	}()

	switch cmd {
	case "migrate":
		if err := runAllMigrations(ctx, postgresDB, mysqlDB); err != nil {
			logger.Fatalf("ошибка миграций: %v", err)
		}
		logger.Println("миграции успешно применены")
		return
	case "seed":
		if err := runAllMigrations(ctx, postgresDB, mysqlDB); err != nil {
			logger.Fatalf("ошибка миграций перед сидированием: %v", err)
		}
		if err := runAllSeeds(ctx, postgresDB, mysqlDB); err != nil {
			logger.Fatalf("ошибка seed-данных: %v", err)
		}
		logger.Println("seed-данные успешно загружены")
		return
	case "serve":
		// continue below
	default:
		logger.Fatalf("неизвестная команда: %s (доступно: serve, migrate, seed)", cmd)
	}

	if cfg.AutoMigrate {
		logger.Println("применение миграций...")
		if err := runAllMigrations(ctx, postgresDB, mysqlDB); err != nil {
			logger.Fatalf("ошибка миграций: %v", err)
		}
	}

	if cfg.AutoSeed {
		logger.Println("загрузка seed-данных...")
		if err := runAllSeeds(ctx, postgresDB, mysqlDB); err != nil {
			logger.Fatalf("ошибка seed-данных: %v", err)
		}
	}

	app, err := handlers.NewApp(cfg, logger, dbManager)
	if err != nil {
		logger.Fatalf("ошибка инициализации приложения: %v", err)
	}

	server := &http.Server{
		Addr:         cfg.Address(),
		Handler:      app.Routes(),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		<-ctx.Done()

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.Println("получен сигнал завершения, останавливаем HTTP сервер...")
		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Printf("ошибка штатного завершения сервера: %v", err)
		}
	}()

	logger.Printf("сервер запущен на http://%s", cfg.Address())
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalf("ошибка HTTP сервера: %v", err)
	}

	logger.Println("сервер остановлен")
}

func runAllMigrations(ctx context.Context, postgresDB, mysqlDB *sql.DB) error {
	postgresDir := "migrations/postgres"
	if _, err := os.Stat(postgresDir); err != nil {
		postgresDir = "migrations"
	}
	if err := db.RunMigrationsFor(ctx, postgresDB, db.DBPostgres, postgresDir); err != nil {
		return err
	}
	return db.RunMigrationsFor(ctx, mysqlDB, db.DBMySQL, "migrations/mysql")
}

func runAllSeeds(ctx context.Context, postgresDB, mysqlDB *sql.DB) error {
	postgresDir := "seeds/postgres"
	if _, err := os.Stat(postgresDir); err != nil {
		postgresDir = "seeds"
	}
	if err := db.RunSeedsFor(ctx, postgresDB, db.DBPostgres, postgresDir); err != nil {
		return err
	}
	return db.RunSeedsFor(ctx, mysqlDB, db.DBMySQL, "seeds/mysql")
}

func ensureMySQLDatabase(ctx context.Context, cfg *config.Config, logger *log.Logger) error {
	if cfg.MySQLURL != "" {
		return nil
	}

	serverDB, err := db.OpenWithRetry(ctx, "mysql", cfg.MySQLServerDSN(), 25, 2*time.Second)
	if err != nil {
		return err
	}
	defer func() {
		if closeErr := serverDB.Close(); closeErr != nil {
			logger.Printf("ошибка закрытия служебного подключения MySQL: %v", closeErr)
		}
	}()

	dbName := strings.ReplaceAll(cfg.MySQLDatabase, "`", "``")
	query := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci", dbName)
	if _, err := serverDB.ExecContext(ctx, query); err != nil {
		return err
	}
	return nil
}
