package orm

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

// OrmFactory is the interface for [gorm.DB] factories.
type OrmFactory interface {
	Create(options ...OrmOption) (*gorm.DB, error)
}

// DefaultOrmFactory is the default [OrmFactory] implementation.
type DefaultOrmFactory struct{}

// NewDefaultOrmFactory returns a [DefaultOrmFactory], implementing [OrmFactory].
func NewDefaultOrmFactory() OrmFactory {
	return &DefaultOrmFactory{}
}

// Create returns a new [gorm.DB], and accepts a list of [OrmOption].
// For example with MySQL driver:
//
//	var db, _ = orm.NewDefaultOrmFactory().Create(
//		orm.WithDriver(orm.Mysql),
//		orm.WithDsn("user:pass@tcp(127.0.0.1:3306)/dbname?parseTime=True"),
//	)
//
// or with SQLite driver:
//
//	var db, _ = orm.NewDefaultOrmFactory().Create(
//		orm.WithDriver(orm.Sqlite),
//		orm.WithDsn("file::memory:?cache=shared"),
//	)
func (f *DefaultOrmFactory) Create(options ...OrmOption) (*gorm.DB, error) {
	appliedOpts := DefaultOrmOptions()
	for _, applyOpt := range options {
		applyOpt(&appliedOpts)
	}

	db, err := f.createDatabase(appliedOpts.Driver.String(), appliedOpts.Dsn, appliedOpts.Config)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (f *DefaultOrmFactory) createDatabase(driver string, dsn string, config gorm.Config) (*gorm.DB, error) {
	var dial gorm.Dialector

	switch FetchDriver(driver) {
	case Sqlite:
		dial = sqlite.Open(dsn)
	case Mysql:
		dial = mysql.Open(dsn)
	case Postgres:
		dial = postgres.Open(dsn)
	case SqlServer:
		dial = sqlserver.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported driver %s", driver)
	}

	return gorm.Open(dial, &config)
}
