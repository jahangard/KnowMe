package telegram

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jahangard/KnowMe/internal/application/roadmap"
	"github.com/jahangard/KnowMe/internal/application/testengine"
)

type Bot struct {
	api        *tgbotapi.BotAPI
	timeout    time.Duration
	db         *sql.DB
	logger     *slog.Logger
	roadmap    *roadmap.Service
	testEngine *testengine.Service
}

func NewBot(token string, timeout time.Duration, db *sql.DB, logger *slog.Logger) (*Bot, error) {
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("create telegram bot: %w", err)
	}

	return &Bot{
		api:        api,
		timeout:    timeout,
		db:         db,
		logger:     logger,
		roadmap:    roadmap.New(db),
		testEngine: testengine.New(db),
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
	if update.CallbackQuery != nil {
		return b.handleCallback(ctx, update.CallbackQuery)
	}
	if update.Message == nil || update.Message.From == nil {
		return nil
	}

	userID, err := b.ensureUser(ctx, update.Message.From)
	if err != nil {
		return err
	}

	switch update.Message.Command() {
	case "start":
		return b.sendHome(update.Message.Chat.ID)
	case "roadmap":
		return b.sendRoadmap(ctx, userID, update.Message.Chat.ID)
	default:
		return b.sendHome(update.Message.Chat.ID)
	}
}

func (b *Bot) handleCallback(ctx context.Context, callback *tgbotapi.CallbackQuery) error {
	if callback.From == nil || callback.Message == nil {
		return nil
	}

	userID, err := b.ensureUser(ctx, callback.From)
	if err != nil {
		return err
	}

	if _, err := b.api.Request(tgbotapi.NewCallback(callback.ID, "")); err != nil {
		b.logger.Warn("callback answer failed", "error", err)
	}

	chatID := callback.Message.Chat.ID
	messageID := callback.Message.MessageID

	switch callback.Data {
	case "menu:home":
		return b.sendHome(chatID)
	case "menu:roadmap":
		return b.sendRoadmap(ctx, userID, chatID)
	case "menu:profile":
		return b.sendHTML(chatID, profileText(), homeKeyboard())
	case "menu:about":
		return b.sendHTML(chatID, aboutText(), homeKeyboard())
	}

	if strings.HasPrefix(callback.Data, "test:start:") {
		rawID := strings.TrimPrefix(callback.Data, "test:start:")
		testID, err := strconv.ParseInt(rawID, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid test id: %w", err)
		}

		view, err := b.testEngine.Start(ctx, userID, testID)
		if err != nil {
			return err
		}
		return b.editHTML(chatID, messageID, questionText(view), questionKeyboard(view))
	}

	if strings.HasPrefix(callback.Data, "answer:") {
		sessionID, questionID, optionID, err := parseAnswerCallback(callback.Data)
		if err != nil {
			return err
		}

		next, result, err := b.testEngine.Answer(ctx, userID, sessionID, questionID, optionID)
		if err != nil {
			return err
		}

		if result != nil {
			return b.editHTML(chatID, messageID, resultText(result), resultKeyboard())
		}
		if next != nil {
			return b.editHTML(chatID, messageID, questionText(next), questionKeyboard(next))
		}
	}

	return nil
}

func parseAnswerCallback(data string) (sessionID, questionID, optionID int64, err error) {
	parts := strings.Split(data, ":")
	if len(parts) != 4 || parts[0] != "answer" {
		return 0, 0, 0, fmt.Errorf("invalid answer callback")
	}

	values := []*int64{&sessionID, &questionID, &optionID}
	for i := 1; i < len(parts); i++ {
		value, parseErr := strconv.ParseInt(parts[i], 10, 64)
		if parseErr != nil {
			return 0, 0, 0, fmt.Errorf("invalid answer callback value: %w", parseErr)
		}
		*values[i-1] = value
	}

	return sessionID, questionID, optionID, nil
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

func (b *Bot) sendHome(chatID int64) error {
	return b.sendHTML(chatID, homeText(), homeKeyboard())
}

func (b *Bot) sendRoadmap(ctx context.Context, userID, chatID int64) error {
	rec, err := b.roadmap.NextTest(ctx, userID)
	if err != nil {
		return err
	}
	if rec == nil {
		return b.sendHTML(
			chatID,
			"<b>🎉 مسیر فعلی کامل شد</b>\n\nفعلاً تست جدیدی برای پیشنهاد ندارم. بعداً مسیرهای بیشتری اضافه می‌کنیم.",
			homeKeyboard(),
		)
	}

	return b.sendHTML(chatID, roadmapText(rec.Title), roadmapKeyboard(rec.TestID))
}

func (b *Bot) sendHTML(chatID int64, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = parseModeHTML
	msg.DisableWebPagePreview = true
	msg.ReplyMarkup = keyboard
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) editHTML(chatID int64, messageID int, text string, keyboard tgbotapi.InlineKeyboardMarkup) error {
	edit := tgbotapi.NewEditMessageTextAndMarkup(chatID, messageID, text, keyboard)
	edit.ParseMode = parseModeHTML
	edit.DisableWebPagePreview = true
	_, err := b.api.Send(edit)
	return err
}
