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
	"github.com/jahangard/KnowMe/internal/application/profileservice"
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
	profile    *profileservice.Service
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
		profile:    profileservice.New(db),
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

	profile, err := b.profile.Get(ctx, userID)
	if err != nil {
		return err
	}

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			return b.startForProfile(ctx, userID, update.Message.Chat.ID, profile)
		case "roadmap":
			return b.sendRoadmapGuarded(ctx, userID, update.Message.Chat.ID)
		default:
			return b.startForProfile(ctx, userID, update.Message.Chat.ID, profile)
		}
	}

	if profile.PendingField != nil {
		switch *profile.PendingField {
		case "age":
			return b.handleAgeMessage(ctx, userID, update.Message.Chat.ID, update.Message.Text)
		case "name":
			return b.handleNameMessage(ctx, userID, update.Message.Chat.ID, update.Message.Text)
		}
	}

	return b.startForProfile(ctx, userID, update.Message.Chat.ID, profile)
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
		return b.sendRoadmapGuarded(ctx, userID, chatID)
	case "menu:profile":
		return b.sendProfile(ctx, userID, chatID)
	case "menu:about":
		return b.sendHTML(chatID, aboutText(), homeKeyboard())
	case "profile:name:start":
		if err := b.profile.BeginNameCapture(ctx, userID); err != nil {
			return err
		}
		return b.editHTML(chatID, messageID, namePromptText(), namePromptKeyboard())
	case "profile:name:skip":
		if err := b.profile.ClearPending(ctx, userID); err != nil {
			return err
		}
		return b.editHTML(chatID, messageID, homeText(), homeKeyboard())
	}

	if strings.HasPrefix(callback.Data, "profile:gender:") {
		gender := strings.TrimPrefix(callback.Data, "profile:gender:")
		switch gender {
		case "male", "female", "prefer_not_say":
		default:
			return fmt.Errorf("invalid gender value")
		}
		if err := b.profile.SetGender(ctx, userID, gender); err != nil {
			return err
		}
		return b.editHTML(chatID, messageID, agePromptText(), emptyInlineKeyboard())
	}

	if strings.HasPrefix(callback.Data, "test:start:") {
		if ok, err := b.baseProfileComplete(ctx, userID, chatID); err != nil || !ok {
			return err
		}

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
			if err := b.editHTML(chatID, messageID, resultText(result), resultKeyboard()); err != nil {
				return err
			}
			return b.maybePromptForName(ctx, userID, chatID)
		}
		if next != nil {
			return b.editHTML(chatID, messageID, questionText(next), questionKeyboard(next))
		}
	}

	return nil
}

func (b *Bot) startForProfile(ctx context.Context, userID, chatID int64, profile *profileservice.Profile) error {
	if profile.Gender == nil {
		return b.sendHTML(chatID, genderPromptText(), genderKeyboard())
	}
	if profile.Age == nil {
		if err := b.profile.BeginAgeCapture(ctx, userID); err != nil {
			return err
		}
		return b.sendHTML(chatID, agePromptText(), emptyInlineKeyboard())
	}
	return b.sendHome(chatID)
}

func (b *Bot) baseProfileComplete(ctx context.Context, userID, chatID int64) (bool, error) {
	profile, err := b.profile.Get(ctx, userID)
	if err != nil {
		return false, err
	}
	if profile.Gender == nil {
		return false, b.sendHTML(chatID, genderPromptText(), genderKeyboard())
	}
	if profile.Age == nil {
		if err := b.profile.BeginAgeCapture(ctx, userID); err != nil {
			return false, err
		}
		return false, b.sendHTML(chatID, agePromptText(), emptyInlineKeyboard())
	}
	return true, nil
}

func (b *Bot) handleAgeMessage(ctx context.Context, userID, chatID int64, raw string) error {
	normalized := normalizeDigits(strings.TrimSpace(raw))
	age, err := strconv.Atoi(normalized)
	if err != nil || age < 13 || age > 100 {
		return b.sendHTML(chatID, invalidAgeText(), emptyInlineKeyboard())
	}

	if err := b.profile.SetAge(ctx, userID, age); err != nil {
		return err
	}

	return b.sendHTML(
		chatID,
		"<b>✅ ثبت شد</b>\n\nحالا می‌تونی مستقیم وارد نقشه راهت بشی.",
		homeKeyboard(),
	)
}

func (b *Bot) handleNameMessage(ctx context.Context, userID, chatID int64, raw string) error {
	name := strings.TrimSpace(raw)
	if name == "" || len([]rune(name)) > 40 {
		return b.sendHTML(
			chatID,
			"<b>اسم یا لقب کوتاه‌تری بفرست</b>\n\nحداکثر ۴۰ کاراکتر.",
			namePromptKeyboard(),
		)
	}

	if err := b.profile.SetName(ctx, userID, name); err != nil {
		return err
	}

	return b.sendHTML(
		chatID,
		"<b>خوشبختم، "+htmlEscape(name)+" 👋</b>\n\nاز این به بعد پروفایلت یک قدم شخصی‌تر شد.",
		homeKeyboard(),
	)
}

func (b *Bot) maybePromptForName(ctx context.Context, userID, chatID int64) error {
	profile, err := b.profile.Get(ctx, userID)
	if err != nil {
		return err
	}
	if profile.Name != nil || profile.NamePrompted {
		return nil
	}

	completed, err := b.profile.CompletedTestCount(ctx, userID)
	if err != nil {
		return err
	}
	if completed < 1 {
		return nil
	}

	if err := b.profile.BeginNameCapture(ctx, userID); err != nil {
		return err
	}
	return b.sendHTML(chatID, namePromptText(), namePromptKeyboard())
}

func (b *Bot) sendProfile(ctx context.Context, userID, chatID int64) error {
	profile, err := b.profile.Get(ctx, userID)
	if err != nil {
		return err
	}
	completed, err := b.profile.CompletedTestCount(ctx, userID)
	if err != nil {
		return err
	}
	return b.sendHTML(chatID, profileText(profile, completed), profileKeyboard(profile))
}

func (b *Bot) sendRoadmapGuarded(ctx context.Context, userID, chatID int64) error {
	ok, err := b.baseProfileComplete(ctx, userID, chatID)
	if err != nil || !ok {
		return err
	}
	return b.sendRoadmap(ctx, userID, chatID)
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

func normalizeDigits(value string) string {
	replacer := strings.NewReplacer(
		"۰", "0", "۱", "1", "۲", "2", "۳", "3", "۴", "4",
		"۵", "5", "۶", "6", "۷", "7", "۸", "8", "۹", "9",
		"٠", "0", "١", "1", "٢", "2", "٣", "3", "٤", "4",
		"٥", "5", "٦", "6", "٧", "7", "٨", "8", "٩", "9",
	)
	return replacer.Replace(value)
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
