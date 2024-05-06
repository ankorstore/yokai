package sql

import (
	"fmt"

	"github.com/go-sql-driver/mysql"
	"github.com/lib/pq"
	"github.com/mattn/go-sqlite3"
)

type DriverFactory interface {
	Create(system System, hooks ...Hook) (*Driver, error)
}

type DefaultDriverFactory struct{}

func NewDefaultDriverFactory() *DefaultDriverFactory {
	return &DefaultDriverFactory{}
}

func (f *DefaultDriverFactory) Create(system System, hooks ...Hook) (*Driver, error) {
	switch system {
	case SqliteSystem:
		return NewDriver(&sqlite3.SQLiteDriver{}, NewConfiguration(system, hooks...)), nil
	case MysqlSystem:
		return NewDriver(&mysql.MySQLDriver{}, NewConfiguration(system, hooks...)), nil
	case PostgresSystem:
		return NewDriver(&pq.Driver{}, NewConfiguration(system, hooks...)), nil
	case UnknownSystem:
		return nil, fmt.Errorf("cannot create database driver for unknown system")
	default:
		return nil, fmt.Errorf("cannot create database driver for system %s", system.String())
	}
}
