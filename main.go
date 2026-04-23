package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
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

	database, err := db.OpenWithRetry(ctx, cfg.DSN(), 25, 2*time.Second)
	if err != nil {
		logger.Fatalf("ошибка подключения к БД: %v", err)
	}
	defer func() {
		if closeErr := database.Close(); closeErr != nil {
			logger.Printf("ошибка закрытия БД: %v", closeErr)
		}
	}()

	switch cmd {
	case "migrate":
		if err := db.RunMigrations(ctx, database, "migrations"); err != nil {
			logger.Fatalf("ошибка миграций: %v", err)
		}
		logger.Println("миграции успешно применены")
		return
	case "seed":
		if err := db.RunMigrations(ctx, database, "migrations"); err != nil {
			logger.Fatalf("ошибка миграций перед сидированием: %v", err)
		}
		if err := db.RunSeeds(ctx, database, "seeds"); err != nil {
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
		if err := db.RunMigrations(ctx, database, "migrations"); err != nil {
			logger.Fatalf("ошибка миграций: %v", err)
		}
	}

	if cfg.AutoSeed {
		logger.Println("загрузка seed-данных...")
		if err := db.RunSeeds(ctx, database, "seeds"); err != nil {
			logger.Fatalf("ошибка seed-данных: %v", err)
		}
	}

	app, err := handlers.NewApp(cfg, logger, database)
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
