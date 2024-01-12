package fxorm

import (
	"context"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/orm"
	"github.com/ankorstore/yokai/orm/plugin"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// ModuleName is the module name.
const ModuleName = "orm"

// FxOrmModule is the [Fx] orm module.
//
// [Fx]: https://github.com/uber-go/fx
var FxOrmModule = fx.Module(
	ModuleName,
	fx.Provide(
		orm.NewDefaultOrmFactory,
		NewFxOrm,
	),
)

// FxOrmParam allows injection of the required dependencies in [NewFxOrm].
type FxOrmParam struct {
	fx.In
	LifeCycle      fx.Lifecycle
	Factory        orm.OrmFactory
	Config         *config.Config
	TracerProvider trace.TracerProvider
}

// RunFxOrmAutoMigrate performs auto migrations for a provided list of models.
func RunFxOrmAutoMigrate(models ...any) fx.Option {
	return fx.Invoke(func(logger *log.Logger, db *gorm.DB) error {
		logger.Info().Msg("starting ORM auto migration")

		err := db.AutoMigrate(models...)
		if err != nil {
			logger.Error().Err(err).Msg("error during ORM auto migration")

			return err
		}

		logger.Info().Msg("ORM auto migration success")

		return nil
	})
}

// NewFxOrm returns a [gorm.DB].
func NewFxOrm(p FxOrmParam) (*gorm.DB, error) {
	ormConfig := gorm.Config{
		Logger:                                   logger.Default.LogMode(logger.Silent),
		DryRun:                                   p.Config.GetBool("modules.orm.config.dry_run"),
		SkipDefaultTransaction:                   p.Config.GetBool("modules.orm.config.skip_default_transaction"),
		FullSaveAssociations:                     p.Config.GetBool("modules.orm.config.full_save_associations"),
		PrepareStmt:                              p.Config.GetBool("modules.orm.config.prepare_stmt"),
		DisableAutomaticPing:                     p.Config.GetBool("modules.orm.config.disable_automatic_ping"),
		DisableForeignKeyConstraintWhenMigrating: p.Config.GetBool("modules.orm.config.disable_foreign_key_constraint_when_migrating"),
		IgnoreRelationshipsWhenMigrating:         p.Config.GetBool("modules.orm.config.ignore_relationships_when_migrating"),
		DisableNestedTransaction:                 p.Config.GetBool("modules.orm.config.disable_nested_transaction"),
		AllowGlobalUpdate:                        p.Config.GetBool("modules.orm.config.allow_global_update"),
		QueryFields:                              p.Config.GetBool("modules.orm.config.query_fields"),
		TranslateError:                           p.Config.GetBool("modules.orm.config.translate_error"),
	}

	if p.Config.GetBool("modules.orm.log.enabled") {
		ormConfig.Logger = orm.NewCtxOrmLogger(
			orm.FetchLogLevel(p.Config.GetString("modules.orm.log.level")),
			p.Config.GetBool("modules.orm.log.values"),
		)
	}

	driver := orm.FetchDriver(p.Config.GetString("modules.orm.driver"))

	db, err := p.Factory.Create(
		orm.WithDsn(p.Config.GetString("modules.orm.dsn")),
		orm.WithDriver(driver),
		orm.WithConfig(ormConfig),
	)

	if err != nil {
		return nil, err
	}

	if p.Config.GetBool("modules.orm.trace.enabled") {
		err = db.Use(plugin.NewOrmTracerPlugin(p.TracerProvider, p.Config.GetBool("modules.orm.trace.values")))
		if err != nil {
			return nil, err
		}
	}

	p.LifeCycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if driver != orm.Sqlite {
				ormDb, err := db.DB()
				if err != nil {
					return err
				}

				return ormDb.Close()
			}

			return nil
		},
	})

	return db, nil
}
