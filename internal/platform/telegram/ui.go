package telegram

import (
	"fmt"
	"math"
	"sort"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jahangard/KnowMe/internal/application/profileservice"
	"github.com/jahangard/KnowMe/internal/application/testcatalog"
	"github.com/jahangard/KnowMe/internal/application/testengine"
	store "github.com/jahangard/KnowMe/internal/platform/database"
)

const parseModeHTML = "HTML"

func homeText() string {
	return "<b>✨ KnowMe</b>\n" +
		"خودت رو بهتر بشناس؛ کوتاه، جذاب و شخصی‌سازی‌شده.\n\n" +
		"می‌تونی از بین موضوعات، تست دلخواهت رو انتخاب کنی یا بذاری نقشه راه، قدم بعدی رو برات انتخاب کنه.\n\n" +
		"<i>از «موضوعات تست‌ها» شروع کن.</i>"
}

func homeKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧩 موضوعات تست‌ها", "menu:topics"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🌱 آموزه‌های زندگی", "menu:life-lessons"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🧭 نقشه راه من", "menu:roadmap"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("👤 پروفایل من", "menu:profile"),
			tgbotapi.NewInlineKeyboardButtonData("✨ درباره KnowMe", "menu:about"),
		),
	)
}

func lifeLessonsText(topics []store.TestCategory) string {
	if len(topics) == 0 {
		return "<b>🌱 آموزه‌های زندگی</b>\n\nفعلاً موضوعی منتشر نشده است."
	}
	return "<b>🌱 آموزه‌های زندگی</b>\n\nموضوع موردنظرت را انتخاب کن:"
}

func lifeLessonsKeyboard(topics []store.TestCategory) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(topics)+1)
	for _, topic := range topics {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(topic.Title, fmt.Sprintf("life-topic:%d", topic.ID)),
		))
	}
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func lifeTopicText(title string) string {
	return "<b>🌱 " + htmlEscape(title) + "</b>\n\nزیرموضوع را انتخاب کن:"
}

func lifeTopicsKeyboard(topics []testcatalog.Category) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(topics)+2)
	for _, topic := range topics {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(topic.Title, fmt.Sprintf("life-topic:%d", topic.ID))))
	}
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ موضوع‌های آموزه‌ها", "menu:life-lessons")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home")),
	)
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func lifeLessonItemsText(lessons []store.LifeLesson) string {
	if len(lessons) == 0 {
		return "<b>🌱 آموزه‌های این موضوع</b>\n\nفعلاً آموزه‌ای در این موضوع منتشر نشده است."
	}
	return "<b>🌱 آموزه‌های این موضوع</b>\n\nبرای خواندن راهنما، یکی را انتخاب کن:"
}

func lifeLessonItemsKeyboard(lessons []store.LifeLesson) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(lessons)+2)
	for _, lesson := range lessons {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(lesson.Title, fmt.Sprintf("life-lesson:%d", lesson.ID))))
	}
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ موضوع‌های آموزه‌ها", "menu:life-lessons")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home")),
	)
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func topicText(title string) string {
	return "<b>🧩 " + htmlEscape(title) + "</b>\n\nزیرموضوع را انتخاب کن:"
}

func topicKeyboard(topics []testcatalog.Category) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(topics)+2)
	for _, topic := range topics {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(topic.Title, fmt.Sprintf("topic:category:%d", topic.ID))))
	}
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ موضوعات", "menu:topics")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home")),
	)
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func lifeLessonText(lesson *store.LifeLesson) string {
	return "🌱 <b>" + htmlEscape(lesson.Title) + "</b>\n" +
		"━━━━━━━━━━━━━━\n" +
		"<i>" + htmlEscape(lesson.Summary) + "</i>\n\n" +
		"🎯 <b>هدف این آموزه</b>\n" + htmlEscape(lesson.Goal) + "\n\n" +
		"🪞 <b>موضوع چیست؟</b>\n" + htmlEscape(lesson.Explanation) + "\n\n" +
		"⚠️ <b>اگر نادیده‌اش بگیریم</b>\n" + htmlEscape(lesson.Consequences) + "\n\n" +
		"👥 <b>مناسب چه سنی است؟</b>\n" + htmlEscape(lesson.AgeRange) + "\n\n" +
		"🧭 <b>چه کارهایی می‌توانم انجام بدهم؟</b>\n" + htmlEscape(lesson.Method) + "\n\n" +
		"🤝 <b>چه زمانی کمک بگیرم؟</b>\n" + htmlEscape(lesson.SeekHelp) + "\n\n" +
		"💚 <b>یادت بماند</b>\n" + htmlEscape(lesson.KeyPoint)
}

func lifeLessonKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ همه آموزه‌ها", "menu:life-lessons")),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home")),
	)
}

func topicsText() string {
	return "<b>🧩 موضوعات تست‌ها</b>\n\n" +
		"موضوعی که بیشتر کنجکاوت می‌کنه انتخاب کن؛ یا اگر نمی‌خوای انتخاب کنی، <b>نقشه راه من</b> خودش تست بعدی رو برات پیدا می‌کنه.\n\n" +
		"<i>هرچی بیشتر تست انجام بدی، پیشنهادهای نقشه راه شخصی‌تر می‌شن.</i>"
}

func topicsKeyboard(categories []testcatalog.Category) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(categories)+3)
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🧭 نقشه راه من", "menu:roadmap"),
	))

	for i := 0; i < len(categories); i += 2 {
		row := []tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData(
				categoryLabel(categories[i]),
				fmt.Sprintf("topic:category:%d", categories[i].ID),
			),
		}
		if i+1 < len(categories) {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(
				categoryLabel(categories[i+1]),
				fmt.Sprintf("topic:category:%d", categories[i+1].ID),
			))
		}
		rows = append(rows, row)
	}

	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home"),
	))
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func categoryLabel(category testcatalog.Category) string {
	icon := "🧩"
	switch category.Code {
	case "love":
		icon = "❤️"
	case "challenge":
		icon = "🎯"
	case "self_knowledge":
		icon = "🧠"
	case "adult":
		icon = "🔞"
	}
	return icon + " " + category.Title
}

func categoryText(category *testcatalog.Category, tests []testcatalog.Test) string {
	if len(tests) == 0 {
		return "<b>" + htmlEscape(categoryLabel(*category)) + "</b>\n\n" +
			"این موضوع آماده‌ست، ولی هنوز تست فعالی داخلش منتشر نشده.\n\n" +
			"<i>به‌زودی تست‌های این بخش اضافه می‌شن.</i>"
	}

	return "<b>" + htmlEscape(categoryLabel(*category)) + "</b>\n\n" +
		"یکی از تست‌های این موضوع رو انتخاب کن:"
}

func categoryKeyboard(tests []testcatalog.Test) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, len(tests)+2)
	for _, test := range tests {
		callbackData := fmt.Sprintf("test:start:%d", test.ID)
		label := test.Title
		if !test.IsReady {
			callbackData = fmt.Sprintf("test:info:%d", test.ID)
			label = "📋 " + label
		}
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				label,
				callbackData,
			),
		))
	}
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ موضوعات", "menu:topics"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home"),
		),
	)
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func testInfoText(test *testcatalog.Test) string {
	text := "<b>" + htmlEscape(test.Title) + "</b>\n\n"
	if test.Description != "" {
		text += htmlEscape(test.Description) + "\n\n"
	}
	if test.IsReady {
		return text + "این آزمون آماده‌ی اجراست."
	}
	return text + "<i>سؤال‌ها و نتیجه‌سنجی این آزمون هنوز در بات آماده نشده است.</i>"
}

func testInfoKeyboard(test *testcatalog.Test) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 3)
	if test.IsReady {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("شروع آزمون", fmt.Sprintf("test:start:%d", test.ID)),
		))
	} else {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🕓 به‌زودی", fmt.Sprintf("test:soon:%d", test.ID)),
		))
	}
	rows = append(rows,
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("⬅️ بازگشت به فهرست", fmt.Sprintf("topic:category:%d", test.CategoryID))),
		tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home")),
	)
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func adultGateText() string {
	return "<b>🔞 این بخش مخصوص ۱۸ سال به بالاست</b>\n\n" +
		"بر اساس سنی که در پروفایلت ثبت شده، فعلاً این موضوع برات نمایش داده نمی‌شه."
}

func adultGateKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("⬅️ برگشت به موضوعات", "menu:topics"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🏠 خانه", "menu:home"),
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

func mobilePromptText() string {
	return "<b>📱 یک گزینه اختیاری</b>\n\n" +
		"اگر دوست داری شماره موبایلت هم روی پروفایل ذخیره بشه، می‌تونی با دکمه پایین شماره خودت رو مستقیم از تلگرام بفرستی.\n\n" +
		"<i>این مرحله کاملاً اختیاریه و شماره‌ات داخل پیام پروفایل نمایش داده نمی‌شه.</i>"
}

func mobileReplyKeyboard() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact("📱 ارسال شماره من"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("فعلاً نه"),
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

	mobile := "ثبت نشده"
	if p.Mobile != nil && strings.TrimSpace(*p.Mobile) != "" {
		mobile = "ثبت شده ✅"
	}

	return "<b>👤 پروفایل KnowMe</b>\n\n" +
		"نام: <b>" + name + "</b>\n" +
		"سن: <b>" + age + "</b>\n" +
		"جنسیت: <b>" + gender + "</b>\n" +
		"موبایل: <b>" + mobile + "</b>\n" +
		fmt.Sprintf("تست‌های کامل‌شده: <b>%d</b>\n\n", completedTests) +
		"<i>پروفایل به‌مرور کامل می‌شه؛ لازم نیست همه اطلاعات رو یک‌جا وارد کنی.</i>"
}

func profileKeyboard(p *profileservice.Profile) tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0, 4)
	if p.Name == nil || strings.TrimSpace(*p.Name) == "" {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✍️ ثبت اسم", "profile:name:start"),
		))
	}
	if p.Mobile == nil || strings.TrimSpace(*p.Mobile) == "" {
		rows = append(rows, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("📱 ثبت موبایل", "profile:mobile:start"),
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
		"<b>" + htmlEscape(result.ResultTitle) + "</b>\n" +
		htmlEscape(result.Subtitle) + "\n\n" +
		htmlEscape(result.Description) + "\n\n" +
		"<b>شدت این سبک:</b> " + strings.Repeat("★", stars) + strings.Repeat("☆", 5-stars) + "\n" +
		fmt.Sprintf("<b>سهم از پاسخ‌ها:</b> %d%%\n\n", percent) +
		scoreBreakdown(result.Scores, result.Labels) +
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

func scoreBreakdown(scores map[string]float64, labels map[string]string) string {
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
		name := labels[item.key]
		if strings.TrimSpace(name) == "" {
			name = item.key
		}
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

func htmlEscape(s string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
	)
	return replacer.Replace(s)
}
