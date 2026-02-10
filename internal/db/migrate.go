package db

import (
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// MigrateConfig holds migration config
type MigrateConfig struct {
	DBURL         string
	MigrationPath string
	Logger        *slog.Logger
}

// RunMigrations executes pending migrations
func RunMigrations(cfg MigrateConfig) error {
	m, err := migrate.New(cfg.MigrationPath, cfg.DBURL)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}
	defer m.Close()

	// set logger if provided
	if cfg.Logger != nil {
		m.Log = &migrateLogger{logger: cfg.Logger}
	}

	// get current version before migration
	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	if err == nil {
		cfg.Logger.Info("current schema version", "version", version, "dirty", dirty)
	}

	// run migration
	if err := m.Up(); err != nil {
		if err == migrate.ErrNoChange {
			cfg.Logger.Info("database schema up to date")
			return nil
		}
		return fmt.Errorf("migration failed %w", err)
	}

	// get version after migration
	newVersion, _, _ := m.Version()
	cfg.Logger.Info("migration applied successfully", "version", newVersion)

	return nil
}

// migrateLogger adapts migrate logs to slog
type migrateLogger struct {
	logger *slog.Logger
}

func (l *migrateLogger) Printf(format string, v ...any) {
	l.logger.Info(fmt.Sprintf(format, v...))
}

func (l *migrateLogger) Verbose() bool {
	return true
}
