package sql_test

import (
	basesql "database/sql"
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/stretchr/testify/assert"
)

func TestNewDefaultDriverRegistry(t *testing.T) {
	t.Parallel()

	registry := sql.NewDefaultDriverRegistry()

	assert.IsType(t, &sql.DefaultDriverRegistry{}, registry)
	assert.Implements(t, (*sql.DriverRegistry)(nil), registry)
}

func TestDriverRegistryLifecycle(t *testing.T) {
	t.Parallel()

	registry := sql.NewDefaultDriverRegistry()

	// empty registry assertions
	assert.False(t, registry.Has("sqlite"))

	_, err := registry.Get("sqlite")
	assert.Error(t, err)
	assert.Equal(t, "cannot find driver sqlite in driver registry", err.Error())

	// create and add a driver
	driver, err := sql.NewDefaultDriverFactory().Create(sql.SqliteSystem)
	assert.NoError(t, err)

	err = registry.Add("sqlite", driver)
	assert.NoError(t, err)

	// populated registry assertions
	assert.True(t, registry.Has("sqlite"))

	fetchedDriver, err := registry.Get("sqlite")
	assert.NoError(t, err)
	assert.Equal(t, driver, fetchedDriver)

	// add again the same driver
	err = registry.Add("sqlite", driver)
	assert.NoError(t, err)
}

func TestDriverRegistryErrorWithAlreadyRegisteredDriverInBaseRegistry(t *testing.T) {
	t.Parallel()

	registry := sql.NewDefaultDriverRegistry()

	// create a driver
	driver, err := sql.NewDefaultDriverFactory().Create(sql.SqliteSystem)
	assert.NoError(t, err)

	// driver base registration
	basesql.Register("test-sqlite", driver)

	// adding the driver again should panic and be recovered as an error
	err = registry.Add("test-sqlite", driver)
	assert.Error(t, err)
	assert.Equal(t, "cannot add driver test-sqlite: sql: Register called twice for driver test-sqlite", err.Error())
}
