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
		&TestSession{},
		&TestAnswer{},
		&TestResult{},
		&UserTrait{},
		&UserEvent{},
	); err != nil {
		return fmt.Errorf("auto migrate database: %w", err)
	}

	if err := seedCategories(ctx, handle.Gorm); err != nil {
		return fmt.Errorf("seed test categories: %w", err)
	}

	if err := seedLoveStyle(ctx, handle.Gorm); err != nil {
		return fmt.Errorf("seed love style test: %w", err)
	}

	return nil
}

func seedCategories(ctx context.Context, db *gorm.DB) error {
	categories := []TestCategory{
		{Code: "love", Title: "عشق و رابطه", SortOrder: 10, IsActive: true},
		{Code: "challenge", Title: "چالشی و باحال", SortOrder: 20, IsActive: true},
		{Code: "self_knowledge", Title: "خودشناسی عمیق‌تر", SortOrder: 30, IsActive: true},
		{Code: "adult", Title: "۱۸+", SortOrder: 40, IsActive: true},
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, category := range categories {
			var row TestCategory
			if err := tx.Where("Code = ?", category.Code).
				Assign(map[string]any{
					"Title":     category.Title,
					"SortOrder": category.SortOrder,
					"IsActive":  category.IsActive,
				}).
				FirstOrCreate(&row, category).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func seedLoveStyle(ctx context.Context, db *gorm.DB) error {
	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		category := TestCategory{
			Code:      "love",
			Title:     "عشق و رابطه",
			SortOrder: 10,
			IsActive:  true,
		}
		if err := tx.Where("Code = ?", category.Code).FirstOrCreate(&category).Error; err != nil {
			return err
		}

		description := "یک تست کوتاه سرگرمی برای شناخت شیوه غالب ابراز علاقه در رابطه."
		test := Test{
			CategoryID:  category.ID,
			Code:        "love_style_v1",
			Title:       "💞 سبک عشق‌ورزی",
			Description: &description,
			SortOrder:   10,
			IsActive:    true,
		}
		if err := tx.Where("Code = ?", test.Code).FirstOrCreate(&test).Error; err != nil {
			return err
		}

		type optionSeed struct {
			Text     string
			TraitKey string
		}
		type questionSeed struct {
			Text    string
			Options []optionSeed
		}

		questions := []questionSeed{
			{
				Text: "وقتی کسی که دوستش داری ناراحته، معمولاً اولین واکنش تو چیه؟",
				Options: []optionSeed{
					{"با حرف زدن آرومش می‌کنم", "words"},
					{"کنارش می‌مونم و وقتم رو بهش می‌دم", "time"},
					{"یه کاری براش انجام می‌دم که حالش بهتر شه", "care"},
					{"با آغوش و نزدیکی بهش آرامش می‌دم", "touch"},
				},
			},
			{
				Text: "اگر بخوای خیلی واضح نشون بدی که کسی برات مهمه، بیشتر چه کار می‌کنی؟",
				Options: []optionSeed{
					{"بهش می‌گم چقدر برام مهمه", "words"},
					{"یه زمان مخصوص فقط برای دوتامون می‌ذارم", "time"},
					{"کارش رو سبک می‌کنم یا کمکش می‌کنم", "care"},
					{"بیشتر بغلش می‌کنم و نزدیکش می‌مونم", "touch"},
				},
			},
			{
				Text: "در یک روز شلوغ، کدوم رفتار طرف مقابل بیشتر به دلت می‌شینه؟",
				Options: []optionSeed{
					{"یک پیام محبت‌آمیز و صمیمی", "words"},
					{"اینکه با وجود شلوغی، وقت برای من باز کنه", "time"},
					{"اینکه بدون گفتن، یک کارم رو انجام بده", "care"},
					{"یک بغل گرم وقتی همدیگه رو می‌بینیم", "touch"},
				},
			},
			{
				Text: "وقتی دلخور می‌شی، کدوم کار بیشتر کمک می‌کنه دوباره احساس نزدیکی کنی؟",
				Options: []optionSeed{
					{"اینکه حرف دلش رو واضح بگه", "words"},
					{"اینکه بشینیم و باهم وقت بگذرونیم", "time"},
					{"اینکه برای جبران، یک کار واقعی انجام بده", "care"},
					{"اینکه با یک آغوش صمیمی فاصله رو کم کنه", "touch"},
				},
			},
			{
				Text: "برای یک مناسبت خاص، کدوم برنامه بیشتر تو رو خوشحال می‌کنه؟",
				Options: []optionSeed{
					{"یک نامه یا پیام خاص و احساسی", "words"},
					{"یک روز کامل فقط با هم بودن", "time"},
					{"یک کار غافلگیرکننده که زندگی‌م رو راحت‌تر کنه", "care"},
					{"یک شب صمیمی و پر از نزدیکی", "touch"},
				},
			},
			{
				Text: "وقتی از کسی خوشت میاد، خودت ناخودآگاه بیشتر کدوم رفتار رو انجام می‌دی؟",
				Options: []optionSeed{
					{"زیاد تعریف و ابراز احساس می‌کنم", "words"},
					{"دنبال فرصت می‌گردم باهاش تنها باشم", "time"},
					{"کارهای کوچیکش رو انجام می‌دم", "care"},
					{"با تماس و نزدیکی علاقه‌م رو نشون می‌دم", "touch"},
				},
			},
			{
				Text: "در یک رابطه، کدوم کمبود بیشتر اذیتت می‌کنه؟",
				Options: []optionSeed{
					{"کمبود حرف‌های محبت‌آمیز", "words"},
					{"کمبود وقت دونفره", "time"},
					{"اینکه همه چیز فقط در حد حرف بمونه", "care"},
					{"فاصله و سردی فیزیکی", "touch"},
				},
			},
			{
				Text: "اگر فقط یکی رو انتخاب کنی، کدوم جمله بیشتر حس دوست‌داشتن بهت می‌ده؟",
				Options: []optionSeed{
					{"«دوستت دارم و بهت افتخار می‌کنم»", "words"},
					{"«امروز رو کامل برای تو خالی کردم»", "time"},
					{"«نگران نباش، من انجامش دادم»", "care"},
					{"«بیا یه بغل طولانی»", "touch"},
				},
			},
		}

		for qIndex, qSeed := range questions {
			question := Question{
				TestID:       test.ID,
				Text:         qSeed.Text,
				Order:        qIndex + 1,
				QuestionType: "single_choice",
				IsActive:     true,
			}
			if err := tx.Where("TestId = ? AND [Order] = ?", test.ID, question.Order).
				FirstOrCreate(&question).Error; err != nil {
				return err
			}

			for oIndex, oSeed := range qSeed.Options {
				traitKey := oSeed.TraitKey
				option := QuestionOption{
					QuestionID: question.ID,
					Text:       oSeed.Text,
					Order:      oIndex + 1,
					Score:      1,
					TraitKey:   &traitKey,
				}
				if err := tx.Where("QuestionId = ? AND [Order] = ?", question.ID, option.Order).
					FirstOrCreate(&option).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
