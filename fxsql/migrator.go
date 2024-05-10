package fxsql

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/log"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db     *sql.DB
	logger *log.Logger
}

func NewMigrator(db *sql.DB, logger *log.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

func (m *Migrator) Migrate(ctx context.Context, dialect string, dir string, command string, args ...string) error {
	// logger
	logger := m.logger.
		With().
		Str("dir", dir).
		Str("command", command).
		Strs("args", args).
		Logger()

	logger.Info().Msg("starting database migration")

	// set dialect
	err := goose.SetDialect(dialect)
	if err != nil {
		logger.Error().Err(err).Msg("database migration dialect error")

		return err
	}

	// apply migration
	err = goose.RunContext(ctx, command, m.db, dir, args...)
	if err != nil {
		logger.Error().Err(err).Msg("database migration error")

		return err
	}

	logger.Info().Msg("database migration success")

	return nil
}
