package sql

import "strings"

// System is an enum for the supported database systems.
type System string

const (
	UnknownSystem  System = "unknown"
	SqliteSystem   System = "sqlite"
	MysqlSystem    System = "mysql"
	PostgresSystem System = "postgres"
)

// String returns a string representation of the System.
func (d System) String() string {
	return string(d)
}

// FetchSystem returns a System for a given name.
func FetchSystem(name string) System {
	//nolint:exhaustive
	switch d := System(strings.ToLower(name)); d {
	case SqliteSystem,
		MysqlSystem,
		PostgresSystem:
		return d
	default:
		return UnknownSystem
	}
}
