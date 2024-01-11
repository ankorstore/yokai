package orm_test

import (
	"testing"

	"github.com/ankorstore/yokai/orm"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func TestDefaultOrmFactory(t *testing.T) {
	t.Parallel()

	factory := orm.NewDefaultOrmFactory()

	assert.IsType(t, &orm.DefaultOrmFactory{}, factory)
	assert.Implements(t, (*orm.OrmFactory)(nil), factory)
}

func TestCreateSuccessWithSqliteDriver(t *testing.T) {
	t.Parallel()

	db, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.Sqlite),
	)
	assert.NoError(t, err)

	d, err := db.DB()
	assert.NoError(t, err)

	err = d.Ping()
	assert.NoError(t, err)

	err = d.Close()
	assert.NoError(t, err)
}

func TestCreateFailureWithMysqlDriverAndInvalidDsn(t *testing.T) {
	t.Parallel()

	_, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.Mysql),
		orm.WithDsn("0.0.0.0/test"),
		orm.WithConfig(gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}),
	)
	assert.Error(t, err)
}

func TestCreateFailureWithPgsqlDriverAndInvalidDsn(t *testing.T) {
	t.Parallel()

	_, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.Postgres),
		orm.WithDsn("0.0.0.0/test"),
		orm.WithConfig(gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}),
	)
	assert.Error(t, err)
}

func TestCreateFailureWithSqlServerDriverAndInvalidDsn(t *testing.T) {
	t.Parallel()

	_, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.SqlServer),
		orm.WithDsn("0.0.0.0/test"),
		orm.WithConfig(gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}),
	)
	assert.Error(t, err)
}

func TestCreateFailureWithUnknownDriver(t *testing.T) {
	t.Parallel()

	_, err := orm.NewDefaultOrmFactory().Create(
		orm.WithDriver(orm.Unknown),
		orm.WithConfig(gorm.Config{
			Logger: logger.Default.LogMode(logger.Silent),
		}),
	)

	assert.Error(t, err)
	assert.Equal(t, "unsupported driver unknown", err.Error())
}
