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

	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresSSLMode  string
	PostgresURL      string

	MySQLHost      string
	MySQLPort      string
	MySQLUser      string
	MySQLPassword  string
	MySQLDatabase  string
	MySQLParseTime string
	MySQLLoc       string
	MySQLURL       string

	DefaultDB string

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
		DBUser:      getEnv("DB_USER", "fleet_user"),
		DBPassword:  getEnv("DB_PASSWORD", "fleet_password"),
		DBName:      getEnv("DB_NAME", "fleet_db"),
		DBSSLMode:   getEnv("DB_SSLMODE", "disable"),
		DatabaseURL: getEnv("DATABASE_URL", ""),

		PostgresHost:     getEnv("POSTGRES_HOST", getEnv("DB_HOST", "localhost")),
		PostgresPort:     getEnv("POSTGRES_PORT", getEnv("DB_PORT", "5432")),
		PostgresUser:     getEnv("POSTGRES_USER", getEnv("DB_USER", "fleet_user")),
		PostgresPassword: getEnv("POSTGRES_PASSWORD", getEnv("DB_PASSWORD", "fleet_password")),
		PostgresDB:       getEnv("POSTGRES_DB", getEnv("DB_NAME", "fleet_db")),
		PostgresSSLMode:  getEnv("POSTGRES_SSLMODE", getEnv("DB_SSLMODE", "disable")),
		PostgresURL:      getEnv("POSTGRES_URL", getEnv("DATABASE_URL", "")),

		MySQLHost:      getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:      getEnv("MYSQL_PORT", "3306"),
		MySQLUser:      getEnv("MYSQL_USER", "fleet_user"),
		MySQLPassword:  getEnv("MYSQL_PASSWORD", "fleet_password"),
		MySQLDatabase:  getEnv("MYSQL_DATABASE", "fleet_db"),
		MySQLParseTime: getEnv("MYSQL_PARSE_TIME", "true"),
		MySQLLoc:       getEnv("MYSQL_LOC", "Local"),
		MySQLURL:       getEnv("MYSQL_URL", ""),

		DefaultDB: getEnv("DEFAULT_DB", "postgres"),

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

// DSN возвращает строку подключения к PostgreSQL (legacy helper).
func (c *Config) DSN() string {
	return c.PostgresDSN()
}

// PostgresDSN возвращает строку подключения к PostgreSQL.
func (c *Config) PostgresDSN() string {
	if c.PostgresURL != "" {
		return c.PostgresURL
	}
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.PostgresHost,
		c.PostgresPort,
		c.PostgresUser,
		c.PostgresPassword,
		c.PostgresDB,
		c.PostgresSSLMode,
	)
}

// MySQLDSN возвращает строку подключения к MySQL.
func (c *Config) MySQLDSN() string {
	if c.MySQLURL != "" {
		return c.MySQLURL
	}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=%s&loc=%s&multiStatements=true&charset=utf8mb4,utf8",
		c.MySQLUser,
		c.MySQLPassword,
		c.MySQLHost,
		c.MySQLPort,
		c.MySQLDatabase,
		c.MySQLParseTime,
		c.MySQLLoc,
	)
}

// MySQLServerDSN возвращает строку подключения к серверу MySQL без выбора БД.
func (c *Config) MySQLServerDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?parseTime=%s&loc=%s&multiStatements=true&charset=utf8mb4,utf8",
		c.MySQLUser,
		c.MySQLPassword,
		c.MySQLHost,
		c.MySQLPort,
		c.MySQLParseTime,
		c.MySQLLoc,
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
