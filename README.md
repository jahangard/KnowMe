# KnowMe

Telegram-based adaptive quiz engine with a data-driven test catalog.

## Current stack

- Go
- Telegram Bot via Long Polling
- GORM
- SQLite for local development
- SQL Server for production
- Test Catalog
- Test Engine
- Scoring Engine
- Roadmap Engine
- Progressive Profile
- Event Tracking
- Structured logging

## Current product flow

```text
Home
  ├── 🧩 موضوعات تست‌ها
  │    ├── 🧭 نقشه راه من
  │    ├── ❤️ عشق و رابطه
  │    ├── 🎯 چالشی و باحال
  │    ├── 🧠 خودشناسی عمیق‌تر
  │    └── 🔞 ۱۸+
  ├── 🧭 نقشه راه من
  ├── 👤 پروفایل من
  └── ✨ درباره KnowMe
```

The 18+ category is age-gated.

## Initial tests

The database bootstraps these tests when missing:

1. 💞 سبک عشق‌ورزی
2. ✨ تیپ جذابیت
3. 💘 سناریوهای قرار
4. 🛡 مرزهای شخصی
5. 🪞 حقیقت تلخ

Each initial test currently contains 8 single-choice questions.

## Data-driven test design

For standard `single_choice + dominant trait` tests, adding a new test requires database content only.

The runtime reads:

- category
- test metadata
- questions
- options
- TraitKey
- score
- result label
- result title
- result subtitle
- result description

from the database.

The test-specific result presentation is not hardcoded in Telegram UI.

See:

- `docs/ADDING-TESTS.md`
- `docs/architecture.md`

## Run locally with SQLite

1. Copy `.env.example` to `.env`.
2. Set `TELEGRAM_BOT_TOKEN`.
3. Keep:

```env
APP_ENV=development
DB_PROVIDER=sqlite
SQLITE_PATH=data/knowme.db
```

4. Run:

```bash
go mod tidy
go run ./cmd/knowme
```

The app creates the SQLite database and runs GORM AutoMigrate.

Bootstrap seed inserts missing initial categories/tests only. Existing database content is not overwritten on restart.

## Run with SQL Server

Set:

```env
DB_PROVIDER=sqlserver
SQLSERVER_DSN=sqlserver://user:password@server:1433?database=KnowMe&encrypt=true&TrustServerCertificate=true
```

Then:

```bash
go run ./cmd/knowme
```

Use an appropriate TLS/certificate policy for production.

## Logs

Runtime logs are written to the terminal and by default to:

```text
logs/knowme.log
```

Logging includes rotation and structured runtime errors.

Do not commit:

- Telegram bot token
- database credentials
- `.env`
- local SQLite databases
- runtime logs

## Documentation

Project decisions and rules are documented here:

- **Architecture:** `docs/architecture.md`
- **Product & content principles:** `docs/PRODUCT-PRINCIPLES.md`
- **Adding tests using database only:** `docs/ADDING-TESTS.md`
- **Decision log:** `docs/DECISIONS.md`
- **Change log:** `docs/CHANGELOG.md`
- **Conversation/decision summary:** `docs/CONVERSATION-SUMMARY.md`

These documents should be updated whenever a product rule, data contract, architecture decision, or important runtime behavior changes.
