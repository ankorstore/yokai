package fxsql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	yokaisql "github.com/ankorstore/yokai/sql"
	yokaisqllog "github.com/ankorstore/yokai/sql/hook/log"
	yokaisqltrace "github.com/ankorstore/yokai/sql/hook/trace"
	"github.com/pressly/goose/v3"
	"go.uber.org/fx"
)

var once sync.Once

// ModuleName is the module name.
const ModuleName = "sql"

// FxSQLModule is the [Fx] SQL module.
//
// [Fx]: https://github.com/uber-go/fx
var FxSQLModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewFxSQLDatabasePool,
		NewFxSQLPrimaryDatabase,
		NewFxSQLMigrator,
		NewFxSQLSeeder,
	),
)

// FxSQLDatabasePoolParam allows injection of the required dependencies in [NewFxSQLDatabasePool].
type FxSQLDatabasePoolParam struct {
	fx.In
	LifeCycle fx.Lifecycle
	Config    *config.Config
	Logger    *log.Logger
	Hooks     []yokaisql.Hook `group:"sql-hooks"`
}

// NewFxSQLDatabasePool returns a DatabasePool instance.
func NewFxSQLDatabasePool(p FxSQLDatabasePoolParam) (*DatabasePool, error) {
	// database drivers hooks
	driverHooks := p.Hooks

	// database drivers trace hook
	if p.Config.GetBool("modules.sql.trace.enabled") {
		driverHooks = append(
			driverHooks,
			yokaisqltrace.NewTraceHook(
				yokaisqltrace.WithArguments(p.Config.GetBool("modules.sql.trace.arguments")),
				yokaisqltrace.WithExcludedOperations(
					yokaisql.FetchOperations(p.Config.GetStringSlice("modules.sql.trace.exclude"))...,
				),
			),
		)
	}

	// database drivers log hook
	if p.Config.GetBool("modules.sql.log.enabled") {
		driverHooks = append(
			driverHooks,
			yokaisqllog.NewLogHook(
				yokaisqllog.WithLevel(log.FetchLogLevel(p.Config.GetString("modules.sql.log.level"))),
				yokaisqllog.WithArguments(p.Config.GetBool("modules.sql.log.arguments")),
				yokaisqllog.WithExcludedOperations(
					yokaisql.FetchOperations(p.Config.GetStringSlice("modules.sql.log.exclude"))...,
				),
			),
		)
	}

	// primary database preparation
	primaryDriverName, err := yokaisql.RegisterNamed(p.Config.GetString("modules.sql.driver"), PrimaryDatabaseName, driverHooks...)
	if err != nil {
		return nil, err
	}

	primaryDB, err := sql.Open(primaryDriverName, p.Config.GetString("modules.sql.dsn"))
	if err != nil {
		return nil, err
	}

	primaryDatabase := NewDatabase(PrimaryDatabaseName, primaryDB)

	// auxiliaries databases preparation
	var auxiliaryDatabases []*Database

	for auxiliaryDatabaseName := range p.Config.GetStringMap("modules.sql.auxiliaries") {
		auxiliaryDriverName, err := yokaisql.RegisterNamed(
			p.Config.GetString(fmt.Sprintf("modules.sql.auxiliaries.%s.driver", auxiliaryDatabaseName)),
			auxiliaryDatabaseName,
			driverHooks...,
		)
		if err != nil {
			return nil, err
		}

		auxiliaryDB, err := sql.Open(auxiliaryDriverName, fmt.Sprintf("modules.sql.auxiliaries.%s.dsn", auxiliaryDatabaseName))
		if err != nil {
			return nil, err
		}

		auxiliaryDatabases = append(auxiliaryDatabases, NewDatabase(auxiliaryDatabaseName, auxiliaryDB))
	}

	// database pool
	databasePool := NewDatabasePool(primaryDatabase, auxiliaryDatabases...)

	// lifecycles
	p.LifeCycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			if !p.Config.IsTestEnv() {
				// close auxiliaries
				for _, auxiliaryDatabase := range databasePool.Auxiliaries() {
					err = auxiliaryDatabase.DB().Close()
					if err != nil {
						return err
					}
				}

				// close primary
				return databasePool.Primary().DB().Close()
			}

			return nil
		},
	})

	return databasePool, nil
}

// FxSQLPrimaryDatabaseParam allows injection of the required dependencies in [FxSQLPrimaryDatabase].
type FxSQLPrimaryDatabaseParam struct {
	fx.In
	Pool *DatabasePool
}

// NewFxSQLPrimaryDatabase returns the primary database.
func NewFxSQLPrimaryDatabase(p FxSQLPrimaryDatabaseParam) *sql.DB {
	return p.Pool.Primary().DB()
}

// FxSQLMigratorParam allows injection of the required dependencies in [NewFxSQLMigrator].
type FxSQLMigratorParam struct {
	fx.In
	Db     *sql.DB
	Config *config.Config
	Logger *log.Logger
}

// NewFxSQLMigrator returns a Migrator instance.
func NewFxSQLMigrator(p FxSQLMigratorParam) *Migrator {
	// set once migrator the logger
	once.Do(func() {
		goose.SetLogger(NewMigratorLogger(p.Logger, p.Config.GetBool("modules.sql.migrations.stdout")))
	})

	// migrator
	return NewMigrator(p.Db, p.Logger)
}

// FxSQLSeederParam allows injection of the required dependencies in [NewFxSQLSeeder].
type FxSQLSeederParam struct {
	fx.In
	Db     *sql.DB
	Logger *log.Logger
	Seeds  []Seed `group:"sql-seeds"`
}

// NewFxSQLSeeder returns a Seeder instance.
func NewFxSQLSeeder(p FxSQLSeederParam) *Seeder {
	return NewSeeder(p.Db, p.Logger, p.Seeds...)
}

// RunFxSQLMigration runs database migrations with a context.
func RunFxSQLMigration(command string, args ...string) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, migrator *Migrator, config *config.Config) error {
			return migrator.Run(
				ctx,
				config.GetString("modules.sql.driver"),
				config.GetString("modules.sql.migrations.path"),
				command,
				args...,
			)
		},
	)
}

// RunFxSQLMigrationAndShutdown runs database migrations with a context and shutdown.
func RunFxSQLMigrationAndShutdown(command string, args ...string) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, migrator *Migrator, config *config.Config, shutdown fx.Shutdowner) error {
			//nolint:errcheck
			defer shutdown.Shutdown()

			return migrator.Run(
				ctx,
				config.GetString("modules.sql.driver"),
				config.GetString("modules.sql.migrations.path"),
				command,
				args...,
			)
		},
	)
}

// RunFxSQLSeeds runs database seeds with a context.
func RunFxSQLSeeds(names ...string) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, seeder *Seeder) error {
			return seeder.Run(ctx, names...)
		},
	)
}

// RunFxSQLSeedsAndShutdown runs database seeds with a context and shutdown.
func RunFxSQLSeedsAndShutdown(names ...string) fx.Option {
	return fx.Invoke(
		func(ctx context.Context, seeder *Seeder, shutdown fx.Shutdowner) error {
			//nolint:errcheck
			defer shutdown.Shutdown()

			return seeder.Run(ctx, names...)
		},
	)
}
