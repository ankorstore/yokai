package fxsql

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/pressly/goose/v3"
)

type Migrator struct {
	db     *sql.DB
	logger *log.Logger
	config *config.Config
}

func NewMigrator(db *sql.DB, logger *log.Logger, config *config.Config) *Migrator {
	return &Migrator{
		db:     db,
		logger: logger,
		config: config,
	}
}

func (m *Migrator) Migrate(ctx context.Context, dialect string, dir string, command string, args ...string) error {
	// logger
	logger := m.logger.With().Str("command", command).Str("dir", dir).Logger()
	logger.Info().Msg("starting database migration")

	// set dialect
	err := goose.SetDialect(m.config.GetString("modules.database.driver"))
	if err != nil {
		logger.Error().Err(err).Msg("database migration dialect error")

		return err
	}

	// apply migration
	err = goose.RunContext(ctx, command, m.db, m.config.GetString("modules.database.migrations"), args...)
	if err != nil {
		logger.Error().Err(err).Msg("database migration error")

		return err
	}

	logger.Info().Msg("database migration success")

	return nil
}
