// Package db はデータベース接続とマイグレーション適用を提供する。
package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	migratepg "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // database/sql 用 "pgx" ドライバを登録

	"github.com/HossyWorlds/next-go-best/backend/migrations"
)

// Migrate は埋め込みSQLで最新までマイグレーションを適用する。
// pgx の database/sql ドライバを使い、SQLスタックを pgx に統一する。
func Migrate(databaseURL string) error {
	m, err := newMigrator(databaseURL)
	if err != nil {
		return err
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("db: migrate up: %w", err)
	}
	return nil
}

// MigrateDown は1ステップだけロールバックする（ロールバック手順の確認用）。
func MigrateDown(databaseURL string) error {
	m, err := newMigrator(databaseURL)
	if err != nil {
		return err
	}
	defer func() { _, _ = m.Close() }()

	if err := m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("db: migrate down: %w", err)
	}
	return nil
}

func newMigrator(databaseURL string) (*migrate.Migrate, error) {
	src, err := iofs.New(migrations.FS, ".")
	if err != nil {
		return nil, fmt.Errorf("db: load embedded migrations: %w", err)
	}

	sqlDB, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("db: open sql: %w", err)
	}

	driver, err := migratepg.WithInstance(sqlDB, &migratepg.Config{})
	if err != nil {
		return nil, fmt.Errorf("db: migrate driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", src, "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("db: new migrator: %w", err)
	}
	return m, nil
}
