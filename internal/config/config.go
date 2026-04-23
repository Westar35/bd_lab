package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config хранит параметры приложения и БД.
type Config struct {
	AppHost       string
	AppPort       string
	AppEnv        string
	ReadTimeout   time.Duration
	WriteTimeout  time.Duration
	SessionKey    string
	AdminUsername string
	AdminPassword string

	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	DatabaseURL string

	AutoMigrate bool
	AutoSeed    bool
}

// Load загружает конфигурацию из переменных окружения и .env.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppHost:       getEnv("APP_HOST", "0.0.0.0"),
		AppPort:       getEnv("APP_PORT", "8080"),
		AppEnv:        getEnv("APP_ENV", "development"),
		ReadTimeout:   time.Duration(getEnvAsInt("READ_TIMEOUT_SEC", 10)) * time.Second,
		WriteTimeout:  time.Duration(getEnvAsInt("WRITE_TIMEOUT_SEC", 20)) * time.Second,
		SessionKey:    getEnv("SESSION_KEY", "change_me_for_production_very_secret_key"),
		AdminUsername: getEnv("ADMIN_USERNAME", "admin"),
		AdminPassword: getEnv("ADMIN_PASSWORD", "admin123"),

		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "5432"),
		DBUser:      getEnv("DB_USER", "postgres"),
		DBPassword:  getEnv("DB_PASSWORD", "Fgrths197+"),
		DBName:      getEnv("DB_NAME", "fleet_db"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		DatabaseURL: getEnv("DATABASE_URL", ""),

		AutoMigrate: getEnvAsBool("AUTO_MIGRATE", true),
		AutoSeed:    getEnvAsBool("AUTO_SEED", false),
	}

	if len(cfg.SessionKey) < 16 {
		return nil, fmt.Errorf("SESSION_KEY должен быть длиной не менее 16 символов")
	}
	if cfg.AdminUsername == "" || cfg.AdminPassword == "" {
		return nil, fmt.Errorf("ADMIN_USERNAME и ADMIN_PASSWORD обязательны")
	}

	return cfg, nil
}

// Address возвращает адрес, на котором слушает HTTP сервер.
func (c *Config) Address() string {
	return fmt.Sprintf("%s:%s", c.AppHost, c.AppPort)
}

// DSN возвращает строку подключения к PostgreSQL.
func (c *Config) DSN() string {
	if c.DatabaseURL != "" {
		return c.DatabaseURL
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	n, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvAsBool(key string, fallback bool) bool {
	value := getEnv(key, "")
	if value == "" {
		return fallback
	}

	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return fallback
	}
	return parsed
}
