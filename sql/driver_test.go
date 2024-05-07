package sql_test

import (
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type baseDriverMock struct {
	mock.Mock
}

func (m *baseDriverMock) Open(dsn string) (driver.Conn, error) {
	args := m.Called(dsn)

	if conn, ok := args.Get(0).(driver.Conn); ok {
		return conn, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (m *baseDriverMock) OpenConnector(dsn string) (driver.Connector, error) {
	args := m.Called(dsn)

	if connector, ok := args.Get(0).(driver.Connector); ok {
		return connector, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func TestNewDriver(t *testing.T) {
	t.Parallel()

	driverMock := new(baseDriverMock)

	config := sql.NewConfiguration(sql.SqliteSystem)

	d := sql.NewDriver(driverMock, config)

	assert.IsType(t, &sql.Driver{}, d)
	assert.Equal(t, driverMock, d.Base())
	assert.Equal(t, config, d.Configuration())

	driverMock.AssertExpectations(t)
}

func TestDriverOpenWithError(t *testing.T) {
	t.Parallel()

	driverMock := new(baseDriverMock)
	driverMock.On("Open", "test dsn").Return(nil, fmt.Errorf("test error"))

	config := sql.NewConfiguration(sql.SqliteSystem)

	d := sql.NewDriver(driverMock, config)

	_, err := d.Open("test dsn")
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	driverMock.AssertExpectations(t)
}

func TestDriverOpenConnector(t *testing.T) {
	t.Parallel()

	driverMock := new(baseDriverMock)
	driverMock.On("OpenConnector", "test dsn").Return(nil, nil)

	config := sql.NewConfiguration(sql.SqliteSystem)

	d := sql.NewDriver(driverMock, config)

	connector, err := d.OpenConnector("test dsn")
	assert.NoError(t, err)

	assert.IsType(t, &sql.Connector{}, connector)

	driverMock.AssertExpectations(t)
}

func TestDriverOpenConnectorWithError(t *testing.T) {
	t.Parallel()

	driverMock := new(baseDriverMock)
	driverMock.On("OpenConnector", "test dsn").Return(nil, fmt.Errorf("test error"))

	config := sql.NewConfiguration(sql.SqliteSystem)

	d := sql.NewDriver(driverMock, config)

	_, err := d.OpenConnector("test dsn")
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	driverMock.AssertExpectations(t)
}
