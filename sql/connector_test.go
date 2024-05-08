package sql_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type baseConnectorMock struct {
	mock.Mock
}

func (m *baseConnectorMock) Connect(ctx context.Context) (driver.Conn, error) {
	args := m.Called(ctx)

	if conn, ok := args.Get(0).(driver.Conn); ok {
		return conn, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (m *baseConnectorMock) Driver() driver.Driver {
	args := m.Called()

	if connector, ok := args.Get(0).(driver.Driver); ok {
		return connector
	} else {
		return nil
	}
}

func TestNewConnector(t *testing.T) {
	t.Parallel()

	connectorMock := new(baseConnectorMock)

	c := sql.NewConnector("test dsn", connectorMock, nil)

	assert.IsType(t, &sql.Connector{}, c)
	assert.Nil(t, c.Driver())

	connectorMock.AssertExpectations(t)
}

func TestConnect(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	driverMock := new(baseDriverMock)

	d := sql.NewDriver(driverMock, sql.NewConfiguration(sql.SqliteSystem))

	connectorMock := new(baseConnectorMock)
	connectorMock.On("Connect", ctx).Return(nil, nil)

	c := sql.NewConnector("test dsn", connectorMock, d)

	conn, err := c.Connect(ctx)
	assert.NoError(t, err)

	assert.IsType(t, &sql.Connection{}, conn)

	driverMock.AssertExpectations(t)
	connectorMock.AssertExpectations(t)
}

func TestConnectWithError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	connectorMock := new(baseConnectorMock)
	connectorMock.On("Connect", ctx).Return(nil, fmt.Errorf("test error"))

	c := sql.NewConnector("test dsn", connectorMock, nil)

	_, err := c.Connect(ctx)
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	connectorMock.AssertExpectations(t)
}
