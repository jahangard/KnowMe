package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/jahangard/KnowMe/internal/config"
	"github.com/jahangard/KnowMe/internal/platform/database"
	"github.com/jahangard/KnowMe/internal/platform/logging"
	telegramplatform "github.com/jahangard/KnowMe/internal/platform/telegram"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	logger, closeLogs, err := logging.New(logging.Config{
		Level:      cfg.LogLevel,
		File:       cfg.LogFile,
		MaxSizeMB:  cfg.LogMaxSizeMB,
		MaxBackups: cfg.LogMaxBackups,
		MaxAgeDays: cfg.LogMaxAgeDays,
		Compress:   cfg.LogCompress,
	})
	if err != nil {
		slog.Error("logging initialization failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		if err := closeLogs(); err != nil {
			fmt.Fprintln(os.Stderr, "log close error:", err)
		}
	}()

	slog.SetDefault(logger)

	defer func() {
		if recovered := recover(); recovered != nil {
			logger.Error(
				"panic recovered at process boundary",
				"panic", recovered,
				"stack", string(debug.Stack()),
			)
			panic(recovered)
		}
	}()

	logger.Info(
		"KnowMe process starting",
		"environment", cfg.Environment,
		"db_provider", cfg.DBProvider,
		"log_file", cfg.LogFile,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	provider, err := database.ParseProvider(cfg.DBProvider)
	if err != nil {
		logger.Error("database provider configuration failed", "error", err)
		os.Exit(1)
	}

	db, err := database.Open(ctx, provider, cfg.SQLServerDSN, cfg.SQLitePath)
	if err != nil {
		logger.Error("database connection failed", "provider", provider, "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection established", "provider", provider)

	if err := database.MigrateAndSeed(ctx, db); err != nil {
		logger.Error("database migration failed", "provider", provider, "error", err)
		os.Exit(1)
	}
	logger.Info("database migration complete", "provider", provider)

	bot, err := telegramplatform.NewBot(cfg.TelegramBotToken, cfg.TelegramPollTimeout, db.Gorm, logger)
	if err != nil {
		logger.Error("telegram bot initialization failed", "error", err)
		os.Exit(1)
	}

	logger.Info("KnowMe started", "environment", cfg.Environment, "db_provider", provider)
	if err := bot.Run(ctx); err != nil {
		logger.Error("KnowMe stopped with error", "error", err)
		os.Exit(1)
	}
	logger.Info("KnowMe stopped cleanly")
}
