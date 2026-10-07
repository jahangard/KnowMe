package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

func MigrateAndSeed(ctx context.Context, handle *Handle) error {
	db := handle.Gorm.WithContext(ctx)

	if err := db.AutoMigrate(
		&User{},
		&UserProfile{},
		&TestCategory{},
		&Test{},
		&Question{},
		&QuestionOption{},
		&TestResultProfile{},
		&TestSession{},
		&TestAnswer{},
		&TestResult{},
		&UserTrait{},
		&UserEvent{},
		&LifeLesson{},
	); err != nil {
		return fmt.Errorf("auto migrate database: %w", err)
	}

	if err := seedCategories(ctx, handle.Gorm); err != nil {
		return fmt.Errorf("seed test categories: %w", err)
	}
	if err := seedTopicTree(ctx, handle.Gorm); err != nil {
		return fmt.Errorf("seed topic tree: %w", err)
	}

	if err := seedInitialTests(ctx, handle.Gorm); err != nil {
		return fmt.Errorf("seed initial tests: %w", err)
	}
	if err := seedPlannedTests(ctx, handle.Gorm); err != nil {
		return fmt.Errorf("seed planned tests: %w", err)
	}
	if err := stripSeededTestTitleEmoji(ctx, handle.Gorm); err != nil {
		return fmt.Errorf("clean test title emoji: %w", err)
	}
	if err := seedLifeLessons(ctx, handle.Gorm); err != nil {
		return fmt.Errorf("seed life lessons: %w", err)
	}

	return nil
}

func stripSeededTestTitleEmoji(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var tests []Test
		if err := tx.Find(&tests).Error; err != nil {
			return err
		}
		for _, test := range tests {
			cleanTitle := stripEmoji(test.Title)
			if cleanTitle != test.Title {
				if err := tx.Model(&Test{}).Where("Id = ?", test.ID).Update("Title", cleanTitle).Error; err != nil {
					return err
				}
			}
		}

		var profiles []TestResultProfile
		if err := tx.Find(&profiles).Error; err != nil {
			return err
		}
		for _, profile := range profiles {
			cleanTitle := stripEmoji(profile.Title)
			if cleanTitle != profile.Title {
				if err := tx.Model(&TestResultProfile{}).Where("Id = ?", profile.ID).Update("Title", cleanTitle).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func seedLifeLessons(ctx context.Context, db *gorm.DB) error {
	lessons := []LifeLesson{
		{
			Code: "living-with-grief", Title: "چطور با غم‌های بزرگ کنار بیاییم؟",
			Summary: "برای روزهایی که فقدان یا اندوه، همه‌چیز را سنگین می‌کند.", SortOrder: 10, IsActive: true,
			Goal:         "هدف، حذف فوری غم یا فراموش کردن اتفاق نیست؛ کمک به توست تا با حفظ پیوندت با زندگی، قدم‌های کوچک و امن برداری.",
			Explanation:  "سوگ می‌تواند بعد از مرگ عزیز، جدایی، بیماری، مهاجرت یا از دست دادن یک فرصت مهم به وجود بیاید. شکل آن برای هرکس فرق دارد: گریه، خشم، بی‌حسی، خستگی یا حتی احساس آسودگی. هیچ زمان‌بندی ثابتی برای خوب شدن وجود ندارد.",
			Consequences: "نادیده گرفتن طولانی‌مدت احساس‌ها یا تنها ماندن ممکن است فشار روانی را بیشتر کند و روی خواب، تمرکز، رابطه‌ها و کارهای روزمره اثر بگذارد. این پیامدها نشانه ضعف نیستند؛ علامتی‌اند که شاید به حمایت بیشتری نیاز داری.",
			AgeRange:     "برای نوجوانان و بزرگسالان؛ نوجوانان بهتر است با همراهی یک بزرگسال امن این مسیر را طی کنند. برای کودکان، توضیح متناسب با سن و حمایت مراقب لازم است.",
			Method:       "۱. احساست را نام ببر؛ لازم نیست آن را قضاوت یا پنهان کنی.\n۲. امروز فقط یک نیاز پایه را انجام بده: آب، غذای ساده، استراحت یا چند دقیقه هوای آزاد.\n۳. از یک آدم امن درخواست مشخصی بکن؛ مثلاً «می‌شود کمی کنارم بمانی؟»\n۴. راه شخصی خودت برای یادآوری یا خداحافظی را پیدا کن: نوشتن، گفت‌وگو، دعا یا یک رسم کوچک.\n۵. تصمیم‌های بزرگ را تا جای ممکن به زمانی موکول کن که فشار اولیه کمتر شده باشد.",
			SeekHelp:     "اگر برای مدتی طولانی انجام کارهای روزمره دشوار مانده، یا احساس می‌کنی از پس فشار برنمی‌آیی، با روان‌شناس، مشاور یا پزشک صحبت کن. اگر فکر آسیب زدن به خودت داری، تنها نمان و همین حالا با فردی مورد اعتماد و خدمات اورژانسی محل زندگی‌ات تماس بگیر.",
			KeyPoint:     "لازم نیست امروز خوب شوی؛ فقط قدم بعدیِ امن را بردار.",
			Content:      "غم با یک جمله ناپدید نمی‌شود. به خودت فرصت بده و از حمایت امن کمک بگیر.",
		},
		{
			Code: "letting-go", Title: "چطور کسی را فراموش کنیم؟",
			Summary: "راهی واقع‌بینانه برای عبور از دلتنگی و پایان یک رابطه.", SortOrder: 20, IsActive: true,
			Goal:         "هدف این نیست که خاطره را پاک کنی؛ هدف این است که خاطره کم‌کم اختیار امروز و انتخاب‌هایت را کمتر بگیرد.",
			Explanation:  "بعد از پایان رابطه، دلتنگی و دوگانگی طبیعی است. ممکن است ذهنت بیشتر لحظه‌های خوب را به یاد بیاورد یا بخواهد برای برگشتن رابطه نشانه پیدا کند. عبور از این دوره برای هرکس سرعت متفاوتی دارد.",
			Consequences: "چسبیدن مداوم به تماس، پیگیری صفحه‌ها یا امیدهای مبهم ممکن است پذیرش پایان را دشوارتر کند و خواب، تمرکز، عزت‌نفس یا رابطه‌های دیگر را تحت فشار بگذارد. این به معنی سرزنش تو نیست؛ فقط می‌تواند نشانه نیاز به مرز و حمایت باشد.",
			AgeRange:     "برای نوجوانان و بزرگسالان. نوجوانان در رابطه‌های آسیب‌زا یا هنگام احساس ناامنی بهتر است از یک بزرگسال قابل اعتماد کمک بگیرند.",
			Method:       "۱. واقعیت پایان را با خودت روشن و مهربانانه مرور کن.\n۲. برای مدتی محرک‌ها را کم کن: دیدن صفحه‌ها و پیام‌ها را محدود کن و اگر لازم است عکس‌ها را از جلوی چشم بردار.\n۳. وقتی میل ناگهانی به پیام دادن می‌آید، کمی مکث کن؛ حرفت را بنویس و با یک دوست امن در میان بگذار.\n۴. رابطه را کامل به یاد بیاور؛ خوبی‌ها و دشواری‌ها را کنار هم ببین.\n۵. برنامه‌های کوچک خودت را دوباره بساز: خواب، حرکت، دوستی و کاری که برایت معنا دارد.",
			SeekHelp:     "اگر این جدایی برای مدت طولانی خواب، تحصیل، کار یا احساس امنیتت را مختل کرده، یا در رابطه خشونت و تهدید وجود داشته، از مشاور یا فردی قابل اعتماد کمک بگیر. در خطر فوری با خدمات اضطراری محل زندگی‌ات تماس بگیر.",
			KeyPoint:     "قرار نیست یک‌شبه فراموش کنی؛ می‌توانی امروز یک مرز کوچک و مهربانانه برای خودت بسازی.",
			Content:      "فراموش کردن هدف لازم نیست؛ هدف این است که خاطره کم‌کم اختیار امروزت را کمتر بگیرد.",
		},
	}
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, lesson := range lessons {
			topicCode := map[string]string{"living-with-grief": "lessons-grief", "letting-go": "lessons-breakups"}[lesson.Code]
			var topic TestCategory
			if err := tx.Where("Code = ? AND CatalogType = ?", topicCode, "lessons").First(&topic).Error; err != nil {
				return err
			}
			lesson.TopicID = topic.ID
			var existing LifeLesson
			if err := tx.Where("Code = ?", lesson.Code).Assign(lesson).FirstOrCreate(&existing).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
func seedCategories(ctx context.Context, db *gorm.DB) error {
	categories := []TestCategory{
		{Code: "love", Title: "عشق و رابطه", CatalogType: "tests", SortOrder: 10, IsActive: true},
		{Code: "challenge", Title: "چالشی و باحال", CatalogType: "tests", SortOrder: 20, IsActive: true},
		{Code: "self_knowledge", Title: "خودشناسی عمیق‌تر", CatalogType: "tests", SortOrder: 30, IsActive: true},
		{Code: "adult", Title: "۱۸+", CatalogType: "tests", SortOrder: 40, IsActive: true},
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, category := range categories {
			var row TestCategory
			if err := tx.Where("Code = ?", category.Code).
				Attrs(category).
				FirstOrCreate(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func seedTopicTree(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		children := []struct {
			code, title, parent string
			sort                int
		}{
			{"love-style", "شناخت رابطه و جذابیت", "love", 10},
			{"love-dating", "آشنایی و قرار", "love", 20},
			{"love-boundaries", "مرزهای شخصی", "love", 30},
			{"challenge-truth", "روبه‌رو شدن با خود", "challenge", 10},
			{"self-emotions", "احساسات و تصمیم‌ها", "self_knowledge", 10},
			{"personality-types", "تیپ‌های شخصیتی", "self_knowledge", 20},
			{"personality-assessments", "آزمون‌های شخصیت", "personality-types", 10},
			{"career-self-knowledge", "خودشناسی شغلی", "self_knowledge", 30},
			{"intelligence-self-knowledge", "خودشناسی هوش", "self_knowledge", 40},
			{"love-marriage-self-knowledge", "عشق و ازدواج", "self_knowledge", 50},
			{"social-self-knowledge", "خودشناسی اجتماعی", "self_knowledge", 60},
			{"adult-relationships", "رابطه‌ی بزرگسالان", "adult", 10},
		}
		for _, child := range children {
			var parent TestCategory
			if err := tx.Where("Code = ?", child.parent).First(&parent).Error; err != nil {
				return err
			}
			row := TestCategory{Code: child.code, Title: child.title, ParentID: &parent.ID, CatalogType: "tests", SortOrder: child.sort, IsActive: true}
			if err := tx.Where("Code = ?", child.code).Assign(row).FirstOrCreate(&TestCategory{}).Error; err != nil {
				return err
			}
		}

		personalityTypes, err := findCategoryByCode(tx, "personality-types")
		if err != nil {
			return err
		}
		for _, child := range []struct {
			code, title string
			sort        int
		}{
			{"personality-archetypes", "کهن‌الگوها", 10},
			{"personality-mbti", "MBTI", 20},
		} {
			row := TestCategory{Code: child.code, Title: child.title, ParentID: &personalityTypes.ID, CatalogType: "tests", SortOrder: child.sort, IsActive: true}
			if err := tx.Where("Code = ?", child.code).Assign(row).FirstOrCreate(&TestCategory{}).Error; err != nil {
				return err
			}
		}

		lessonTopics := []struct {
			code, title string
			parent      *string
			sort        int
		}{
			{code: "lessons-feelings", title: "احساسات و سوگ", sort: 10},
			{code: "lessons-grief", title: "غم و فقدان", parent: stringPtr("lessons-feelings"), sort: 10},
			{code: "lessons-relationships", title: "رابطه و دلتنگی", sort: 20},
			{code: "lessons-breakups", title: "جدایی و دل‌کندن", parent: stringPtr("lessons-relationships"), sort: 10},
		}
		for _, topic := range lessonTopics {
			var parentID *int64
			if topic.parent != nil {
				var parent TestCategory
				if err := tx.Where("Code = ?", *topic.parent).First(&parent).Error; err != nil {
					return err
				}
				parentID = &parent.ID
			}
			row := TestCategory{Code: topic.code, Title: topic.title, ParentID: parentID, CatalogType: "lessons", SortOrder: topic.sort, IsActive: true}
			if err := tx.Where("Code = ?", topic.code).Assign(row).FirstOrCreate(&TestCategory{}).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func findCategoryByCode(tx *gorm.DB, code string) (TestCategory, error) {
	var category TestCategory
	err := tx.Where("Code = ?", code).First(&category).Error
	return category, err
}

func stringPtr(value string) *string { return &value }
