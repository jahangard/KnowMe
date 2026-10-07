# KnowMe Change Log

## 2026-10-07 — Telegram Mini App prototype

### Product experience

- Added a standalone mobile-first Persian RTL Mini App prototype under `telegram-mini-app/`.
- Added home, topic, category, roadmap, profile, quiz, result, and age-gate screens in KnowMe's dark plum and mint visual direction.
- The eight «سبک عشق‌ورزی» questions, answer options, and result profiles mirror the current Go seed content.
- Added Telegram Web App lifecycle, back-button, and haptic integration for use inside Telegram.
- Other catalog tests are visible but remain unavailable in the standalone preview; answers and profiles are not persisted to the Go backend.

### Production follow-up

- Add a public HTTPS launch URL and Telegram Web App entry point in the bot.
- Add backend endpoints and signed `initData` validation before using the Mini App with real user data or age-restricted content.

## 2026-10-07 — Data-driven Quiz Platform

### Product

- دو مسیر اصلی Test Experience تثبیت شد:
  - انتخاب آزاد از Topic
  - «نقشه راه من»
- Category Browser به Telegram اضافه شد.
- Roadmap کنار Topicها نمایش داده می‌شود.
- Categoryهای اصلی:
  - عشق و رابطه
  - چالشی و باحال
  - خودشناسی عمیق‌تر
  - ۱۸+
- Age Gate برای Category ۱۸+ اضافه شد.
- UX بر Edit Message متمرکز شد تا Chat شلوغ نشود.

### Initial test catalog

پنج Test اولیه اضافه شد:

1. 💞 سبک عشق‌ورزی — `love_style_v1`
2. ✨ تیپ جذابیت — `attraction_style_v1`
3. 💘 سناریوهای قرار — `dating_scenarios_v1`
4. 🛡 مرزهای شخصی — `personal_boundaries_v1`
5. 🪞 حقیقت تلخ — `bitter_truth_v1`

هر Test فعلاً 8 سؤال دارد.

### Data-driven Test Engine

- Questionها از DB خوانده می‌شوند.
- Optionها از DB خوانده می‌شوند.
- `TraitKey` و `Score` از DB خوانده می‌شوند.
- Result text از DB خوانده می‌شود.
- Score breakdown label از DB خوانده می‌شود.
- Telegram UI دیگر Result presentation مخصوص Test را Hard-code نمی‌کند.

### New table

`TestResultProfiles` اضافه شد.

کاربرد:

- TraitKey
- Label
- Title
- Subtitle
- Description
- SortOrder
- IsActive

### Database as source of truth

Seed behavior اصلاح شد:

- فقط رکورد Missing ایجاد می‌شود.
- Existing Test content overwrite نمی‌شود.
- Restart برنامه تغییرات دستی DB را از بین نمی‌برد.

این تغییر برای آماده‌سازی Admin/Test Studio آینده انجام شد.

### Database abstraction

Persistence از Raw SQL به GORM منتقل شد.

Providerهای پشتیبانی‌شده:

- SQLite
- SQL Server

Development default:

`SQLite`

Production target:

`SQL Server`

### SQLite Development

- local file DB
- WAL
- foreign keys
- busy timeout
- no SQL Server dependency for local development

### SQL Server

- Remote SQL Server support
- TLS DSN support
- اتصال قبلی با TrustServerCertificate برای محیط فعلی رفع اشکال شد
- Production باید TLS policy مناسب داشته باشد

### Configuration

`.env` loading اضافه شد.

متغیرهای اصلی:

- TELEGRAM_BOT_TOKEN
- APP_ENV
- DB_PROVIDER
- SQLITE_PATH
- SQLSERVER_DSN
- TELEGRAM_POLL_TIMEOUT_SECONDS
- LOG_LEVEL
- LOG_FILE
- LOG_MAX_SIZE_MB
- LOG_MAX_BACKUPS
- LOG_MAX_AGE_DAYS
- LOG_COMPRESS

### Logging

Logging تقویت شد:

- structured `slog`
- terminal output
- persistent file
- rotation
- panic stack
- startup/shutdown
- DB errors
- Telegram update errors
- duration

### Profile

Progressive Profiling تثبیت شد:

- Gender
- Age
- optional Name
- optional Mobile

Mobile در Profile به شکل masked state نمایش داده می‌شود، نه مقدار واقعی.

### Runtime

- Go
- Telegram Long Polling
- no inbound public port
- Windows-friendly
- Linux/ARM friendly

### CI

GitHub Actions شامل:

- Go setup
- dependency resolution
- gofmt validation
- go test ./...

آخرین تغییرات Test Engine و Data-driven architecture Compile/Test موفق داشته‌اند.

### Documentation

اسناد اصلی:

- `docs/architecture.md`
- `docs/ADDING-TESTS.md`
- `docs/PRODUCT-PRINCIPLES.md`
- `docs/CHANGELOG.md`
- `docs/CONVERSATION-SUMMARY.md`

## Previous implementation milestones

### Initial MVP

- Telegram bot skeleton
- Users/Profile
- Test sessions
- Answer persistence
- Result persistence
- UserTraits
- UserEvents schema
- Roadmap v0
- first Love Style test

### Runtime fixes

رفع موارد:

- missing Telegram token due to env loading
- SQL Server TLS handshake configuration
- profile column mismatch around mobile prompt flags
- local SQLite support
- CI formatting/compile issues
