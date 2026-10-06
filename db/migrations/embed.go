package migrations

import "embed"

// Files contains SQL Server migrations in this directory and SQLite-specific
// migrations in the sqlite subdirectory.
//
//go:embed *.sql sqlite/*.sql
var Files embed.FS
