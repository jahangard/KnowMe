package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type resultSeed struct {
	TraitKey    string
	Label       string
	Title       string
	Subtitle    string
	Description string
	SortOrder   int
}

type optionSeed struct {
	Text     string
	TraitKey string
	Score    float64
}

type questionSeed struct {
	Text    string
	Options []optionSeed
}

type testSeedDefinition struct {
	CategoryCode string
	Code         string
	Title        string
	Description  string
	SortOrder    int
	Results      []resultSeed
	Questions    []questionSeed
}

func seedInitialTests(ctx context.Context, db *gorm.DB) error {
	definitions := []testSeedDefinition{
		loveStyleSeed(),
		attractionStyleSeed(),
		datingScenariosSeed(),
		personalBoundariesSeed(),
		bitterTruthSeed(),
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, definition := range definitions {
			if err := upsertTestDefinition(tx, definition); err != nil {
				return fmt.Errorf("seed test %s: %w", definition.Code, err)
			}
		}
		return nil
	})
}

func upsertTestDefinition(tx *gorm.DB, definition testSeedDefinition) error {
	var category TestCategory
	if err := tx.Where("Code = ? AND IsActive = ?", definition.CategoryCode, true).
		First(&category).Error; err != nil {
		return fmt.Errorf("load category %s: %w", definition.CategoryCode, err)
	}

	description := definition.Description
	var test Test
	if err := tx.Where("Code = ?", definition.Code).
		Assign(map[string]any{
			"CategoryId":  category.ID,
			"Title":       definition.Title,
			"Description": &description,
			"SortOrder":   definition.SortOrder,
			"IsActive":    true,
		}).
		FirstOrCreate(&test, Test{
			CategoryID:  category.ID,
			Code:        definition.Code,
			Title:       definition.Title,
			Description: &description,
			SortOrder:   definition.SortOrder,
			IsActive:    true,
		}).Error; err != nil {
		return err
	}

	for resultIndex, seed := range definition.Results {
		row := TestResultProfile{
			TestID:      test.ID,
			TraitKey:    seed.TraitKey,
			Label:       seed.Label,
			Title:       seed.Title,
			Subtitle:    seed.Subtitle,
			Description: seed.Description,
			SortOrder:   seed.SortOrder,
			IsActive:    true,
		}
		if row.SortOrder == 0 {
			row.SortOrder = resultIndex + 1
		}
		if err := tx.Where("TestId = ? AND TraitKey = ?", test.ID, seed.TraitKey).
			Assign(map[string]any{
				"Label":       row.Label,
				"Title":       row.Title,
				"Subtitle":    row.Subtitle,
				"Description": row.Description,
				"SortOrder":   row.SortOrder,
				"IsActive":    true,
			}).
			FirstOrCreate(&row).Error; err != nil {
			return err
		}
	}

	for qIndex, qSeed := range definition.Questions {
		question := Question{
			TestID:       test.ID,
			Text:         qSeed.Text,
			Order:        qIndex + 1,
			QuestionType: "single_choice",
			IsActive:     true,
		}
		if err := tx.Where("TestId = ? AND [Order] = ?", test.ID, question.Order).
			Assign(map[string]any{
				"Text":         question.Text,
				"QuestionType": question.QuestionType,
				"IsActive":     true,
			}).
			FirstOrCreate(&question).Error; err != nil {
			return err
		}

		for oIndex, oSeed := range qSeed.Options {
			score := oSeed.Score
			if score == 0 {
				score = 1
			}
			traitKey := oSeed.TraitKey
			option := QuestionOption{
				QuestionID: question.ID,
				Text:       oSeed.Text,
				Order:      oIndex + 1,
				Score:      score,
				TraitKey:   &traitKey,
			}
			if err := tx.Where("QuestionId = ? AND [Order] = ?", question.ID, option.Order).
				Assign(map[string]any{
					"Text":     option.Text,
					"Score":    option.Score,
					"TraitKey": option.TraitKey,
				}).
				FirstOrCreate(&option).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

func loveStyleSeed() testSeedDefinition {
	return testSeedDefinition{
		CategoryCode: "love",
		Code:         "love_style_v1",
		Title:        "💞 سبک عشق‌ورزی",
		Description:  "شیوه غالب ابراز علاقه و دریافت محبتت را پیدا کن.",
		SortOrder:    10,
		Results: []resultSeed{
			{
				TraitKey: "words",
				Label:    "کلمات",
				Title:    "💬 عاشقِ کلمات",
				Subtitle: "برای تو، حرف خوب فقط حرف نیست.",
				Description: "ابراز مستقیم احساس، تعریف و جمله‌های صمیمی خیلی زود به قلبت راه پیدا می‌کنن. " +
					"احتمالاً خودت هم وقتی کسی برات مهمه، از کلمات برای نشان‌دادن علاقه استفاده می‌کنی.",
			},
			{
				TraitKey: "time",
				Label:    "حضور",
				Title:    "⏳ عاشقِ حضور",
				Subtitle: "برای تو، وقت گذاشتن یعنی انتخاب کردن.",
				Description: "حضور واقعی، توجه بدون حواس‌پرتی و وقت دونفره بیشتر از کارهای نمایشی روی تو اثر می‌ذاره. " +
					"وقتی کسی زمانش رو به تو می‌ده، احساس ارزشمندی بیشتری می‌کنی.",
			},
			{
				TraitKey: "care",
				Label:    "عمل و مراقبت",
				Title:    "🛠 عاشقِ عمل",
				Subtitle: "برای تو، دوست داشتن باید دیده بشه.",
				Description: "کمک کردن، مسئولیت برداشتن و کارهای کوچک واقعی برای تو معنی زیادی دارن. " +
					"احتمالاً بیشتر به رفتار نگاه می‌کنی تا وعده‌ها.",
			},
			{
				TraitKey: "touch",
				Label:    "نزدیکی",
				Title:    "🤍 عاشقِ نزدیکی",
				Subtitle: "برای تو، فاصله کم یعنی احساس بیشتر.",
				Description: "آغوش، تماس و نزدیکی فیزیکی محترمانه برای تو یکی از روشن‌ترین نشانه‌های محبت و امنیت عاطفیه.",
			},
		},
		Questions: []questionSeed{
			{
				Text: "وقتی کسی که دوستش داری ناراحته، معمولاً اولین واکنش تو چیه؟",
				Options: []optionSeed{
					{Text: "با حرف زدن آرومش می‌کنم", TraitKey: "words"},
					{Text: "کنارش می‌مونم و وقتم رو بهش می‌دم", TraitKey: "time"},
					{Text: "یه کاری براش انجام می‌دم که حالش بهتر شه", TraitKey: "care"},
					{Text: "با آغوش و نزدیکی بهش آرامش می‌دم", TraitKey: "touch"},
				},
			},
			{
				Text: "اگر بخوای خیلی واضح نشون بدی که کسی برات مهمه، بیشتر چه کار می‌کنی؟",
				Options: []optionSeed{
					{Text: "بهش می‌گم چقدر برام مهمه", TraitKey: "words"},
					{Text: "یه زمان مخصوص فقط برای دوتامون می‌ذارم", TraitKey: "time"},
					{Text: "کارش رو سبک می‌کنم یا کمکش می‌کنم", TraitKey: "care"},
					{Text: "بیشتر بغلش می‌کنم و نزدیکش می‌مونم", TraitKey: "touch"},
				},
			},
			{
				Text: "در یک روز شلوغ، کدوم رفتار طرف مقابل بیشتر به دلت می‌شینه؟",
				Options: []optionSeed{
					{Text: "یک پیام محبت‌آمیز و صمیمی", TraitKey: "words"},
					{Text: "اینکه با وجود شلوغی، وقت برای من باز کنه", TraitKey: "time"},
					{Text: "اینکه بدون گفتن، یک کارم رو انجام بده", TraitKey: "care"},
					{Text: "یک بغل گرم وقتی همدیگه رو می‌بینیم", TraitKey: "touch"},
				},
			},
			{
				Text: "وقتی دلخور می‌شی، کدوم کار بیشتر کمک می‌کنه دوباره احساس نزدیکی کنی؟",
				Options: []optionSeed{
					{Text: "اینکه حرف دلش رو واضح بگه", TraitKey: "words"},
					{Text: "اینکه بشینیم و باهم وقت بگذرونیم", TraitKey: "time"},
					{Text: "اینکه برای جبران، یک کار واقعی انجام بده", TraitKey: "care"},
					{Text: "اینکه با یک آغوش صمیمی فاصله رو کم کنه", TraitKey: "touch"},
				},
			},
			{
				Text: "برای یک مناسبت خاص، کدوم برنامه بیشتر تو رو خوشحال می‌کنه؟",
				Options: []optionSeed{
					{Text: "یک نامه یا پیام خاص و احساسی", TraitKey: "words"},
					{Text: "یک روز کامل فقط با هم بودن", TraitKey: "time"},
					{Text: "یک کار غافلگیرکننده که زندگی‌م رو راحت‌تر کنه", TraitKey: "care"},
					{Text: "یک شب صمیمی و پر از نزدیکی", TraitKey: "touch"},
				},
			},
			{
				Text: "وقتی از کسی خوشت میاد، خودت ناخودآگاه بیشتر کدوم رفتار رو انجام می‌دی؟",
				Options: []optionSeed{
					{Text: "زیاد تعریف و ابراز احساس می‌کنم", TraitKey: "words"},
					{Text: "دنبال فرصت می‌گردم باهاش تنها باشم", TraitKey: "time"},
					{Text: "کارهای کوچیکش رو انجام می‌دم", TraitKey: "care"},
					{Text: "با تماس و نزدیکی علاقه‌م رو نشون می‌دم", TraitKey: "touch"},
				},
			},
			{
				Text: "در یک رابطه، کدوم کمبود بیشتر اذیتت می‌کنه؟",
				Options: []optionSeed{
					{Text: "کمبود حرف‌های محبت‌آمیز", TraitKey: "words"},
					{Text: "کمبود وقت دونفره", TraitKey: "time"},
					{Text: "اینکه همه چیز فقط در حد حرف بمونه", TraitKey: "care"},
					{Text: "فاصله و سردی فیزیکی", TraitKey: "touch"},
				},
			},
			{
				Text: "اگر فقط یکی رو انتخاب کنی، کدوم جمله بیشتر حس دوست‌داشتن بهت می‌ده؟",
				Options: []optionSeed{
					{Text: "«دوستت دارم و بهت افتخار می‌کنم»", TraitKey: "words"},
					{Text: "«امروز رو کامل برای تو خالی کردم»", TraitKey: "time"},
					{Text: "«نگران نباش، من انجامش دادم»", TraitKey: "care"},
					{Text: "«بیا یه بغل طولانی»", TraitKey: "touch"},
				},
			},
		},
	}
}

func attractionStyleSeed() testSeedDefinition {
	return testSeedDefinition{
		CategoryCode: "love",
		Code:         "attraction_style_v1",
		Title:        "✨ تیپ جذابیت",
		Description:  "ببین جذابیت تو بیشتر از چه جنسیه؛ حضور، گرما، رازآلودگی یا انرژی.",
		SortOrder:    20,
		Results: []resultSeed{
			{
				TraitKey: "attraction:magnetic",
				Label:    "مغناطیسی",
				Title:    "🧲 جذابیت مغناطیسی",
				Subtitle: "لازم نیست زیاد تلاش کنی؛ حضورت خودش دیده می‌شه.",
				Description: "اعتمادبه‌نفس، قاطعیت و حضور تو معمولاً باعث می‌شه دیگران ناخودآگاه متوجهت بشن. " +
					"جذابیتت بیشتر از جنس «حضور»ه تا نمایش.",
			},
			{
				TraitKey: "attraction:warm",
				Label:    "گرم",
				Title:    "☀️ جذابیت گرم و صمیمی",
				Subtitle: "کنارت بودن حس راحتی می‌ده.",
				Description: "لبخند، توجه و حس پذیرفته‌شدن بخش اصلی جذابیت توئه. آدم‌ها احتمالاً سریع‌تر از معمول کنار تو خود واقعی‌شون می‌شن.",
			},
			{
				TraitKey: "attraction:mysterious",
				Label:    "رازآلود",
				Title:    "🌙 جذابیت رازآلود",
				Subtitle: "همه چیز را همان اول رو نمی‌کنی.",
				Description: "تو بخشی از خودت را برای کشف شدن نگه می‌داری و همین کنجکاوی ایجاد می‌کنه. " +
					"جذابیتت آرام‌تر شروع می‌شه اما می‌تونه ماندگارتر باشه.",
			},
			{
				TraitKey: "attraction:playful",
				Label:    "بازیگوش",
				Title:    "✨ جذابیت بازیگوش و پرانرژی",
				Subtitle: "با تو فضا زودتر زنده می‌شه.",
				Description: "شوخ‌طبعی، انرژی و واکنش‌های سریع، بخش مهمی از جذابیت تو هستند. " +
					"آدم‌ها معمولاً تو را با حس خوب و لحظه‌های به‌یادماندنی به خاطر می‌سپارن.",
			},
		},
		Questions: []questionSeed{
			{
				Text: "وقتی وارد جمعی می‌شی که بیشتر آدم‌ها رو نمی‌شناسی، معمولاً چطور دیده می‌شی؟",
				Options: []optionSeed{
					{Text: "با اعتمادبه‌نفس وارد می‌شم و حضورم حس می‌شه", TraitKey: "attraction:magnetic"},
					{Text: "زود لبخند می‌زنم و با آدم‌ها گرم می‌گیرم", TraitKey: "attraction:warm"},
					{Text: "اول بیشتر نگاه می‌کنم و کم‌کم وارد می‌شم", TraitKey: "attraction:mysterious"},
					{Text: "با شوخی یا انرژی خوب یخ جمع رو آب می‌کنم", TraitKey: "attraction:playful"},
				},
			},
			{
				Text: "کدوم تعریف بیشتر به دلت می‌شینه؟",
				Options: []optionSeed{
					{Text: "«خیلی کاریزماتیکی»", TraitKey: "attraction:magnetic"},
					{Text: "«کنارت خیلی راحت می‌شه بود»", TraitKey: "attraction:warm"},
					{Text: "«یه چیزی توت هست که آدم می‌خواد بیشتر بشناسدت»", TraitKey: "attraction:mysterious"},
					{Text: "«با تو هیچ‌وقت حوصله‌م سر نمی‌ره»", TraitKey: "attraction:playful"},
				},
			},
			{
				Text: "در اولین برخورد با کسی که برات جذابه، بیشتر چه رفتاری داری؟",
				Options: []optionSeed{
					{Text: "مستقیم و مطمئن ارتباط می‌گیرم", TraitKey: "attraction:magnetic"},
					{Text: "سعی می‌کنم احساس امنیت و راحتی بده", TraitKey: "attraction:warm"},
					{Text: "کمی فاصله نگه می‌دارم تا کنجکاوی ایجاد شه", TraitKey: "attraction:mysterious"},
					{Text: "با شوخی و بازی کلامی فضا رو جذاب می‌کنم", TraitKey: "attraction:playful"},
				},
			},
			{
				Text: "در چت، سبک تو بیشتر شبیه کدومه؟",
				Options: []optionSeed{
					{Text: "کم ولی واضح و مطمئن", TraitKey: "attraction:magnetic"},
					{Text: "صمیمی، مهربان و پیگیر حال طرف", TraitKey: "attraction:warm"},
					{Text: "کم‌حرف‌تر و انتخاب‌شده؛ همه چیز رو نمی‌گم", TraitKey: "attraction:mysterious"},
					{Text: "میم، شوخی، ویس و جواب‌های سریع", TraitKey: "attraction:playful"},
				},
			},
			{
				Text: "اگر توجه همه ناگهان روی تو باشه، چه حسی داری؟",
				Options: []optionSeed{
					{Text: "راحت؛ می‌تونم فضا رو مدیریت کنم", TraitKey: "attraction:magnetic"},
					{Text: "سعی می‌کنم توجه رو با بقیه تقسیم کنم", TraitKey: "attraction:warm"},
					{Text: "ترجیح می‌دم بخشی از من همچنان خصوصی بمونه", TraitKey: "attraction:mysterious"},
					{Text: "ازش استفاده می‌کنم تا فضا رو بامزه‌تر کنم", TraitKey: "attraction:playful"},
				},
			},
			{
				Text: "در استایل شخصی بیشتر دنبال چه اثری هستی؟",
				Options: []optionSeed{
					{Text: "تمیز، قاطع و اثرگذار", TraitKey: "attraction:magnetic"},
					{Text: "دوست‌داشتنی و قابل نزدیک‌شدن", TraitKey: "attraction:warm"},
					{Text: "خاص و کمی غیرقابل‌پیش‌بینی", TraitKey: "attraction:mysterious"},
					{Text: "خلاق، زنده و متفاوت", TraitKey: "attraction:playful"},
				},
			},
			{
				Text: "وقتی کسی بهت علاقه نشون می‌ده، بیشتر کدوم واکنش رو داری؟",
				Options: []optionSeed{
					{Text: "اگر منم علاقه داشته باشم، واضح جلو می‌رم", TraitKey: "attraction:magnetic"},
					{Text: "با توجه و مهربونی جواب می‌دم", TraitKey: "attraction:warm"},
					{Text: "کمی زمان می‌دم تا مطمئن شم", TraitKey: "attraction:mysterious"},
					{Text: "با شیطنت و شوخی واکنش نشون می‌دم", TraitKey: "attraction:playful"},
				},
			},
			{
				Text: "دوست داری آدم‌ها بعد از یک دیدار کوتاه بیشتر چه چیزی ازت یادشون بمونه؟",
				Options: []optionSeed{
					{Text: "اعتمادبه‌نفس و حضورم", TraitKey: "attraction:magnetic"},
					{Text: "حس خوب و آرامشی که دادم", TraitKey: "attraction:warm"},
					{Text: "کنجکاوی درباره اینکه واقعاً چه جور آدمی‌ام", TraitKey: "attraction:mysterious"},
					{Text: "خنده و انرژی‌ای که به فضا دادم", TraitKey: "attraction:playful"},
				},
			},
		},
	}
}

func datingScenariosSeed() testSeedDefinition {
	return testSeedDefinition{
		CategoryCode: "love",
		Code:         "dating_scenarios_v1",
		Title:        "💘 سناریوهای قرار",
		Description:  "در موقعیت‌های واقعی قرار و آشنایی، غریزه تو چطور تصمیم می‌گیره؟",
		SortOrder:    30,
		Results: []resultSeed{
			{
				TraitKey: "dating:planner",
				Label:    "برنامه‌ریز",
				Title:    "🗺 قرارساز حسابگر",
				Subtitle: "تو دوست داری بدانـی کجا داری می‌ری.",
				Description: "برای تو کیفیت قرار با توجه، برنامه و قابل‌اتکا بودن بالا می‌ره. " +
					"غافلگیری بد نیست، ولی دوست داری حداقل چارچوب دستت باشه.",
			},
			{
				TraitKey: "dating:spontaneous",
				Label:    "ماجراجو",
				Title:    "⚡ ماجراجوی لحظه‌ای",
				Subtitle: "بهترین قرار برای تو می‌تونه اصلاً از قبل قرار نباشه.",
				Description: "انرژی، تجربه تازه و تصمیم‌های لحظه‌ای برات جذابن. " +
					"وقتی فضا طبیعی و بدون فشار پیش می‌ره، بیشتر خودت می‌شی.",
			},
			{
				TraitKey: "dating:deep",
				Label:    "عمیق",
				Title:    "💬 جوینده اتصال عمیق",
				Subtitle: "قرار خوب برای تو یعنی یک گفت‌وگوی واقعی.",
				Description: "تو بیشتر از ظاهر برنامه، دنبال ارتباط ذهنی و عاطفی هستی. " +
					"اگر گفت‌وگو سطحی بمونه، حتی یک قرار لوکس هم ممکنه برات خالی باشه.",
			},
			{
				TraitKey: "dating:selective",
				Label:    "انتخاب‌گر",
				Title:    "🎯 انتخاب‌گر دقیق",
				Subtitle: "تو زود وارد بازی نمی‌شی؛ اول می‌سنجی.",
				Description: "نشانه‌های رفتاری، احترام و سازگاری برای تو مهمن. " +
					"ممکنه دیرتر جذب بشی، ولی وقتی انتخاب کنی معمولاً دلیل محکمی داری.",
			},
		},
		Questions: []questionSeed{
			{
				Text: "طرف مقابل می‌گه «امشب بریم بیرون؟» ولی هیچ برنامه‌ای نداره. واکنش تو؟",
				Options: []optionSeed{
					{Text: "اول می‌پرسم کجا و چه ساعتی؛ یه برنامه جمع‌وجور می‌چینم", TraitKey: "dating:planner"},
					{Text: "می‌گم بریم، تو راه تصمیم می‌گیریم!", TraitKey: "dating:spontaneous"},
					{Text: "اگر فرصت حرف زدن خوب باشه، مکان خیلی مهم نیست", TraitKey: "dating:deep"},
					{Text: "قبلش می‌خوام بدونم چقدر جدی و قابل‌اعتماده", TraitKey: "dating:selective"},
				},
			},
			{
				Text: "در اولین قرار، کدوم اتفاق بیشتر امتیاز می‌گیره؟",
				Options: []optionSeed{
					{Text: "به جزئیات فکر کرده و همه چیز مرتب پیش می‌ره", TraitKey: "dating:planner"},
					{Text: "یه اتفاق غیرمنتظره باحال پیش میاد", TraitKey: "dating:spontaneous"},
					{Text: "گفت‌وگو از سطح معمولی رد می‌شه", TraitKey: "dating:deep"},
					{Text: "رفتارش محترمانه و بدون Red Flagه", TraitKey: "dating:selective"},
				},
			},
			{
				Text: "بعد از یک قرار خوب، معمولاً دوست داری چه اتفاقی بیفته؟",
				Options: []optionSeed{
					{Text: "برای قرار بعدی یک زمان مشخص کنیم", TraitKey: "dating:planner"},
					{Text: "بذاریم حس و حال خودش جلو بره", TraitKey: "dating:spontaneous"},
					{Text: "یه پیام واقعی درباره حسی که داشتیم", TraitKey: "dating:deep"},
					{Text: "کمی فاصله تا ببینم رفتار بعدیش چیه", TraitKey: "dating:selective"},
				},
			},
			{
				Text: "اگر طرف مقابل نیم ساعت دیر کنه و درست توضیح بده، تو بیشتر کدومی؟",
				Options: []optionSeed{
					{Text: "می‌پذیرم ولی دوست دارم دفعه بعد هماهنگ‌تر باشه", TraitKey: "dating:planner"},
					{Text: "مشکلی نیست، برنامه رو عوض می‌کنیم", TraitKey: "dating:spontaneous"},
					{Text: "اگر توضیحش صادقانه باشه برام مهم‌تره", TraitKey: "dating:deep"},
					{Text: "این رفتار رو کنار بقیه نشانه‌ها یادم نگه می‌دارم", TraitKey: "dating:selective"},
				},
			},
			{
				Text: "کدوم قرار جذاب‌تره؟",
				Options: []optionSeed{
					{Text: "رستوران خوب با رزرو و برنامه مشخص", TraitKey: "dating:planner"},
					{Text: "یه مسیر بی‌برنامه، کافه تصادفی و قدم‌زدن", TraitKey: "dating:spontaneous"},
					{Text: "جایی آروم برای ساعت‌ها حرف زدن", TraitKey: "dating:deep"},
					{Text: "فعالیتی که بشه رفتار واقعی همدیگه رو دید", TraitKey: "dating:selective"},
				},
			},
			{
				Text: "طرف مقابل خیلی سریع ابراز علاقه شدید می‌کنه. واکنش تو؟",
				Options: []optionSeed{
					{Text: "سرعت رو تنظیم می‌کنم تا مرحله‌به‌مرحله جلو بریم", TraitKey: "dating:planner"},
					{Text: "اگر حسش خوب باشه، می‌ذارم اتفاق بیفته", TraitKey: "dating:spontaneous"},
					{Text: "می‌خوام بدونم پشت این احساس واقعاً چی هست", TraitKey: "dating:deep"},
					{Text: "کمی محتاط می‌شم تا ببینم حرف و عملش یکیه یا نه", TraitKey: "dating:selective"},
				},
			},
			{
				Text: "در انتخاب پارتنر، کدوم مورد برای تو سریع‌تر تصمیم‌سازه؟",
				Options: []optionSeed{
					{Text: "ثبات و قابل‌پیش‌بینی بودن مثبت", TraitKey: "dating:planner"},
					{Text: "هیجان و شیمی لحظه‌ای", TraitKey: "dating:spontaneous"},
					{Text: "عمق فکری و احساسی", TraitKey: "dating:deep"},
					{Text: "کیفیت شخصیت و مرزهای سالم", TraitKey: "dating:selective"},
				},
			},
			{
				Text: "اگر قرار اول «بد نبود ولی واو هم نبود»، چه می‌کنی؟",
				Options: []optionSeed{
					{Text: "اگر پتانسیل داشت، یک قرار دوم برنامه‌ریزی می‌کنم", TraitKey: "dating:planner"},
					{Text: "اگر حس لحظه‌ای نداشتم احتمالاً ادامه نمی‌دم", TraitKey: "dating:spontaneous"},
					{Text: "اگر گفت‌وگوی خوبی داشتیم، حتماً فرصت دوم می‌دم", TraitKey: "dating:deep"},
					{Text: "نکات مثبت و منفی رو سبک‌سنگین می‌کنم", TraitKey: "dating:selective"},
				},
			},
		},
	}
}

func personalBoundariesSeed() testSeedDefinition {
	return testSeedDefinition{
		CategoryCode: "love",
		Code:         "personal_boundaries_v1",
		Title:        "🛡 مرزهای شخصی",
		Description:  "ببین وقتی پای نه گفتن، احترام و فضای شخصی وسطه، سبک تو چیه.",
		SortOrder:    40,
		Results: []resultSeed{
			{
				TraitKey: "boundaries:firm",
				Label:    "روشن",
				Title:    "🛡 مرزبان روشن",
				Subtitle: "نه گفتن برای تو بی‌احترامی نیست.",
				Description: "معمولاً می‌دونی چه چیزی برات قابل‌قبوله و می‌تونی آن را نسبتاً شفاف بیان کنی. " +
					"قدرت تو وضوحه؛ فقط حواست باشه انعطاف رو با عقب‌نشینی اشتباه نگیری.",
			},
			{
				TraitKey: "boundaries:flexible",
				Label:    "منعطف",
				Title:    "🌿 منعطف اما آگاه",
				Subtitle: "هم مرز داری، هم جا برای مذاکره.",
				Description: "تو معمولاً شرایط و آدم‌ها را در نظر می‌گیری و بعد تصمیم می‌گیری. " +
					"این انعطاف می‌تونه نقطه قوت بزرگی باشه، تا وقتی نیازهای خودت گم نشن.",
			},
			{
				TraitKey: "boundaries:peacekeeper",
				Label:    "صلح‌طلب",
				Title:    "🤝 صلح‌طلبِ بیش‌ازحد",
				Subtitle: "گاهی آرام نگه داشتن فضا را به خودت ترجیح می‌دی.",
				Description: "تو از تنش خوشت نمیاد و ممکنه برای حفظ رابطه بیشتر از حد لازم کوتاه بیای. " +
					"حقیقت مهم برای تو اینه که صلح واقعی بدون احترام به مرزهای خودت دوام نمیاره.",
			},
			{
				TraitKey: "boundaries:guarded",
				Label:    "محافظ",
				Title:    "🚪 محافظِ محتاط",
				Subtitle: "ورود به فضای نزدیک تو ساده نیست.",
				Description: "وقتی احساس خطر یا فشار کنی سریع فاصله می‌گیری و کنترل فضای شخصی برات مهمه. " +
					"این محافظت ارزشمنده، اما گاهی توضیح روشن بهتر از قطع ارتباط ناگهانی عمل می‌کنه.",
			},
		},
		Questions: []questionSeed{
			{
				Text: "دوستی در آخرین لحظه ازت می‌خواد کاری انجام بدی که واقعاً وقتش رو نداری. چه می‌کنی؟",
				Options: []optionSeed{
					{Text: "محترمانه می‌گم نمی‌تونم", TraitKey: "boundaries:firm"},
					{Text: "اگر خیلی ضروری باشه راه میانه پیدا می‌کنم", TraitKey: "boundaries:flexible"},
					{Text: "احتمالاً قبول می‌کنم که ناراحت نشه", TraitKey: "boundaries:peacekeeper"},
					{Text: "اگر فشار بیاره، سریع فاصله می‌گیرم", TraitKey: "boundaries:guarded"},
				},
			},
			{
				Text: "پارتنرت رمز گوشی‌ات رو می‌خواد چون می‌گه «اگه چیزی نداری، چرا ندی؟»",
				Options: []optionSeed{
					{Text: "می‌گم حریم خصوصی ربطی به پنهان‌کاری نداره", TraitKey: "boundaries:firm"},
					{Text: "درباره دلیلش صحبت می‌کنم و مرز مشترک می‌سازیم", TraitKey: "boundaries:flexible"},
					{Text: "برای اینکه دعوا نشه شاید بدم", TraitKey: "boundaries:peacekeeper"},
					{Text: "این درخواست باعث می‌شه خیلی محتاط و بسته بشم", TraitKey: "boundaries:guarded"},
				},
			},
			{
				Text: "وقتی کسی شوخی‌ای می‌کنه که واقعاً اذیتت می‌کنه، واکنش طبیعی تو چیه؟",
				Options: []optionSeed{
					{Text: "همون موقع می‌گم این شوخی برام اوکی نیست", TraitKey: "boundaries:firm"},
					{Text: "با لحن آروم توضیح می‌دم چرا خوشم نیومد", TraitKey: "boundaries:flexible"},
					{Text: "می‌خندم ولی توی دلم می‌مونه", TraitKey: "boundaries:peacekeeper"},
					{Text: "کم‌کم از اون آدم فاصله می‌گیرم", TraitKey: "boundaries:guarded"},
				},
			},
			{
				Text: "یکی از نزدیکانت مرتب بدون هماهنگی وارد برنامه روزانه‌ات می‌شه.",
				Options: []optionSeed{
					{Text: "قانون مشخص می‌ذارم که قبلش هماهنگ کنه", TraitKey: "boundaries:firm"},
					{Text: "بعضی وقت‌ها اوکیه، بعضی وقت‌ها نه؛ توضیح می‌دم", TraitKey: "boundaries:flexible"},
					{Text: "سخت می‌تونم بگم مزاحمه", TraitKey: "boundaries:peacekeeper"},
					{Text: "ممکنه کلاً دسترسی‌اش رو خیلی محدود کنم", TraitKey: "boundaries:guarded"},
				},
			},
			{
				Text: "در یک بحث، طرف مقابل صدایش را بالا می‌برد. تو؟",
				Options: []optionSeed{
					{Text: "می‌گم با این لحن ادامه نمی‌دم", TraitKey: "boundaries:firm"},
					{Text: "سعی می‌کنم فضا رو آروم کنم و بعد ادامه بدیم", TraitKey: "boundaries:flexible"},
					{Text: "برای تمام شدن دعوا زودتر کوتاه میام", TraitKey: "boundaries:peacekeeper"},
					{Text: "بحث رو قطع می‌کنم و کامل عقب می‌کشم", TraitKey: "boundaries:guarded"},
				},
			},
			{
				Text: "وقتی ازت ناراحت می‌شن چون «نه» گفتی، چه حسی پیدا می‌کنی؟",
				Options: []optionSeed{
					{Text: "ناراحت می‌شم ولی تصمیمم رو عوض نمی‌کنم", TraitKey: "boundaries:firm"},
					{Text: "دوباره بررسی می‌کنم که آیا راه منصفانه‌تری هست", TraitKey: "boundaries:flexible"},
					{Text: "احساس گناه می‌کنم و ممکنه نظرم عوض شه", TraitKey: "boundaries:peacekeeper"},
					{Text: "فکر می‌کنم بهتره مدتی فاصله بگیرم", TraitKey: "boundaries:guarded"},
				},
			},
			{
				Text: "کدام جمله بیشتر شبیه توست؟",
				Options: []optionSeed{
					{Text: "«دوستت دارم، ولی این برام قابل‌قبول نیست.»", TraitKey: "boundaries:firm"},
					{Text: "«بیا یه راهی پیدا کنیم که برای هر دومون خوب باشه.»", TraitKey: "boundaries:flexible"},
					{Text: "«ولش کن، مهم نیست.» حتی وقتی مهمه", TraitKey: "boundaries:peacekeeper"},
					{Text: "«اگر مجبورم کنی، کلاً کنار می‌کشم.»", TraitKey: "boundaries:guarded"},
				},
			},
			{
				Text: "بیشتر از چی می‌ترسی؟",
				Options: []optionSeed{
					{Text: "اینکه کسی مرزم رو جدی نگیره", TraitKey: "boundaries:firm"},
					{Text: "اینکه نتونیم به توافق منصفانه برسیم", TraitKey: "boundaries:flexible"},
					{Text: "اینکه نه گفتن باعث از دست دادن آدم‌ها بشه", TraitKey: "boundaries:peacekeeper"},
					{Text: "اینکه زیادی نزدیک شدن باعث آسیب بشه", TraitKey: "boundaries:guarded"},
				},
			},
		},
	}
}

func bitterTruthSeed() testSeedDefinition {
	return testSeedDefinition{
		CategoryCode: "challenge",
		Code:         "bitter_truth_v1",
		Title:        "🪞 حقیقت تلخ",
		Description:  "یک تست چالشی برای پیدا کردن الگویی که شاید درباره خودت کمتر دوست داشته باشی ببینی.",
		SortOrder:    50,
		Results: []resultSeed{
			{
				TraitKey: "truth:approval",
				Label:    "تأییدطلب",
				Title:    "🎭 حقیقت تلخ: زیادی دنبال تأییدی",
				Subtitle: "گاهی نظر بقیه صدای خودت را کم‌رنگ می‌کنه.",
				Description: "ممکنه بیشتر از چیزی که فکر می‌کنی واکنش دیگران روی تصمیم‌ها و حال خوبت اثر بذاره. " +
					"نقطه رشد تو اینه که قبل از پرسیدن «اونا چی فکر می‌کنن؟» بپرسی «من واقعاً چی می‌خوام؟».",
			},
			{
				TraitKey: "truth:control",
				Label:    "کنترل‌گر",
				Title:    "🎛 حقیقت تلخ: دوست داری کنترل دست تو باشه",
				Subtitle: "ابهام و بی‌برنامگی سریع‌تر از بقیه خسته‌ات می‌کنه.",
				Description: "توان مدیریت و پیش‌بینی تو نقطه قوته، اما وقتی همه‌چیز باید طبق نقشه تو جلو بره، " +
					"ممکنه آزادی و خودجوشی دیگران کمتر جا داشته باشه.",
			},
			{
				TraitKey: "truth:avoidance",
				Label:    "اجتنابی",
				Title:    "🌫 حقیقت تلخ: بعضی چیزها را عقب می‌اندازی تا مجبور نشی باهاشون روبه‌رو شی",
				Subtitle: "سکوت کوتاه‌مدت آرامش می‌ده، ولی مسئله همیشه محو نمی‌شه.",
				Description: "تو احتمالاً در تنش‌ها اول دنبال کم‌کردن فشار هستی. " +
					"گاهی همین باعث می‌شه گفت‌وگوی لازم دیرتر اتفاق بیفته و مسئله بزرگ‌تر شه.",
			},
			{
				TraitKey: "truth:intensity",
				Label:    "همه یا هیچ",
				Title:    "🔥 حقیقت تلخ: گاهی همه‌چیز برایت صفر یا صده",
				Subtitle: "وقتی چیزی برات مهم شه، نصفه‌نیمه بودن سخت می‌شه.",
				Description: "شدت احساس و تعهدت می‌تونه خیلی جذاب باشه، اما نگاه همه‌یا‌هیچ گاهی اجازه نمی‌ده خاکستری‌ها و تغییر تدریجی رو ببینی.",
			},
		},
		Questions: []questionSeed{
			{
				Text: "یک نفر پیام تو را دیده ولی چند ساعت جواب نداده. اولین داستانی که ذهنت می‌سازه چیه؟",
				Options: []optionSeed{
					{Text: "نکنه چیزی گفتم که بد برداشت کرده؟", TraitKey: "truth:approval"},
					{Text: "دوست دارم بدونم دقیقاً چرا و کی جواب می‌ده", TraitKey: "truth:control"},
					{Text: "بی‌خیال، منم فعلاً جواب نمی‌دم", TraitKey: "truth:avoidance"},
					{Text: "یا علاقه داره یا نداره؛ این وسط‌بازی‌ها رو دوست ندارم", TraitKey: "truth:intensity"},
				},
			},
			{
				Text: "وقتی کسی ازت انتقاد می‌کنه، کدوم واکنش بهت نزدیک‌تره؟",
				Options: []optionSeed{
					{Text: "مدتی ذهنم درگیر می‌شه که نکنه واقعاً بد دیده شدم", TraitKey: "truth:approval"},
					{Text: "سریع دنبال دلیل و منطق می‌گردم تا بفهمم حق با کیه", TraitKey: "truth:control"},
					{Text: "ترجیح می‌دم بحث رو عوض کنم یا بعداً بهش فکر کنم", TraitKey: "truth:avoidance"},
					{Text: "اگر حس کنم ناعادلانه‌ست، خیلی شدید واکنش می‌گیرم", TraitKey: "truth:intensity"},
				},
			},
			{
				Text: "وقتی برنامه‌ای که برایش ذوق داشتی ناگهان عوض می‌شود، بیشتر چه چیزی اذیتت می‌کند؟",
				Options: []optionSeed{
					{Text: "اینکه نکنه من برای بقیه به اندازه کافی مهم نبودم", TraitKey: "truth:approval"},
					{Text: "اینکه کنترل و نظم ماجرا از دست رفته", TraitKey: "truth:control"},
					{Text: "ترجیح می‌دم کلاً بی‌خیالش بشم", TraitKey: "truth:avoidance"},
					{Text: "خیلی سریع از ذوق کامل به ناامیدی کامل می‌رسم", TraitKey: "truth:intensity"},
				},
			},
			{
				Text: "کدام رفتار خودت بیشتر ممکنه بعداً حرصت بده؟",
				Options: []optionSeed{
					{Text: "اینکه برای خوشحال کردن بقیه چیزی رو قبول کردم", TraitKey: "truth:approval"},
					{Text: "اینکه زیادی روی جزئیات و نتیجه پافشاری کردم", TraitKey: "truth:control"},
					{Text: "اینکه حرف لازم رو نزدم و گذاشتم بگذره", TraitKey: "truth:avoidance"},
					{Text: "اینکه از روی احساس تصمیم خیلی قطعی گرفتم", TraitKey: "truth:intensity"},
				},
			},
			{
				Text: "در یک اختلاف رابطه‌ای، کدوم جمله بیشتر شبیه ذهن توئه؟",
				Options: []optionSeed{
					{Text: "«فقط می‌خوام مطمئن شم هنوز دوستم داره.»", TraitKey: "truth:approval"},
					{Text: "«باید مشخص کنیم دقیقاً قرار چطور ادامه بدیم.»", TraitKey: "truth:control"},
					{Text: "«الان حوصله این بحث رو ندارم؛ بعداً.»", TraitKey: "truth:avoidance"},
					{Text: "«اگر اینطوریه، شاید اصلاً نباید ادامه بدیم.»", TraitKey: "truth:intensity"},
				},
			},
			{
				Text: "وقتی تصمیم مهمی می‌گیری، کدوم دام بیشتر ممکنه گیرت بندازه؟",
				Options: []optionSeed{
					{Text: "نظر همه رو می‌پرسم و خودم گیج‌تر می‌شم", TraitKey: "truth:approval"},
					{Text: "آنقدر اطلاعات جمع می‌کنم که تصمیم عقب می‌افته", TraitKey: "truth:control"},
					{Text: "تا وقتی مجبور نشم تصمیم رو عقب می‌ندازم", TraitKey: "truth:avoidance"},
					{Text: "یک لحظه کاملاً مطمئن می‌شم و سریع می‌پرم", TraitKey: "truth:intensity"},
				},
			},
			{
				Text: "کدام تعریف پنهانی بیشتر خوشحالت می‌کند؟",
				Options: []optionSeed{
					{Text: "«همه دوستت دارن.»", TraitKey: "truth:approval"},
					{Text: "«تو همیشه می‌دونی باید چی کار کرد.»", TraitKey: "truth:control"},
					{Text: "«تو هیچ‌وقت وارد drama نمی‌شی.»", TraitKey: "truth:avoidance"},
					{Text: "«تو برای چیزهایی که می‌خوای واقعاً می‌جنگی.»", TraitKey: "truth:intensity"},
				},
			},
			{
				Text: "اگر بخوای یک عادتت رو همین امروز کمتر کنی، کدوم بیشتر به کارت میاد؟",
				Options: []optionSeed{
					{Text: "کمتر دنبال رضایت همه باشم", TraitKey: "truth:approval"},
					{Text: "بذارم بعضی چیزها بدون کنترل من جلو بره", TraitKey: "truth:control"},
					{Text: "گفت‌وگوهای سخت رو عقب نندازم", TraitKey: "truth:avoidance"},
					{Text: "قبل از تصمیم قطعی، کمی خاکستری‌ها رو ببینم", TraitKey: "truth:intensity"},
				},
			},
		},
	}
}
