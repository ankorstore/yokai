package fxsql

import (
	"context"
	"database/sql"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	yokaisql "github.com/ankorstore/yokai/sql"
	yokaisqllog "github.com/ankorstore/yokai/sql/hook/log"
	yokaisqltrace "github.com/ankorstore/yokai/sql/hook/trace"
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

// NewFxSQLDatabase returns a sql.DB instance.
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

// FxSQLMigratorParam allows injection of the required dependencies in [NewFxSQLMigrator].
type FxSQLMigratorParam struct {
	fx.In
	Db     *sql.DB
	Logger *log.Logger
	Config *config.Config
}

// NewFxSQLMigrator returns a Migrator instance.
func NewFxSQLMigrator(p FxSQLMigratorParam) *Migrator {
	return NewMigrator(p.Db, p.Logger, p.Config)
}

// RunFxSQLMigration runs database migrations.
func RunFxSQLMigration(command string, args ...string) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, migrator *Migrator, config *config.Config) error {
			return migrator.Migrate(
				ctx,
				config.GetString("modules.database.driver"),
				config.GetString("modules.database.migrations"),
				command,
				args...,
			)

		},
	)
}

// RunFxSQLMigrationAndShutdown runs database migrations and shutdown.
func RunFxSQLMigrationAndShutdown(command string, args ...string) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, migrator *Migrator, config *config.Config, shutdown fx.Shutdowner) error {
			defer shutdown.Shutdown()

			return migrator.Migrate(
				ctx,
				config.GetString("modules.database.driver"),
				config.GetString("modules.database.migrations"),
				command,
				args...,
			)
		},
	)
}
