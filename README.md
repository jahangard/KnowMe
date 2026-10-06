# KnowMe

Telegram-based adaptive quiz engine.

## v0 architecture
- Go
- Telegram Bot (Long Polling)
- Remote SQL Server
- Test Engine
- Scoring Engine
- Roadmap Engine
- User Profile / Progressive Profiling
- Event Tracking

## Run locally

1. Copy `.env.example` values into your environment.
2. Set `TELEGRAM_BOT_TOKEN`.
3. Set `SQLSERVER_DSN`.
4. Run:

```bash
go mod tidy
go run ./cmd/knowme
```

No bot token or database password should be committed to Git.
