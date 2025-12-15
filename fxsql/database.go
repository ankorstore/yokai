package fxsql

import "database/sql"

const PrimaryDatabaseName = "primary"

type Database struct {
	name string
	db   *sql.DB
}

func NewDatabase(name string, db *sql.DB) *Database {
	return &Database{
		name: name,
		db:   db,
	}
}

func (d *Database) Name() string {
	return d.name
}

func (d *Database) DB() *sql.DB {
	return d.db
}
