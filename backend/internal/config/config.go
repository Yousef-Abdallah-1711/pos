package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Email    EmailConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host            string
	Port            int
	User            string
	Password        string
	DBName          string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxIdleTime string
	ConnMaxLifetime string
}

type JWTConfig struct {
	SuperAdminSecret string
	Secret           string
	Expiry           string
}

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
}

func (c EmailConfig) SMTPConfigured() bool {
	return c.SMTPHost != "" || c.SMTPUsername != "" || c.SMTPPassword != "" || c.FromEmail != ""
}

func Load() (*Config, error) {
	// Load .env file
	_ = godotenv.Load()

	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	emailConfig, err := loadEmailConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            dbPort,
			User:            getEnv("DB_USER", "postgres"),
			Password:        getEnv("DB_PASSWORD", "postgres"),
			DBName:          getEnv("DB_NAME", "pos_saas"),
			SSLMode:         getEnv("DB_SSLMODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 75),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 25),
			ConnMaxIdleTime: getEnv("DB_CONN_MAX_IDLE_TIME", "2m"),
			ConnMaxLifetime: getEnv("DB_CONN_MAX_LIFETIME", "5m"),
		},
		JWT: JWTConfig{
			Secret:           getEnv("JWT_SECRET", "your-secret-key-change-this"),
			SuperAdminSecret: getEnv("SUPER_ADMIN_JWT_SECRET", ""),
			Expiry:           getEnv("JWT_EXPIRY", "24h"),
		},
		Email: emailConfig,
	}

	if cfg.Server.Env == "production" {
		if cfg.JWT.Secret == "your-secret-key-change-this" || len(cfg.JWT.Secret) < 32 {
			return nil, fmt.Errorf("JWT_SECRET must be set to a secure value (>= 32 chars) in production")
		}
	}

	return cfg, nil
}

func loadEmailConfig() (EmailConfig, error) {
	port := 587
	if rawPort := strings.TrimSpace(os.Getenv("SMTP_PORT")); rawPort != "" {
		parsedPort, err := strconv.Atoi(rawPort)
		if err != nil || parsedPort < 1 || parsedPort > 65535 {
			return EmailConfig{}, fmt.Errorf("SMTP_PORT must be between 1 and 65535")
		}
		port = parsedPort
	}
	config := EmailConfig{
		SMTPHost: strings.TrimSpace(os.Getenv("SMTP_HOST")), SMTPPort: port,
		SMTPUsername: strings.TrimSpace(os.Getenv("SMTP_USERNAME")), SMTPPassword: os.Getenv("SMTP_PASSWORD"),
		FromEmail: strings.TrimSpace(os.Getenv("SMTP_FROM_EMAIL")), FromName: strings.TrimSpace(os.Getenv("SMTP_FROM_NAME")),
	}
	if !config.SMTPConfigured() {
		return config, nil
	}
	if config.SMTPHost == "" || config.SMTPUsername == "" || config.SMTPPassword == "" || config.FromEmail == "" {
		return EmailConfig{}, fmt.Errorf("SMTP_HOST, SMTP_USERNAME, SMTP_PASSWORD, and SMTP_FROM_EMAIL must be configured together")
	}
	return config, nil
}

func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.User, c.Password, c.DBName, c.SSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return parsed
}
