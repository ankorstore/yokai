package orm_test

import (
	"testing"

	"github.com/ankorstore/yokai/orm"
	"github.com/stretchr/testify/assert"
	base "gorm.io/gorm"
)

func TestWithDsn(t *testing.T) {
	t.Parallel()

	opt := orm.DefaultOrmOptions()
	orm.WithDsn("dsn")(&opt)

	assert.Equal(t, "dsn", opt.Dsn)
}

func TestWithDriver(t *testing.T) {
	t.Parallel()

	opt := orm.DefaultOrmOptions()
	orm.WithDriver(orm.Sqlite)(&opt)

	assert.Equal(t, orm.Sqlite, opt.Driver)
}

func TestWithConfig(t *testing.T) {
	t.Parallel()

	cfg := base.Config{}

	opt := orm.DefaultOrmOptions()
	orm.WithConfig(cfg)(&opt)

	assert.Equal(t, cfg, opt.Config)
}
