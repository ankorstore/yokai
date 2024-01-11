package orm

import (
	"gorm.io/gorm"
)

// Options are options for the [OrmFactory] implementations.
type Options struct {
	Dsn    string
	Driver Driver
	Config gorm.Config
}

// DefaultOrmOptions are the default options used in the [DefaultOrmFactory].
func DefaultOrmOptions() Options {
	return Options{
		Dsn:    "",
		Driver: Unknown,
		Config: gorm.Config{},
	}
}

// OrmOption are functional options for the [OrmFactory] implementations.
type OrmOption func(o *Options)

// WithDsn is used to specify the database DSN to use.
func WithDsn(d string) OrmOption {
	return func(o *Options) {
		o.Dsn = d
	}
}

// WithDriver is used to specify the database driver to use.
func WithDriver(d Driver) OrmOption {
	return func(o *Options) {
		o.Driver = d
	}
}

// WithConfig is used to specify the [gorm.Config] to use.
func WithConfig(c gorm.Config) OrmOption {
	return func(o *Options) {
		o.Config = c
	}
}
