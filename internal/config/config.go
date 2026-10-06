package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment         string
	TelegramBotToken    string
	TelegramPollTimeout time.Duration
	SQLServerDSN        string
	LogLevel            slog.Level
	LogFile             string
	LogMaxSizeMB        int
	LogMaxBackups       int
	LogMaxAgeDays       int
	LogCompress         bool
}

func Load() (Config, error) {
	// Load .env for local development. Existing OS environment variables
	// keep priority because godotenv.Load does not overwrite them.
	_ = godotenv.Load()

	cfg := Config{
		Environment:         valueOrDefault("APP_ENV", "development"),
		TelegramBotToken:    strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")),
		SQLServerDSN:        strings.TrimSpace(os.Getenv("SQLSERVER_DSN")),
		TelegramPollTimeout: 30 * time.Second,
		LogLevel:            parseLogLevel(valueOrDefault("LOG_LEVEL", "info")),
		LogFile:             valueOrDefault("LOG_FILE", "logs/knowme.log"),
		LogMaxSizeMB:        20,
		LogMaxBackups:       10,
		LogMaxAgeDays:       14,
		LogCompress:         true,
	}

	var err error
	if cfg.TelegramPollTimeout, err = durationSecondsFromEnv("TELEGRAM_POLL_TIMEOUT_SECONDS", 30); err != nil {
		return Config{}, err
	}
	if cfg.LogMaxSizeMB, err = positiveIntFromEnv("LOG_MAX_SIZE_MB", 20); err != nil {
		return Config{}, err
	}
	if cfg.LogMaxBackups, err = positiveIntFromEnv("LOG_MAX_BACKUPS", 10); err != nil {
		return Config{}, err
	}
	if cfg.LogMaxAgeDays, err = positiveIntFromEnv("LOG_MAX_AGE_DAYS", 14); err != nil {
		return Config{}, err
	}
	if cfg.LogCompress, err = boolFromEnv("LOG_COMPRESS", true); err != nil {
		return Config{}, err
	}

	if cfg.TelegramBotToken == "" {
		return Config{}, fmt.Errorf("TELEGRAM_BOT_TOKEN is required")
	}
	if cfg.SQLServerDSN == "" {
		return Config{}, fmt.Errorf("SQLSERVER_DSN is required")
	}

	return cfg, nil
}

func valueOrDefault(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func durationSecondsFromEnv(key string, fallback int) (time.Duration, error) {
	seconds, err := positiveIntFromEnv(key, fallback)
	if err != nil {
		return 0, err
	}
	return time.Duration(seconds) * time.Second, nil
}

func positiveIntFromEnv(key string, fallback int) (int, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.Atoi(raw)
	if err != nil || value < 1 {
		return 0, fmt.Errorf("%s must be a positive integer", key)
	}
	return value, nil
}

func boolFromEnv(key string, fallback bool) (bool, error) {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback, nil
	}
	value, err := strconv.ParseBool(raw)
	if err != nil {
		return false, fmt.Errorf("%s must be true or false", key)
	}
	return value, nil
}

func parseLogLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
