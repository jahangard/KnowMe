package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/jahangard/KnowMe/internal/config"
	"github.com/jahangard/KnowMe/internal/platform/database"
	telegramplatform "github.com/jahangard/KnowMe/internal/platform/telegram"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("configuration error", "error", err)
		os.Exit(1)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	db, err := database.OpenSQLServer(ctx, cfg.SQLServerDSN)
	if err != nil {
		logger.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer db.Close()

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
	logger.Info("KnowMe stopped")
}
