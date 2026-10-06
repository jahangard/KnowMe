package telegram

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const parseModeHTML = "HTML"

func homeText() string {
	return "<b>✨ KnowMe</b>\n" +
		"خودت رو بهتر بشناس؛ کوتاه، جذاب و شخصی‌سازی‌شده.\n\n" +
		"من جواب‌هات رو مرحله‌به‌مرحله به خاطر می‌سپارم و بر اساس شناختی که ازت می‌سازم، تست بعدی رو پیشنهاد می‌دم.\n\n" +
		"<i>برای شروع، نقشه راه شخصی‌ات رو باز کن.</i>"
}

func homeKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧭 نقشه راه من", "menu:roadmap"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👤 پروفایل من", "menu:profile"),
			tgbotapi.NewInlineKeyboardButtonData("✨ درباره KnowMe", "menu:about"),
		),
	)
}

func roadmapText(title string) string {
	return "<b>🧭 پیشنهاد بعدی برای تو</b>\n\n" +
		"<b>" + htmlEscape(title) + "</b>\n" +
		"این تست بر اساس مسیر فعلی تو انتخاب شده.\n\n" +
		"⏱ کوتاه و سریع\n" +
		"🎯 نتیجه شخصی‌سازی‌شده\n" +
		"🖼 خروجی قابل اشتراک‌گذاری"
}

func roadmapKeyboard(testID int64) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				"🚀 شروع تست",
				fmt.Sprintf("test:start:%d", testID),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home"),
		),
	)
}

func profileText() string {
	return "<b>👤 پروفایل KnowMe</b>\n\n" +
		"پروفایل تو به‌مرور و با انجام تست‌ها کامل می‌شه؛ بدون فرم ثبت‌نام طولانی.\n\n" +
		"فعلاً از اطلاعات پایه شروع می‌کنیم و بقیه چیزها رو فقط وقتی لازم باشه می‌پرسیم."
}

func aboutText() string {
	return "<b>✨ KnowMe چیه؟</b>\n\n" +
		"یک ربات تست و خودشناسیه که هر بار بیشتر با سبک پاسخ‌هات آشنا می‌شه و مسیر بعدی رو هوشمندتر انتخاب می‌کنه.\n\n" +
		"نتیجه‌ها برای سرگرمی و خودشناسی طراحی می‌شن و تشخیص پزشکی یا روان‌شناختی نیستند."
}

func testSelectedText() string {
	return "<b>🚀 آماده‌ای؟</b>\n\n" +
		"تست انتخاب شد. سؤال‌ها یکی‌یکی نمایش داده می‌شن تا هم سریع باشه، هم خسته‌کننده نشه."
}

func progressBar(current, total int) string {
	if total <= 0 {
		return ""
	}
	if current < 0 {
		current = 0
	}
	if current > total {
		current = total
	}
	return strings.Repeat("●", current) + strings.Repeat("○", total-current)
}

func questionText(title string, current, total int, question string) string {
	return fmt.Sprintf(
		"<b>%s</b>\n\n%s\n<b>%d از %d</b>\n\n%s",
		htmlEscape(title),
		progressBar(current, total),
		current,
		total,
		htmlEscape(question),
	)
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)
	return replacer.Replace(s)
}
