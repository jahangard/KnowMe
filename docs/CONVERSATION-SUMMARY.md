# KnowMe — جمع‌بندی تصمیم‌ها و نتایج طراحی

> این سند جمع‌بندی مکالمات طراحی و توسعه پروژه KnowMe تا 2026-10-07 است و به‌عنوان مرجع تصمیم‌های محصول و فنی نگهداری می‌شود.

## 1. هدف محصول

KnowMe یک تجربه تعاملی مبتنی بر Telegram برای تست‌های سرگرم‌کننده و خودشناسی است. سیستم باید هم برای زن‌ها و هم مردها جذاب باشد، نتیجه قابل‌فهم و قابل‌اشتراک تولید کند و به‌مرور از پاسخ‌ها برای ساخت پروفایل و پیشنهاد مسیر بعدی استفاده کند.

دو شیوه اصلی استفاده از تست‌ها داریم:

1. **نقشه راه من** — سیستم بر اساس وضعیت کاربر، تست بعدی را خودش پیشنهاد می‌کند.
2. **انتخاب آزاد تست** — کاربر وارد فهرست دسته‌بندی‌ها می‌شود و تست دلخواهش را انتخاب می‌کند.

«نقشه راه من» در UI مانند یک دسته اصلی نمایش داده می‌شود، ولی رفتار آن با دسته‌های عادی متفاوت و هوشمند است.

## 2. ساختار صفحه اصلی تست‌ها

دسته‌بندی‌های توافق‌شده:

- 🧭 **نقشه راه من** — انتخاب خودکار تست بعدی متناسب با مسیر کاربر
- ❤️ **عشق و رابطه**
- 🎯 **چالشی و باحال**
- 🧠 **خودشناسی عمیق‌تر**
- 🔞 **۱۸+** — فقط برای کاربران واجد شرایط سنی

صفحه اصلی باید ابتدا دسته‌ها را نشان دهد. ورود به دسته عادی، فهرست تست‌های آن دسته را نمایش می‌دهد؛ ورود به «نقشه راه من» مستقیماً Recommendation Engine را اجرا می‌کند.

## 3. تست اولیه

اولین تست Seed شده فعلی:

**💞 سبک عشق‌ورزی** (`love_style_v1`)

این تست برای هر دو جنس طراحی شده و پاسخ‌ها را در چهار Trait اصلی امتیازدهی می‌کند:

- `words` — ابراز علاقه با کلام
- `time` — وقت باکیفیت
- `care` — توجه و عمل
- `touch` — نزدیکی و تماس

تست فعلی ۸ سؤال Single Choice دارد و نتیجه غالب از مجموع امتیاز Traitها تعیین می‌شود.

## 4. تجربه کاربری و Progressive Profiling

هدف این است که در شروع، کاربر با فرم طولانی مواجه نشود.

جریان فعلی:

- ورود کاربر از Telegram
- دریافت جنسیت
- دریافت سن
- ورود به Home / Roadmap
- بعد از حداقل یک تست تکمیل‌شده، درخواست اختیاری نام/لقب
- اطلاعات بیشتر فقط در زمانی گرفته شود که برای تجربه کاربر ارزش داشته باشد.

فیلدهای Profile فعلی شامل Name، Age، Gender، Mobile، ProfileCompletionLevel، PendingField و وضعیت Promptها هستند.

## 5. Telegram UX

Bot بر پایه Long Polling کار می‌کند.

منوی اصلی شامل مسیرهایی مانند:

- خانه
- نقشه راه
- پروفایل
- درباره

برای تست‌ها از Inline Keyboard و Callback Query استفاده می‌شود. سؤال بعدی در همان پیام Edit می‌شود تا چت بی‌دلیل شلوغ نشود.

هدف UI بعدی: اضافه‌شدن ورودی «تست‌ها» به صفحه اصلی و نمایش Category Browser مطابق ساختار بخش ۲.

## 6. معماری داده

مدل‌های اصلی فعلی:

- `Users`
- `UserProfiles`
- `TestCategories`
- `Tests`
- `Questions`
- `QuestionOptions`
- `TestSessions`
- `TestAnswers`
- `TestResults`
- `UserTraits`
- `UserEvents`

هر Test به Category وابسته است. Session وضعیت اجرای تست را نگه می‌دارد، Answer پاسخ‌های کاربر را ثبت می‌کند، Result خروجی نهایی را نگه می‌دارد و UserTrait شناخت تجمعی از کاربر را تشکیل می‌دهد.

## 7. Database و ORM

تصمیم نهایی این است که برنامه هر دو Provider را پشتیبانی کند:

### Development
- SQLite
- مسیر پیش‌فرض: `data/knowme.db`
- WAL
- Foreign Keys فعال
- Busy timeout
- بدون نیاز به SQL Server برای توسعه روزمره

### Production
- SQL Server
- اتصال از طریق `SQLSERVER_DSN`

ORM انتخاب‌شده: **GORM**

برنامه در Startup از `AutoMigrate` استفاده می‌کند و Seed اولیه نیز Idempotent است. Provider از طریق `DB_PROVIDER=sqlite|sqlserver` انتخاب می‌شود.

در Development، SQLite انتخاب پیش‌فرض است.

## 8. Logging و Observability

نیاز مهم پروژه این بود که خطاهایی که در Terminal دیده نمی‌شوند بعداً قابل بررسی باشند.

Logging فعلی:

- Structured logging با `slog`
- خروجی همزمان Terminal و فایل
- فایل پیش‌فرض: `logs/knowme.log`
- Rotation با Lumberjack
- تنظیم Level
- Max Size
- Max Backups
- Max Age
- Compression
- ثبت خطاهای Startup، Database، Telegram update handling و Panic در مرز Process

نباید Token تلگرام، Password دیتابیس یا Secretها در Log ثبت شوند.

## 9. Configuration

متغیرهای مهم:

```env
TELEGRAM_BOT_TOKEN=

APP_ENV=development

DB_PROVIDER=sqlite
SQLITE_PATH=data/knowme.db

SQLSERVER_DSN=

TELEGRAM_POLL_TIMEOUT_SECONDS=30

LOG_LEVEL=info
LOG_FILE=logs/knowme.log
LOG_MAX_SIZE_MB=20
LOG_MAX_BACKUPS=10
LOG_MAX_AGE_DAYS=14
LOG_COMPRESS=true
```

فایل `.env` نباید Commit شود.

## 10. مشکلاتی که در توسعه حل شدند

در مسیر راه‌اندازی چند مشکل عملی دیده و اصلاح شد:

- خطای `TELEGRAM_BOT_TOKEN is required` با اصلاح Configuration/Environment
- خطای SQL Server TLS handshake با تنظیم مناسب DSN و Certificate/Encryption برای محیط موردنظر
- خطای Load User Profile ناشی از ناسازگاری نام Column مربوط به Mobile/Prompt state
- انتقال Persistence از Raw SQL به GORM
- اضافه‌شدن SQLite برای توسعه بدون وابستگی به SQL Server
- اضافه‌شدن Logging ماندگار
- اصلاح CI برای Go/GORM/SQLite
- CI نهایی روی Commit مربوط به Refactor با موفقیت Pass شد.

## 11. Git و Repository

Repository:

`jahangard/KnowMe`

Repository به حالت **Private** تغییر داده شد.

Branch فعلی توسعه در این مرحله: `main`.

فایل‌های Local Database، Log و Environment در `.gitignore` قرار گرفته‌اند.

## 12. تصمیم‌های محصول برای ادامه

اولویت بعدی پیاده‌سازی Category Browser است، نه اضافه‌کردن تعداد زیادی تست بدون ساختار.

جریان پیشنهادی:

```text
Home
  └── تست‌ها
       ├── 🧭 نقشه راه من
       │     └── انتخاب خودکار تست بعدی
       ├── ❤️ عشق و رابطه
       │     └── لیست تست‌ها
       ├── 🎯 چالشی و باحال
       │     └── لیست تست‌ها
       ├── 🧠 خودشناسی عمیق‌تر
       │     └── لیست تست‌ها
       └── 🔞 ۱۸+
             └── Age Gate + لیست تست‌ها
```

Categoryها باید Data-driven باشند؛ یعنی عنوان، ترتیب، فعال/غیرفعال بودن و تست‌های هر دسته از دیتابیس خوانده شوند، نه اینکه منطق دسته‌ها در Telegram Handler هاردکد شود.

«نقشه راه من» می‌تواند در UI کنار Categoryها دیده شود ولی بهتر است یک Virtual/System Category باشد و لزوماً به رکورد عادی TestCategory وابسته نباشد.

## 13. جهت طراحی تست‌های آینده

تست‌ها باید:

- برای کاربر سرگرم‌کننده باشند، نه شبیه فرم اداری
- نتیجه واضح و جذاب داشته باشند
- تا جای ممکن برای زن و مرد قابل استفاده باشند مگر ماهیت تست خلاف آن را ایجاب کند
- نتیجه‌ای تولید کنند که ارزش Share کردن داشته باشد
- پاسخ‌ها فقط مصرف یک تست نباشند و در صورت امکان به UserTraitهای قابل استفاده در Recommendation Engine تبدیل شوند
- کوتاه شروع شوند و Completion Rate بالا داشته باشند
- از سؤال‌های تکراری و قابل حدس پرهیز کنند

## 14. تعریف Done برای مرحله بعد

مرحله Category Browser زمانی Done محسوب می‌شود که:

1. Home دکمه «تست‌ها» داشته باشد.
2. صفحه تست‌ها پنج مسیر توافق‌شده را نشان دهد.
3. Categoryهای واقعی از دیتابیس Load شوند.
4. انتخاب Category، تست‌های فعال همان Category را نشان دهد.
5. انتخاب Test، Test Engine فعلی را Start کند.
6. «نقشه راه من» همچنان Recommendation خودکار داشته باشد.
7. دسته ۱۸+ Age Gate داشته باشد.
8. Back/Home navigation بدون تولید پیام‌های اضافه کار کند.
9. SQLite و SQL Server هر دو Build/Run شوند.
10. CI سبز باقی بماند.

---

این سند باید با هر تصمیم مهم محصول، تغییر معماری یا تغییر جریان کاربر به‌روزرسانی شود تا گفتگوهای پراکنده به مرجع اجرایی واحد تبدیل شوند.
