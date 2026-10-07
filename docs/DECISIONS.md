# KnowMe Decision Log

> تصمیم‌های اصلی محصول و معماری. هدف این فایل ثبت «چرا»ی تصمیم‌هاست تا بعداً دوباره از صفر بحث نشوند.

## D-001 — Go برای Runtime

**Status:** Accepted

Go برای Bot انتخاب شد چون:

- مصرف حافظه کم
- binary مستقل
- مناسب Windows Service
- مناسب Linux/ARM
- deployment ساده
- concurrency مناسب برای Bot

## D-002 — Telegram Long Polling برای MVP

**Status:** Accepted

Webhook در MVP لازم نیست.

دلیل:

- IP ثابت لازم نیست
- inbound port لازم نیست
- setup ساده روی Windows workstation
- مناسب حجم اولیه محصول

## D-003 — SQLite برای Development، SQL Server برای Production

**Status:** Accepted

Development باید بدون وابستگی به Server خارجی ممکن باشد.

SQLite سرعت توسعه محلی را بالا می‌برد و SQL Server همچنان Provider Production است.

## D-004 — GORM به‌عنوان ORM

**Status:** Accepted

هدف اصلی:

- یک persistence model برای SQLite و SQL Server
- حذف duplication کوئری‌های provider-specific
- AutoMigrate برای schema evolution
- ساده‌تر شدن توسعه Feature

## D-005 — Database-first Test Content

**Status:** Accepted

تعریف Test نباید در Bot code باشد.

Database شامل:

- Category
- Test
- Question
- Option
- Score
- Trait
- Result Profile

باشد.

نتیجه:

افزودن Test استاندارد نیاز به Deploy کد ندارد.

## D-006 — Seed فقط Bootstrap است

**Status:** Accepted

Seed:

- Missing row را ایجاد می‌کند
- Existing content را overwrite نمی‌کند

دلیل:

اگر DB قرار است Source of Truth باشد، Restart برنامه نباید Editing محتوا را خنثی کند.

## D-007 — Result Presentation در DB

**Status:** Accepted

جدول `TestResultProfiles` اضافه شد.

قبل از این، Result Presentation برای Love Style در Telegram UI Hard-code بود.

اکنون Result title/subtitle/description/label از DB خوانده می‌شود.

## D-008 — مدل Scoring اولیه: Dominant Trait

**Status:** Accepted for MVP

هر Option یک TraitKey و Score دارد.

Scoreها Group می‌شوند و Trait غالب Result اصلی است.

این مدل برای اکثر Testهای اولیه کافی است.

Rule Engine پیچیده‌تر بعداً در صورت نیاز اضافه می‌شود.

## D-009 — Trait Namespace

**Status:** Accepted

Testهای جدید باید TraitKey namespaced داشته باشند.

مثال:

- `attraction:magnetic`
- `dating:planner`
- `boundaries:firm`

دلیل:

جلوگیری از Collision بین Testها.

## D-010 — دو مسیر Test Discovery

**Status:** Accepted

User باید همزمان دو تجربه داشته باشد:

- Free browsing by topic
- Personalized roadmap

Roadmap جای Category Browser را نمی‌گیرد و برعکس.

## D-011 — Roadmap یک Virtual Entry است

**Status:** Accepted

«نقشه راه من» کنار Categoryها نمایش داده می‌شود، اما TestCategory واقعی نیست.

دلیل:

Roadmap behavior یک Recommendation workflow است، نه Content category.

## D-012 — Category Browser از DB

**Status:** Accepted

Category و Test list از DB Load می‌شوند.

در صورت اضافه شدن Test به Category موجود، UI بدون تغییر کد آن را نمایش می‌دهد.

برای Category جدید، Generic icon قابل استفاده است؛ Icon اختصاصی یک enhancement نمایشی است و نباید مانع نمایش Category شود.

## D-013 — 18+ Age Gate

**Status:** Accepted

Category ۱۸+ فقط برای Age >= 18 قابل مشاهده/ورود است.

نتیجه‌های این بخش باید طیفی و غیرقضاوتی باشند.

## D-014 — Progressive Profiling

**Status:** Accepted

User نباید هنگام ورود فرم بلند ببیند.

اطلاعات به‌مرور و Contextual گرفته می‌شوند.

## D-015 — Edit Message برای Navigation

**Status:** Accepted

برای جلوگیری از شلوغی Telegram Chat، Inline flow تا جای ممکن با Edit Message انجام می‌شود.

## D-016 — Raw Answers همیشه ذخیره شوند

**Status:** Accepted

TestResult و UserTrait مشتق‌شده‌اند.

Answer history داده اصلی و قابل بازپردازش است.

## D-017 — Structured Persistent Logging

**Status:** Accepted

Terminal به‌تنهایی برای Error investigation کافی نیست.

Logging باید:

- structured
- persistent
- rotated
- secret-safe

باشد.

## D-018 — پنج Test اولیه

**Status:** Accepted

ترتیب اولیه Roadmap:

1. سبک عشق‌ورزی
2. تیپ جذابیت
3. سناریوهای قرار
4. مرزهای شخصی
5. حقیقت تلخ

## D-019 — Testها برای همه جنسیت‌ها تا جای ممکن مشترک

**Status:** Accepted

Gender-specific branching فقط در صورت نیاز واقعی Content استفاده شود.

## D-020 — Entertainment tests تشخیص نیستند

**Status:** Accepted

KnowMe می‌تواند Test سرگرمی و self-discovery داشته باشد، اما نباید آن‌ها را Diagnosis پزشکی/روان‌شناختی معرفی کند.

## Future decisions

مواردی که هنوز نیاز به تصمیم نهایی دارند:

- Admin/Test Studio
- Branching Question model
- Multi-select support
- adaptive question engine
- result share-card renderer
- consent/privacy/delete-data flow
- mobile collection timing
- advanced Roadmap scoring
- event taxonomy finalization
