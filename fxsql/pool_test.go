package fxsql_test

import (
	"database/sql"
	"testing"

	"github.com/ankorstore/yokai/fxsql"
	"github.com/stretchr/testify/assert"
)

func TestNewDatabasePool(t *testing.T) {
	t.Parallel()

	primary := fxsql.NewDatabase("primary", &sql.DB{})
	pool := fxsql.NewDatabasePool(primary)

	assert.NotNil(t, pool)
}

func TestNewDatabasePoolWithAuxiliaries(t *testing.T) {
	t.Parallel()

	primary := fxsql.NewDatabase("primary", &sql.DB{})
	aux1 := fxsql.NewDatabase("aux1", &sql.DB{})
	aux2 := fxsql.NewDatabase("aux2", &sql.DB{})

	pool := fxsql.NewDatabasePool(primary, aux1, aux2)

	assert.NotNil(t, pool)
	assert.Len(t, pool.Auxiliaries(), 2)
}

func TestDatabasePool_Primary(t *testing.T) {
	t.Parallel()

	primary := fxsql.NewDatabase("primary", &sql.DB{})
	pool := fxsql.NewDatabasePool(primary)

	assert.Same(t, primary, pool.Primary())
}

func TestDatabasePool_Auxiliary(t *testing.T) {
	t.Parallel()

	primary := fxsql.NewDatabase("primary", &sql.DB{})
	aux1 := fxsql.NewDatabase("aux1", &sql.DB{})
	aux2 := fxsql.NewDatabase("aux2", &sql.DB{})

	pool := fxsql.NewDatabasePool(primary, aux1, aux2)

	result, err := pool.Auxiliary("aux1")
	assert.NoError(t, err)
	assert.Same(t, aux1, result)

	result, err = pool.Auxiliary("aux2")
	assert.NoError(t, err)
	assert.Same(t, aux2, result)
}

func TestDatabasePool_AuxiliaryNotFound(t *testing.T) {
	t.Parallel()

	primary := fxsql.NewDatabase("primary", &sql.DB{})
	pool := fxsql.NewDatabasePool(primary)

	result, err := pool.Auxiliary("nonexistent")
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, "database with name nonexistent was not found", err.Error())
}

func TestDatabasePool_Auxiliaries(t *testing.T) {
	t.Parallel()

	primary := fxsql.NewDatabase("primary", &sql.DB{})
	aux1 := fxsql.NewDatabase("aux1", &sql.DB{})
	aux2 := fxsql.NewDatabase("aux2", &sql.DB{})

	pool := fxsql.NewDatabasePool(primary, aux1, aux2)

	auxiliaries := pool.Auxiliaries()
	assert.Len(t, auxiliaries, 2)
	assert.Same(t, aux1, auxiliaries["aux1"])
	assert.Same(t, aux2, auxiliaries["aux2"])
}

func TestDatabasePool_EmptyAuxiliaries(t *testing.T) {
	t.Parallel()

	primary := fxsql.NewDatabase("primary", &sql.DB{})
	pool := fxsql.NewDatabasePool(primary)

	auxiliaries := pool.Auxiliaries()
	assert.NotNil(t, auxiliaries)
	assert.Empty(t, auxiliaries)
}
