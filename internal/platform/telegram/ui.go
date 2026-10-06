package telegram

import (
	"fmt"
	"math"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jahangard/KnowMe/internal/application/profileservice"
	"github.com/jahangard/KnowMe/internal/application/testengine"
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

func genderPromptText() string {
	return "<b>👋 شروع خیلی کوتاه</b>\n\n" +
		"برای اینکه پیشنهادها کمی بهتر شخصی‌سازی بشن، فقط دو چیز رو اول می‌پرسم.\n\n" +
		"<b>جنسیتت رو انتخاب کن:</b>"
}

func genderKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👨 مرد", "profile:gender:male"),
			tgbotapi.NewInlineKeyboardButtonData("👩 زن", "profile:gender:female"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("ترجیح می‌دم نگم", "profile:gender:prefer_not_say"),
		),
	)
}

func agePromptText() string {
	return "<b>🎂 فقط یک مورد دیگه</b>\n\n" +
		"سنت رو فقط به‌صورت عدد بفرست.\n" +
		"مثلاً: <code>32</code>\n\n" +
		"<i>بعد از این مستقیم وارد تست‌ها می‌شی.</i>"
}

func invalidAgeText() string {
	return "<b>سن رو به‌صورت عدد بفرست</b>\n\nمثلاً: <code>32</code>"
}

func namePromptText() string {
	return "<b>✨ یه قدم کوچیک برای شخصی‌تر شدن</b>\n\n" +
		"حالا که اولین تستت رو انجام دادی، دوست داری با چه اسمی صدات کنم؟\n\n" +
		"<i>فقط اسم یا لقبی که خودت دوست داری کافیه.</i>"
}

func namePromptKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("فعلاً نه", "profile:name:skip"),
		),
	)
}

func emptyInlineKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
	}
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

func profileText(p *profileservice.Profile, completedTests int) string {
	name := "ثبت نشده"
	if p.Name != nil && strings.TrimSpace(*p.Name) != "" {
		name = htmlEscape(*p.Name)
	}

	age := "ثبت نشده"
	if p.Age != nil {
		age = fmt.Sprintf("%d سال", *p.Age)
	}

	gender := "ثبت نشده"
	if p.Gender != nil {
		switch *p.Gender {
		case "male":
			gender = "مرد"
		case "female":
			gender = "زن"
		case "prefer_not_say":
			gender = "ترجیح داده نشده"
		default:
			gender = htmlEscape(*p.Gender)
		}
	}

	return "<b>👤 پروفایل KnowMe</b>\n\n" +
		"نام: <b>" + name + "</b>\n" +
		"سن: <b>" + age + "</b>\n" +
		"جنسیت: <b>" + gender + "</b>\n" +
		fmt.Sprintf("تست‌های کامل‌شده: <b>%d</b>\n\n", completedTests) +
		"<i>پروفایل به‌مرور کامل می‌شه؛ لازم نیست همه اطلاعات رو یک‌جا وارد کنی.</i>"
}

func profileKeyboard(p *profileservice.Profile) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
	if p.Name == nil || strings.TrimSpace(*p.Name) == "" {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✍️ ثبت اسم", "profile:name:start"),
		))
	}
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧭 نقشه راه", "menu:roadmap"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home"),
		),
	)
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func aboutText() string {
	return "<b>✨ KnowMe چیه؟</b>\n\n" +
		"یک ربات تست و خودشناسیه که هر بار بیشتر با سبک پاسخ‌هات آشنا می‌شه و مسیر بعدی رو هوشمندتر انتخاب می‌کنه.\n\n" +
		"نتیجه‌ها برای سرگرمی و خودشناسی طراحی می‌شن و تشخیص پزشکی یا روان‌شناختی نیستند."
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

func questionText(view *testengine.QuestionView) string {
	return fmt.Sprintf(
		"<b>%s</b>\n\n%s\n<b>%d از %d</b>\n\n%s",
		htmlEscape(view.TestTitle),
		progressBar(view.Current, view.Total),
		view.Current,
		view.Total,
		htmlEscape(view.Text),
	)
}

func questionKeyboard(view *testengine.QuestionView) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(view.Options)+1)
	for _, option := range view.Options {
		callback := fmt.Sprintf("answer:%d:%d:%d", view.SessionID, view.QuestionID, option.ID)
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(option.Text, callback),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🏠 خروج از تست", "menu:home"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func resultText(result *testengine.Result) string {
	title, subtitle, description := traitPresentation(result.TraitKey)

	total := 0.0
	for _, score := range result.Scores {
		total += score
	}
	percent := 0
	if total > 0 {
		percent = int(math.Round(result.Score / total * 100))
	}

	stars := int(math.Round(float64(percent) / 20.0))
	if stars < 1 {
		stars = 1
	}
	if stars > 5 {
		stars = 5
	}

	return "<b>✨ نتیجه تست تو</b>\n\n" +
		"<b>" + htmlEscape(result.TestTitle) + "</b>\n\n" +
		"<b>" + title + "</b>\n" +
		subtitle + "\n\n" +
		description + "\n\n" +
		"<b>شدت این سبک:</b> " + strings.Repeat("★", stars) + strings.Repeat("☆", 5-stars) + "\n" +
		fmt.Sprintf("<b>سهم از پاسخ‌ها:</b> %d%%\n\n", percent) +
		scoreBreakdown(result.Scores) +
		"\n\n<i>این نتیجه برای سرگرمی و خودشناسی طراحی شده و تشخیص روان‌شناختی نیست.</i>"
}

func resultKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧭 تست بعدی", "menu:roadmap"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👤 پروفایل من", "menu:profile"),
			tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home"),
		),
	)
}

func scoreBreakdown(scores map[string]float64) string {
	type item struct {
		key   string
		score float64
	}
	items := make([]item, 0, len(scores))
	total := 0.0
	for key, score := range scores {
		items = append(items, item{key: key, score: score})
		total += score
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].score == items[j].score {
			return items[i].key < items[j].key
		}
		return items[i].score > items[j].score
	})

	lines := []string{"<b>نقشه سبک‌های تو</b>"}
	for _, item := range items {
		name, _, _ := traitPresentation(item.key)
		percent := 0
		if total > 0 {
			percent = int(math.Round(item.score / total * 100))
		}
		bars := percent / 10
		if bars > 10 {
			bars = 10
		}
		lines = append(lines, fmt.Sprintf(
			"%s %s%s %d%%",
			name,
			strings.Repeat("■", bars),
			strings.Repeat("□", 10-bars),
			percent,
		))
	}
	return strings.Join(lines, "\n")
}

func traitPresentation(key string) (title, subtitle, description string) {
	switch key {
	case "words":
		return "💬 عاشقِ کلمات", "برای تو، حرف خوب فقط حرف نیست.", "ابراز مستقیم احساس، تعریف، تأیید و جمله‌های صمیمی خیلی زود به قلبت راه پیدا می‌کنن. احتمالاً خودت هم وقتی کسی برات مهمه، بیشتر از زبان و کلمات استفاده می‌کنی."
	case "time":
		return "⏳ عاشقِ حضور", "برای تو، وقت گذاشتن یعنی انتخاب کردن.", "حضور واقعی، توجه بدون حواس‌پرتی و وقت دونفره بیشتر از کارهای نمایشی روی تو اثر می‌ذاره. وقتی کسی زمانش رو به تو می‌ده، احساس ارزشمندی بیشتری می‌کنی."
	case "care":
		return "🛠 عاشقِ عمل", "برای تو، دوست داشتن باید دیده بشه.", "کمک کردن، مسئولیت برداشتن و کارهای کوچکِ واقعی برای تو معنی زیادی دارن. احتمالاً بیشتر به رفتار نگاه می‌کنی تا وعده‌ها."
	case "touch":
		return "🤍 عاشقِ نزدیکی", "برای تو، فاصله کم یعنی احساس بیشتر.", "آغوش، تماس و نزدیکی فیزیکیِ محترمانه برای تو یکی از واضح‌ترین نشانه‌های محبت و امنیت عاطفیه."
	default:
		return "✨ سبک ترکیبی", "تو یک الگوی تک‌بعدی نداری.", "چند شیوه مختلف برای دریافت و ابراز علاقه در تو نزدیک به هم هستند."
	}
}

func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)
	return replacer.Replace(s)
}
