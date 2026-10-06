package database

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	sqlite "github.com/libtnb/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Provider string

const (
	ProviderSQLServer Provider = "sqlserver"
	ProviderSQLite    Provider = "sqlite"
)

type Handle struct {
	Gorm     *gorm.DB
	SQL      *sql.DB
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
	var dialector gorm.Dialector

	switch provider {
	case ProviderSQLServer:
		if strings.TrimSpace(sqlServerDSN) == "" {
			return nil, fmt.Errorf("SQLSERVER_DSN is required when DB_PROVIDER=sqlserver")
		}
		dialector = sqlserver.Open(sqlServerDSN)

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
			sqlitePath += "?_txlock=immediate&_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)"
		}
		dialector = sqlite.Open(sqlitePath)

	default:
		return nil, fmt.Errorf("unsupported database provider %q", provider)
	}

	orm, err := gorm.Open(dialector, &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   gormlogger.Default.LogMode(gormlogger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("open %s with gorm: %w", provider, err)
	}

	sqlDB, err := orm.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql db from gorm: %w", err)
	}

	if provider == ProviderSQLite {
		sqlDB.SetMaxOpenConns(1)
		sqlDB.SetMaxIdleConns(1)
		sqlDB.SetConnMaxLifetime(0)
	} else {
		sqlDB.SetMaxOpenConns(10)
		sqlDB.SetMaxIdleConns(5)
		sqlDB.SetConnMaxLifetime(30 * time.Minute)
		sqlDB.SetConnMaxIdleTime(5 * time.Minute)
	}

	pingCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		_ = sqlDB.Close()
		return nil, fmt.Errorf("ping %s: %w", provider, err)
	}

	return &Handle{
		Gorm:     orm,
		SQL:      sqlDB,
		Provider: provider,
	}, nil
}

func (h *Handle) Close() error {
	if h == nil || h.SQL == nil {
		return nil
	}
	return h.SQL.Close()
}
