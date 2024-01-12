package factory

import (
	"github.com/ankorstore/yokai/orm"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestOrmFactory struct{}

func NewTestOrmFactory() orm.OrmFactory {
	return &TestOrmFactory{}
}

func (f *TestOrmFactory) Create(options ...orm.OrmOption) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		DryRun: true,
	})
}
