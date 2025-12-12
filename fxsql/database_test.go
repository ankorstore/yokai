package fxsql_test

import (
	"database/sql"
	"testing"

	"github.com/ankorstore/yokai/fxsql"
	"github.com/stretchr/testify/assert"
)

func TestPrimaryDatabaseName(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "primary", fxsql.PrimaryDatabaseName)
}

func TestNewDatabase(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	database := fxsql.NewDatabase("test-db", db)

	assert.NotNil(t, database)
}

func TestDatabase_Name(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	database := fxsql.NewDatabase("test-db", db)

	assert.Equal(t, "test-db", database.Name())
}

func TestDatabase_DB(t *testing.T) {
	t.Parallel()

	db := &sql.DB{}
	database := fxsql.NewDatabase("test-db", db)

	assert.Same(t, db, database.DB())
}
