package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Provider string

const (
	ProviderSQLServer Provider = "sqlserver"
	ProviderSQLite    Provider = "sqlite"
)

type Handle struct {
	*sql.DB
	Provider Provider
}

func ParseProvider(raw string) (Provider, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "", string(ProviderSQLServer):
		return ProviderSQLServer, nil
	case string(ProviderSQLite):
		return ProviderSQLite, nil
	default:
		return "", fmt.Errorf("unsupported database provider %q", raw)
	}
}

func Open(ctx context.Context, provider Provider, sqlServerDSN, sqlitePath string) (*Handle, error) {
	switch provider {
	case ProviderSQLServer:
		db, err := OpenSQLServer(ctx, sqlServerDSN)
		if err != nil {
			return nil, err
		}
		return &Handle{DB: db, Provider: provider}, nil

	case ProviderSQLite:
		if strings.TrimSpace(sqlitePath) == "" {
			sqlitePath = "data/knowme.db"
		}
		if sqlitePath != ":memory:" {
			dir := filepath.Dir(sqlitePath)
			if dir != "." && dir != "" {
				if err := os.MkdirAll(dir, 0o755); err != nil {
					return nil, fmt.Errorf("create sqlite directory: %w", err)
				}
			}
		}

		db, err := sql.Open("sqlite", sqlitePath)
		if err != nil {
			return nil, fmt.Errorf("open sqlite: %w", err)
		}

		// A single connection keeps connection-level PRAGMAs deterministic and
		// is plenty for the local development workload.
		db.SetMaxOpenConns(1)
		db.SetMaxIdleConns(1)
		db.SetConnMaxLifetime(0)

		pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()

		if err := db.PingContext(pingCtx); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("ping sqlite: %w", err)
		}

		for _, pragma := range []string{
			"PRAGMA foreign_keys = ON",
			"PRAGMA journal_mode = WAL",
			"PRAGMA busy_timeout = 5000",
		} {
			if _, err := db.ExecContext(pingCtx, pragma); err != nil {
				_ = db.Close()
				return nil, fmt.Errorf("configure sqlite (%s): %w", pragma, err)
			}
		}

		return &Handle{DB: db, Provider: provider}, nil

	default:
		return nil, fmt.Errorf("unsupported database provider %q", provider)
	}
}

func (h *Handle) Table(name string) string {
	if h.Provider == ProviderSQLServer {
		return "dbo." + name
	}
	return name
}

func (h *Handle) NowExpr() string {
	if h.Provider == ProviderSQLServer {
		return "SYSUTCDATETIME()"
	}
	return "CURRENT_TIMESTAMP"
}

func (h *Handle) IsSQLite() bool {
	return h.Provider == ProviderSQLite
}
