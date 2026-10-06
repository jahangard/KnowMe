package telegram

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jahangard/KnowMe/internal/application/roadmap"
)

type Bot struct {
	api     *tgbotapi.BotAPI
	timeout time.Duration
	db      *sql.DB
	logger  *slog.Logger
	roadmap *roadmap.Service
}

func NewBot(token string, timeout time.Duration, db *sql.DB, logger *slog.Logger) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	return &Bot{
		api:     api,
		timeout: timeout,
		db:      db,
		logger:  logger,
		roadmap: roadmap.New(db),
	}, nil
}

func (b *Bot) Run(ctx context.Context) error {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = int(b.timeout.Seconds())
	updates := b.api.GetUpdatesChan(updateConfig)

	b.logger.Info("telegram long polling started", "bot", b.api.Self.UserName)

	for {
		select {
		case <-ctx.Done():
			b.api.StopReceivingUpdates()
			return nil
		case update, ok := <-updates:
			if !ok {
				return nil
			}
			if err := b.handleUpdate(ctx, update); err != nil {
				b.logger.Error("update handling failed", "error", err, "update_id", update.UpdateID)
			}
		}
	}
}

func (b *Bot) handleUpdate(ctx context.Context, update tgbotapi.Update) error {
	if update.Message == nil || update.Message.From == nil {
		return nil
	}

	userID, err := b.ensureUser(ctx, update.Message.From)
	if err != nil {
		return err
	}

	switch update.Message.Command() {
	case "start":
		return b.sendWelcome(update.Message.Chat.ID)
	case "roadmap":
		return b.sendRoadmap(ctx, userID, update.Message.Chat.ID)
	default:
		return b.sendText(update.Message.Chat.ID, "فعلاً دو دستور آماده است: /start و /roadmap")
	}
}

func (b *Bot) ensureUser(ctx context.Context, tgUser *tgbotapi.User) (int64, error) {
	const query = "SET NOCOUNT ON; " +
		"UPDATE dbo.Users SET Username = @Username, LastSeenAt = SYSUTCDATETIME() WHERE TelegramUserId = @TelegramUserId; " +
		"IF @@ROWCOUNT = 0 BEGIN " +
		"INSERT INTO dbo.Users (TelegramUserId, Username, FirstSeenAt, LastSeenAt) " +
		"VALUES (@TelegramUserId, @Username, SYSUTCDATETIME(), SYSUTCDATETIME()); END; " +
		"SELECT Id FROM dbo.Users WHERE TelegramUserId = @TelegramUserId;"

	var username any
	if tgUser.UserName != "" {
		username = tgUser.UserName
	}

	var id int64
	if err := b.db.QueryRowContext(
		ctx,
		query,
		sql.Named("TelegramUserId", tgUser.ID),
		sql.Named("Username", username),
	).Scan(&id); err != nil {
		return 0, fmt.Errorf("ensure user: %w", err)
	}
	return id, nil
}

func (b *Bot) sendWelcome(chatID int64) error {
	return b.sendText(chatID,
		"به KnowMe خوش اومدی.\n\n"+
			"اینجا تست‌ها رو یکی‌یکی انجام می‌دی و سیستم کم‌کم شناخت بهتری ازت پیدا می‌کنه.\n\n"+
			"برای شروع نقشه راه: /roadmap")
}

func (b *Bot) sendRoadmap(ctx context.Context, userID, chatID int64) error {
	rec, err := b.roadmap.NextTest(ctx, userID)
	if err != nil {
		return err
	}
	if rec == nil {
		return b.sendText(chatID, "فعلاً تست جدیدی برای نقشه راه باقی نمونده.")
	}

	return b.sendText(chatID, fmt.Sprintf(
		"تست پیشنهادی بعدی برای تو:\n%s\n\nTest ID: %d",
		rec.Title,
		rec.TestID,
	))
}

func (b *Bot) sendText(chatID int64, text string) error {
	_, err := b.api.Send(tgbotapi.NewMessage(chatID, text))
	return err
}
