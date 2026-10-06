package database

import (
	"context"
	"database/sql"
	"fmt"
	"io/fs"
	"sort"
	"strings"

	migrationfiles "github.com/jahangard/KnowMe/db/migrations"
)

func RunMigrations(ctx context.Context, db *Handle) error {
	if err := ensureMigrationTable(ctx, db); err != nil {
		return err
	}

	if db.Provider == ProviderSQLServer {
		if err := baselineLegacySQLServer(ctx, db); err != nil {
			return err
		}
	}

	pattern := "*.sql"
	if db.Provider == ProviderSQLite {
		pattern = "sqlite/*.sql"
	}

	files, err := fs.Glob(migrationfiles.Files, pattern)
	if err != nil {
		return fmt.Errorf("list migrations: %w", err)
	}
	sort.Strings(files)

	for _, name := range files {
		applied, err := migrationApplied(ctx, db, name)
		if err != nil {
			return err
		}
		if applied {
			continue
		}

		body, err := migrationfiles.Files.ReadFile(name)
		if err != nil {
			return fmt.Errorf("read migration %s: %w", name, err)
		}

		if _, err := db.ExecContext(ctx, string(body)); err != nil {
			return fmt.Errorf("apply migration %s: %w", name, err)
		}

		if err := recordMigration(ctx, db, name); err != nil {
			return err
		}
	}

	return nil
}

func ensureMigrationTable(ctx context.Context, db *Handle) error {
	var query string
	if db.Provider == ProviderSQLServer {
		query = `IF OBJECT_ID(N'dbo.SchemaMigrations', N'U') IS NULL
BEGIN
    CREATE TABLE dbo.SchemaMigrations (
        MigrationName NVARCHAR(255) NOT NULL CONSTRAINT PK_SchemaMigrations PRIMARY KEY,
        AppliedAt DATETIME2 NOT NULL CONSTRAINT DF_SchemaMigrations_AppliedAt DEFAULT SYSUTCDATETIME()
    );
END;`
	} else {
		query = `CREATE TABLE IF NOT EXISTS SchemaMigrations (
    MigrationName TEXT PRIMARY KEY,
    AppliedAt TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);`
	}

	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("ensure schema migrations table: %w", err)
	}
	return nil
}

func migrationApplied(ctx context.Context, db *Handle, name string) (bool, error) {
	query := fmt.Sprintf("SELECT COUNT(1) FROM %s WHERE MigrationName=@MigrationName", db.Table("SchemaMigrations"))
	var count int
	if err := db.QueryRowContext(ctx, query, sql.Named("MigrationName", name)).Scan(&count); err != nil {
		return false, fmt.Errorf("check migration %s: %w", name, err)
	}
	return count > 0, nil
}

func recordMigration(ctx context.Context, db *Handle, name string) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (MigrationName, AppliedAt) VALUES (@MigrationName, %s)",
		db.Table("SchemaMigrations"),
		db.NowExpr(),
	)
	if _, err := db.ExecContext(ctx, query, sql.Named("MigrationName", name)); err != nil {
		return fmt.Errorf("record migration %s: %w", name, err)
	}
	return nil
}

func baselineLegacySQLServer(ctx context.Context, db *Handle) error {
	if db.Provider != ProviderSQLServer {
		return nil
	}

	var usersExists int
	if err := db.QueryRowContext(ctx,
		"SELECT CASE WHEN OBJECT_ID(N'dbo.Users', N'U') IS NULL THEN 0 ELSE 1 END",
	).Scan(&usersExists); err != nil {
		return fmt.Errorf("detect legacy schema: %w", err)
	}
	if usersExists == 0 {
		return nil
	}

	marks := []string{"001_initial.sql"}

	var loveStyleExists int
	if err := db.QueryRowContext(ctx,
		"SELECT COUNT(1) FROM dbo.Tests WHERE Code=N'love_style_v1'",
	).Scan(&loveStyleExists); err == nil && loveStyleExists > 0 {
		marks = append(marks, "002_seed_love_style.sql")
	}

	var pendingField int
	if err := db.QueryRowContext(ctx,
		"SELECT CASE WHEN COL_LENGTH('dbo.UserProfiles','PendingField') IS NULL THEN 0 ELSE 1 END",
	).Scan(&pendingField); err == nil && pendingField == 1 {
		marks = append(marks, "003_progressive_profile.sql")
	}

	var promptFlags int
	if err := db.QueryRowContext(ctx,
		"SELECT CASE WHEN COL_LENGTH('dbo.UserProfiles','NamePrompted') IS NOT NULL AND COL_LENGTH('dbo.UserProfiles','MobilePrompted') IS NOT NULL THEN 1 ELSE 0 END",
	).Scan(&promptFlags); err == nil && promptFlags == 1 {
		marks = append(marks, "004_profile_prompt_flags.sql")
	}

	for _, name := range marks {
		applied, err := migrationApplied(ctx, db, name)
		if err != nil {
			return err
		}
		if applied {
			continue
		}
		if err := recordMigration(ctx, db, name); err != nil {
			return err
		}
	}

	return nil
}

func QuoteIdentifier(provider Provider, name string) string {
	if provider == ProviderSQLServer {
		return "[" + strings.ReplaceAll(name, "]", "]]") + "]"
	}
	return `"` + strings.ReplaceAll(name, `"`, `""`) + `"`
}
