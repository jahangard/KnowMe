# KnowMe

Telegram-based adaptive quiz engine.

## Architecture
- Go
- Telegram Bot (Long Polling)
- GORM
- SQLite for local development
- SQL Server for production
- Test Engine
- Scoring Engine
- Roadmap Engine
- User Profile / Progressive Profiling
- Event Tracking

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

The app creates the SQLite database, runs GORM AutoMigrate, and seeds the first quiz automatically.

## Run with SQL Server

Set:

```env
DB_PROVIDER=sqlserver
SQLSERVER_DSN=sqlserver://user:password@server:1433?database=KnowMe&encrypt=true&TrustServerCertificate=true
```

Then run the same command:

```bash
go run ./cmd/knowme
```

The existing SQL Server schema is preserved; GORM AutoMigrate only adds missing schema elements and does not drop unused columns.

## Logs

Runtime logs are written to the terminal and to:

```text
logs/knowme.log
```

No bot token, database password, local SQLite database, or runtime log should be committed to Git.
