package fxsql

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/log"
	"github.com/pressly/goose/v3"
)

// Migrator is a database migrator based on [Goose].
//
// [Fx]: https://github.com/pressly/goose
type Migrator struct {
	db     *sql.DB
	logger *log.Logger
}

// NewMigrator returns a new Migrator instance.
func NewMigrator(db *sql.DB, logger *log.Logger) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
	}
}

// Run executes a database migration command.
func (m *Migrator) Run(ctx context.Context, dialect string, dir string, command string, args ...string) error {
	m.logger.Info().
		Str("dir", dir).
		Str("command", command).
		Strs("args", args).
		Msg("starting database migration")

	// set dialect
	err := goose.SetDialect(dialect)
	if err != nil {
		m.logger.Error().Err(err).Msg("database dialect error")

		return err
	}

	// apply migration
	err = goose.RunContext(ctx, command, m.db, dir, args...)
	if err != nil {
		m.logger.Error().Err(err).Msg("database migration error")

		return err
	}

	m.logger.Info().Msg("database migration success")

	return nil
}
