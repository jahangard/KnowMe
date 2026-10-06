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
		"log_file", cfg.LogFile,
	)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := database.OpenSQLServer(ctx, cfg.SQLServerDSN)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("database connection established", "provider", "sqlserver")

	bot, err := telegramplatform.NewBot(cfg.TelegramBotToken, cfg.TelegramPollTimeout, db, logger)
	if err != nil {
		logger.Error("telegram bot initialization failed", "error", err)
		os.Exit(1)
	}

	logger.Info("KnowMe started", "environment", cfg.Environment)
	if err := bot.Run(ctx); err != nil {
		logger.Error("KnowMe stopped with error", "error", err)
		os.Exit(1)
	}
	logger.Info("KnowMe stopped cleanly")
}
