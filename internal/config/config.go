package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Environment         string
	TelegramBotToken    string
	TelegramPollTimeout time.Duration
	SQLServerDSN        string
	LogLevel            slog.Level
}

func Load() (Config, error) {
	cfg := Config{
		Environment:         valueOrDefault("APP_ENV", "development"),
		TelegramBotToken:    strings.TrimSpace(os.Getenv("TELEGRAM_BOT_TOKEN")),
		SQLServerDSN:        strings.TrimSpace(os.Getenv("SQLSERVER_DSN")),
		TelegramPollTimeout: 30 * time.Second,
		LogLevel:            parseLogLevel(valueOrDefault("LOG_LEVEL", "info")),
	}

	if raw := strings.TrimSpace(os.Getenv("TELEGRAM_POLL_TIMEOUT_SECONDS")); raw != "" {
		seconds, err := strconv.Atoi(raw)
		if err != nil || seconds < 1 {
			return Config{}, fmt.Errorf("TELEGRAM_POLL_TIMEOUT_SECONDS must be a positive integer")
		}
		cfg.TelegramPollTimeout = time.Duration(seconds) * time.Second
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
