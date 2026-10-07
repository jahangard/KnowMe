package database

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type plannedTestSeed struct {
	categoryCode string
	code         string
	title        string
	description  string
	sortOrder    int
}

func seedPlannedTests(ctx context.Context, db *gorm.DB) error {
	tests := []plannedTestSeed{
		{categoryCode: "personality-assessments", code: "mbti_v1", title: "تست شخصیت‌شناسی MBTI", description: "روش اجرا: پرسش‌نامه‌ی خودشناسی بر پایه‌ی چهار ترجیح درون‌گرایی/برون‌گرایی، حسی/شهودی، منطقی/احساسی و قضاوت‌گر/ادراکی؛ نتیجه یکی از ۱۶ تیپ است. ۶۰ سؤال، حدود ۱۵ دقیقه.", sortOrder: 10},
		{categoryCode: "personality-assessments", code: "disc_v1", title: "تست شخصیت DISC", description: "روش اجرا: ترجیح فرد در دو محورِ کار در برابر افراد و سرعت بالا در برابر متوسط بررسی می‌شود؛ خروجی چهار گرایش D، I، S و C را نشان می‌دهد.", sortOrder: 20},
		{categoryCode: "personality-assessments", code: "cattell_16_v1", title: "تست شخصیت ۱۶ عاملی کتل", description: "شناخت چندبعدی ویژگی‌های شخصیتی بر پایه‌ی ۱۶ عامل کتل.", sortOrder: 30},
		{categoryCode: "personality-assessments", code: "neo_ffi_v1", title: "تست شخصیت نئو (NEO-FFI)", description: "روش اجرا: پرسش‌نامه‌ی خودگزارشی پنج عامل روان‌رنجوری، برون‌گرایی، توافق‌پذیری، گشودگی به تجربه و وظیفه‌شناسی را می‌سنجد. ۶۰ سؤال، حدود ۱۵ دقیقه.", sortOrder: 40},
		{categoryCode: "personality-assessments", code: "hartman_v1", title: "تست شخصیت هارتمن", description: "روش اجرا: پرسش‌نامه‌ی ارزش‌های محوری؛ نتیجه در چهار رنگ قرمز، آبی، سفید و زرد دسته‌بندی می‌شود. ۴۵ سؤال، حدود ۱۰ دقیقه.", sortOrder: 50},

		{categoryCode: "career-self-knowledge", code: "organizational_commitment_v1", title: "تست تعهد سازمانی", description: "میزان پیوند و تعهد فرد به سازمان و محیط کاری را بررسی می‌کند.", sortOrder: 10},
		{categoryCode: "career-self-knowledge", code: "job_satisfaction_v1", title: "تست رضایت شغلی", description: "نگاه فرد به وظایف، محیط کار و تجربه‌ی شغلی را ارزیابی می‌کند.", sortOrder: 20},
		{categoryCode: "career-self-knowledge", code: "job_burnout_v1", title: "تست فرسودگی شغلی", description: "نشانه‌های خستگی و فشار مزمن در تجربه‌ی کاری را بررسی می‌کند.", sortOrder: 30},
		{categoryCode: "career-self-knowledge", code: "entrepreneur_personality_v1", title: "تست شخصیت کارآفرین", description: "گرایش‌ها و توانمندی‌های شخصیتی مرتبط با کارآفرینی را بررسی می‌کند.", sortOrder: 40},
		{categoryCode: "career-self-knowledge", code: "disc_career_v1", title: "تست شخصیت DISC", description: "روش اجرا: ترجیح فرد در دو محورِ کار در برابر افراد و سرعت بالا در برابر متوسط بررسی می‌شود؛ خروجی چهار گرایش D، I، S و C را نشان می‌دهد.", sortOrder: 50},
		{categoryCode: "career-self-knowledge", code: "mbti_career_v1", title: "تست شخصیت‌شناسی MBTI", description: "روش اجرا: پرسش‌نامه‌ی خودشناسی بر پایه‌ی چهار ترجیح رفتاری؛ نتیجه یکی از ۱۶ تیپ شخصیتی است. ۶۰ سؤال، حدود ۱۵ دقیقه.", sortOrder: 60},
		{categoryCode: "career-self-knowledge", code: "vocational_interest_v1", title: "تست رغبت‌سنج شغلی", description: "زمینه‌های کاری و فعالیت‌هایی را که با علاقه‌ها و ترجیحات فرد سازگارند بررسی می‌کند.", sortOrder: 70},
		{categoryCode: "career-self-knowledge", code: "neo_ffi_career_v1", title: "تست شخصیت نئو (NEO-FFI)", description: "روش اجرا: پرسش‌نامه‌ی خودگزارشی پنج عامل شخصیتی را برای شناخت تناسب‌های حرفه‌ای بررسی می‌کند. ۶۰ سؤال، حدود ۱۵ دقیقه.", sortOrder: 80},
		{categoryCode: "career-self-knowledge", code: "adapto_career_path_v1", title: "تست مسیر شغلی ادپتو (Career Path)", description: "روش اجرا: چهار بُعد شخصیت شغلی، علایق، استعدادها و ارزش‌های شغلی را ترکیب می‌کند و مسیرها و شغل‌های متناسب پیشنهاد می‌دهد. چهار ارزیابی، حدود ۴۰ دقیقه.", sortOrder: 90},
		{categoryCode: "career-self-knowledge", code: "dmsi_motivation_v1", title: "تست سبک انگیزشی غالب (DMSI)", description: "روش اجرا: بر پایه‌ی سه محرک مک‌کللند؛ پیشرفت، قدرت و پیوندجویی. ۱۵ سؤال.", sortOrder: 100},
		{categoryCode: "career-self-knowledge", code: "tki_conflict_style_v1", title: "تست سبک حل تعارض توماس–کیلمن (TKI)", description: "روش اجرا: پرسش‌نامه‌ی موقعیت‌محور؛ سبک‌های همکاری، رقابت، مصالحه، سازش و اجتناب را نشان می‌دهد. ۳۰ سؤال.", sortOrder: 110},
		{categoryCode: "career-self-knowledge", code: "csi_stress_coping_v1", title: "تست سبک‌های مقابله با استرس (CSI)", description: "روش اجرا: پرسش‌نامه‌ی خودگزارشی درباره‌ی واکنش به موقعیت‌های استرس‌زا؛ هشت سبک مقابله را بررسی می‌کند. ۷۲ سؤال، حدود ۱۰ دقیقه.", sortOrder: 120},
		{categoryCode: "career-self-knowledge", code: "cognitive_abilities_career_v1", title: "تست توانایی‌های شناختی", description: "روش اجرا: ۳۰ سؤال برای ارزیابی حوزه‌هایی مانند توجه، برنامه‌ریزی، تصمیم‌گیری، حافظه و انعطاف‌پذیری شناختی.", sortOrder: 130},
		{categoryCode: "career-self-knowledge", code: "career_anchors_v1", title: "تست لنگرگاه‌های شغلی (Career Anchors)", description: "روش اجرا: خودارزیابی شایستگی‌ها، ارزش‌ها و انگیزه‌ها برای شناخت گرایش‌های اصلی مسیر شغلی. ۴۰ سؤال، حدود ۱۰ دقیقه.", sortOrder: 140},
		{categoryCode: "career-self-knowledge", code: "clifton_strengths_v1", title: "تست استعدادیابی کلیفتون (گالوپ)", description: "روش اجرا: ۳۴ استعداد را رتبه‌بندی می‌کند و نقاط قوت برجسته و کاربردهای حرفه‌ای آن‌ها را گزارش می‌دهد.", sortOrder: 150},
		{categoryCode: "career-self-knowledge", code: "holland_vocational_interest_v1", title: "تست رغبت شغلی هالند (RIASEC)", description: "روش اجرا: علاقه‌های شغلی را در شش تیپ واقع‌گرا، جست‌وجوگر، هنری، اجتماعی، متهور و قراردادی بررسی می‌کند.", sortOrder: 160},
		{categoryCode: "career-self-knowledge", code: "neo_pir_v1", title: "تست شخصیت نئو (NEO PI-R)", description: "روش اجرا: پرسش‌نامه‌ی گسترده‌ی پنج‌عاملی شخصیت. ۲۴۰ سؤال، حدود ۶۰ دقیقه.", sortOrder: 170},

		{categoryCode: "intelligence-self-knowledge", code: "raven_adult_v1", title: "تست هوش ریون بزرگسالان", description: "روش اجرا: ۶۰ پرسش تصویری چندگزینه‌ای از آسان به دشوار؛ با کامل‌کردن ماتریس‌های هندسی استدلال غیرکلامی سنجیده می‌شود. حدود ۴۵ دقیقه.", sortOrder: 10},
		{categoryCode: "intelligence-self-knowledge", code: "gardner_multiple_intelligences_v1", title: "تست هوش‌های چندگانه گاردنر", description: "توانمندی‌های فرد را در حوزه‌های مختلف هوش بررسی می‌کند.", sortOrder: 20},
		{categoryCode: "intelligence-self-knowledge", code: "cattell_adult_a_v1", title: "تست هوش کتل بزرگسال A", description: "فرم A برای بررسی استدلال و توانایی شناختی بزرگسالان.", sortOrder: 30},
		{categoryCode: "intelligence-self-knowledge", code: "cattell_adult_b_v1", title: "تست هوش کتل بزرگسال B", description: "فرم B برای بررسی استدلال و توانایی شناختی بزرگسالان.", sortOrder: 40},
		{categoryCode: "intelligence-self-knowledge", code: "bonnardel_v1", title: "تست هوش بوناردل", description: "توانایی تحلیل و استدلال در الگوهای تصویری را بررسی می‌کند.", sortOrder: 50},
		{categoryCode: "intelligence-self-knowledge", code: "baron_emotional_intelligence_v1", title: "تست هوش هیجانی بارآن", description: "روش اجرا: پرسش‌نامه‌ی خودگزارشی در ۱۵ بُعد هیجانی و اجتماعی. ۹۰ سؤال، حدود ۲۵ دقیقه.", sortOrder: 60},
		{categoryCode: "intelligence-self-knowledge", code: "schutte_emotional_intelligence_v1", title: "تست هوش هیجانی شات", description: "برداشت و مدیریت هیجان‌ها در خود و دیگران را ارزیابی می‌کند.", sortOrder: 70},
		{categoryCode: "intelligence-self-knowledge", code: "child_intelligence_v1", title: "تست هوش کودک", description: "ارزیابی متناسب با سن کودک؛ گروه سنی و شیوه‌ی اجرا پیش از انتشار تکمیل می‌شود.", sortOrder: 80},
		{categoryCode: "intelligence-self-knowledge", code: "cognitive_abilities_v1", title: "تست توانایی‌های شناختی", description: "روش اجرا: ۳۰ سؤال برای ارزیابی حوزه‌هایی مانند توجه، برنامه‌ریزی، تصمیم‌گیری، حافظه و انعطاف‌پذیری شناختی.", sortOrder: 90},

		{categoryCode: "love-marriage-self-knowledge", code: "sternberg_triangular_love_v1", title: "تست مثلث عشق استرنبرگ", description: "صمیمیت، شور و تعهد را در تجربه‌ی رابطه بررسی می‌کند.", sortOrder: 10},
		{categoryCode: "love-marriage-self-knowledge", code: "premarital_fears_v1", title: "تست ترس‌های قبل از ازدواج", description: "نگرانی‌ها و دغدغه‌های پیش از تصمیم به ازدواج را بررسی می‌کند.", sortOrder: 20},
		{categoryCode: "love-marriage-self-knowledge", code: "enrich_marital_satisfaction_v1", title: "تست رضایت زناشویی انریچ (ENRICH)", description: "ابعاد گوناگون رضایت و کیفیت رابطه‌ی زناشویی را ارزیابی می‌کند.", sortOrder: 30},
		{categoryCode: "love-marriage-self-knowledge", code: "cattell_16_love_v1", title: "تست شخصیت ۱۶ عاملی کتل", description: "نگاهی چندبعدی به ویژگی‌های شخصیتی و تفاوت‌های فردی.", sortOrder: 40},

		{categoryCode: "social-self-knowledge", code: "emotional_intelligence_v1", title: "تست هوش هیجانی", description: "شناخت و مدیریت هیجان‌ها را در زندگی روزمره بررسی می‌کند.", sortOrder: 10},
		{categoryCode: "social-self-knowledge", code: "loneliness_v1", title: "تست احساس تنهایی", description: "تجربه‌ی تنهایی و احساس پیوند اجتماعی را بررسی می‌کند.", sortOrder: 20},
		{categoryCode: "social-self-knowledge", code: "premarital_fears_social_v1", title: "تست ترس‌های قبل از ازدواج", description: "نگرانی‌ها و دغدغه‌های پیش از تصمیم به ازدواج را بررسی می‌کند.", sortOrder: 30},
		{categoryCode: "social-self-knowledge", code: "general_health_v1", title: "تست سلامت عمومی (GHQ)", description: "روش اجرا: پرسش‌نامه‌ی خودگزارشی درباره‌ی تجربه‌ی سلامت عمومی در یک ماه گذشته.", sortOrder: 40},
		{categoryCode: "social-self-knowledge", code: "social_anxiety_v1", title: "تست اضطراب اجتماعی", description: "میزان نگرانی و دشواری در موقعیت‌های اجتماعی را بررسی می‌کند.", sortOrder: 50},
		{categoryCode: "social-self-knowledge", code: "depression_v1", title: "تست افسردگی", description: "نشانه‌های خلقی را برای خودآگاهی اولیه بررسی می‌کند.", sortOrder: 60},
		{categoryCode: "social-self-knowledge", code: "tki_conflict_style_social_v1", title: "تست سبک حل تعارض توماس–کیلمن (TKI)", description: "روش اجرا: پرسش‌نامه‌ی موقعیت‌محور؛ پنج سبک واکنش در تعارض را نشان می‌دهد. ۳۰ سؤال.", sortOrder: 70},
		{categoryCode: "social-self-knowledge", code: "csi_stress_coping_social_v1", title: "تست سبک‌های مقابله با استرس (CSI)", description: "روش اجرا: خودگزارشی درباره‌ی واکنش به فشار؛ هشت سبک مقابله را بررسی می‌کند. ۷۲ سؤال، حدود ۱۰ دقیقه.", sortOrder: 80},
	}

	return db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, definition := range tests {
			var category TestCategory
			if err := tx.Where("Code = ? AND IsActive = ?", definition.categoryCode, true).First(&category).Error; err != nil {
				return fmt.Errorf("load category %s: %w", definition.categoryCode, err)
			}

			description := definition.description
			planned := Test{
				CategoryID:  category.ID,
				Code:        definition.code,
				Title:       definition.title,
				Description: &description,
				SortOrder:   definition.sortOrder,
				IsActive:    true,
				IsReady:     false,
			}
			var existing Test
			if err := tx.Where("Code = ?", definition.code).Assign(planned).FirstOrCreate(&existing).Error; err != nil {
				return fmt.Errorf("upsert planned test %s: %w", definition.code, err)
			}
			if err := tx.Model(&Test{}).Where("Code = ?", definition.code).Updates(map[string]any{
				"CategoryId":  category.ID,
				"Title":       definition.title,
				"Description": definition.description,
				"SortOrder":   definition.sortOrder,
				"IsActive":    true,
				"IsReady":     false,
			}).Error; err != nil {
				return fmt.Errorf("update planned test %s: %w", definition.code, err)
			}
		}
		return nil
	})
}
