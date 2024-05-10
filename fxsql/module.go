package fxsql

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	yokaisql "github.com/ankorstore/yokai/sql"
	yokaisqllog "github.com/ankorstore/yokai/sql/hook/log"
	yokaisqltrace "github.com/ankorstore/yokai/sql/hook/trace"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "sql"

// FxSQLModule is the [Fx] SQL module.
//
// [Fx]: https://github.com/uber-go/fx
var FxSQLModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewFxSQLDatabase,
	),
)

// FxSQLDatabaseParam allows injection of the required dependencies in [NewFxSQLDatabase].
type FxSQLDatabaseParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
	Logger    *log.Logger
	Hooks     []yokaisql.Hook `group:"sql-hooks"`
}

// NewFxSQLDatabase returns a *sql.DB instance.
func NewFxSQLDatabase(p FxSQLDatabaseParam) (*sql.DB, error) {
	// custom hooks
	driverHooks := p.Hooks

	// trace hook
	if p.Config.GetBool("modules.database.trace.enabled") {
		driverHooks = append(
			driverHooks,
			yokaisqltrace.NewTraceHook(
				yokaisqltrace.WithArguments(p.Config.GetBool("modules.database.trace.arguments")),
				yokaisqltrace.WithExcludedOperations(
					yokaisql.FetchOperations(p.Config.GetStringSlice("modules.database.trace.exclude"))...,
				),
			),
		)
	}

	// log hook
	if p.Config.GetBool("modules.database.log.enabled") {
		driverHooks = append(
			driverHooks,
			yokaisqllog.NewLogHook(
				yokaisqllog.WithLevel(log.FetchLogLevel(p.Config.GetString("modules.database.log.level"))),
				yokaisqllog.WithArguments(p.Config.GetBool("modules.database.log.arguments")),
				yokaisqllog.WithExcludedOperations(
					yokaisql.FetchOperations(p.Config.GetStringSlice("modules.database.log.exclude"))...,
				),
			),
		)
	}

	// driver registration
	driverName, err := yokaisql.Register(p.Config.GetString("modules.database.driver"), driverHooks...)
	if err != nil {
		return nil, err
	}

	// database preparation
	db, err := sql.Open(driverName, p.Config.GetString("modules.database.dsn"))
	if err != nil {
		return nil, err
	}

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			if yokaisql.FetchSystem(p.Config.GetString("modules.database.driver")) != yokaisql.SqliteSystem {
				return db.Close()
			}

			return nil
		},
	})

	return db, nil
}

// RunFxSQLDatabaseMigration runs database migrations.
func RunFxSQLDatabaseMigration(command string, shutdown bool) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, db *sql.DB, cfg *config.Config, lgr *log.Logger, sd fx.Shutdowner) error {
			// source dir
			dir := cfg.GetString("modules.database.migrations")

			logger := lgr.With().Str("command", command).Str("dir", dir).Logger()
			logger.Info().Msg("starting database migration")

			// set dialect
			err := goose.SetDialect(cfg.GetString("modules.database.driver"))
			if err != nil {
				logger.Error().Err(err).Msg("database migration dialect error")

				return err
			}

			// apply migration
			err = goose.RunContext(ctx, command, db, cfg.GetString("modules.database.migrations"))
			if err != nil {
				logger.Error().Err(err).Msg("database migration error")

				return err
			}

			logger.Info().Msg("database migration success")

			// shutdown
			if shutdown {
				return sd.Shutdown()
			}

			return nil
		},
	)
}
