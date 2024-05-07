package sql_test

import (
	"context"
	"database/sql/driver"
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/ankorstore/yokai/sql/hook/hooktest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type baseConnMock struct {
	mock.Mock
}

func (m *baseConnMock) Exec(query string, args []driver.Value) (driver.Result, error) {
	callArgs := m.Called(query, args)

	if res, ok := callArgs.Get(0).(driver.Result); ok {
		return res, callArgs.Error(1)
	} else {
		return nil, callArgs.Error(1)
	}
}

func (m *baseConnMock) Query(query string, args []driver.Value) (driver.Rows, error) {
	callArgs := m.Called(query, args)

	if rows, ok := callArgs.Get(0).(driver.Rows); ok {
		return rows, callArgs.Error(1)
	} else {
		return nil, callArgs.Error(1)
	}
}

func (m *baseConnMock) Prepare(query string) (driver.Stmt, error) {
	args := m.Called(query)

	if stmt, ok := args.Get(0).(driver.Stmt); ok {
		return stmt, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (m *baseConnMock) Begin() (driver.Tx, error) {
	args := m.Called()

	if tx, ok := args.Get(0).(driver.Tx); ok {
		return tx, args.Error(1)
	} else {
		return nil, args.Error(1)
	}
}

func (m *baseConnMock) ResetSession(ctx context.Context) error {
	return m.Called(ctx).Error(0)
}

func (m *baseConnMock) Close() error {
	return m.Called().Error(0)
}

func TestConnExec(t *testing.T) {
	t.Parallel()

	query := hooktest.TestHookEventQuery
	arguments := []driver.Value{hooktest.TestHookEventArgument}

	resultMock := new(baseResultMock)
	resultMock.On("LastInsertId").Return(int64(0), nil)
	resultMock.On("RowsAffected").Return(int64(1), nil)

	connMock := new(baseConnMock)
	connMock.On("Exec", query, arguments).Return(resultMock, nil)

	c := sql.NewConnection(connMock, &sql.Configuration{})

	res, err := c.Exec(query, arguments)
	assert.NoError(t, err)

	lastInsertId, err := res.LastInsertId()
	assert.NoError(t, err)
	assert.Equal(t, int64(0), lastInsertId)

	rowsAffected, err := res.RowsAffected()
	assert.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	resultMock.AssertExpectations(t)
	connMock.AssertExpectations(t)
}

func TestConnExecError(t *testing.T) {
	t.Parallel()

	query := hooktest.TestHookEventQuery
	arguments := []driver.Value{hooktest.TestHookEventArgument}

	connMock := new(baseConnMock)
	connMock.On("Exec", query, arguments).Return(nil, fmt.Errorf("test error"))

	c := sql.NewConnection(connMock, &sql.Configuration{})

	_, err := c.Exec(query, arguments)
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	connMock.AssertExpectations(t)
}

func TestConnQuery(t *testing.T) {
	t.Parallel()

	query := hooktest.TestHookEventQuery
	arguments := []driver.Value{hooktest.TestHookEventArgument}

	rowsMock := new(baseRowsMock)
	connMock := new(baseConnMock)
	connMock.On("Query", query, arguments).Return(rowsMock, nil)

	c := sql.NewConnection(connMock, &sql.Configuration{})

	_, err := c.Query(query, arguments)
	assert.NoError(t, err)

	rowsMock.AssertExpectations(t)
	connMock.AssertExpectations(t)
}

func TestConnQueryError(t *testing.T) {
	t.Parallel()

	query := hooktest.TestHookEventQuery
	arguments := []driver.Value{hooktest.TestHookEventArgument}

	connMock := new(baseConnMock)
	connMock.On("Query", query, arguments).Return(nil, fmt.Errorf("test error"))

	c := sql.NewConnection(connMock, &sql.Configuration{})

	_, err := c.Query(query, arguments)
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	connMock.AssertExpectations(t)
}

func TestConnPrepare(t *testing.T) {
	t.Parallel()

	query := "SELECT * FROM foo WHERE id = ?"

	stmtMock := new(baseStmtMock)
	connMock := new(baseConnMock)
	connMock.On("Prepare", query).Return(stmtMock, nil)

	c := sql.NewConnection(connMock, &sql.Configuration{})

	_, err := c.Prepare(query)
	assert.NoError(t, err)

	stmtMock.AssertExpectations(t)
	connMock.AssertExpectations(t)
}

func TestConnPrepareError(t *testing.T) {
	t.Parallel()

	query := hooktest.TestHookEventQuery

	connMock := new(baseConnMock)
	connMock.On("Prepare", query).Return(nil, fmt.Errorf("test error"))

	c := sql.NewConnection(connMock, &sql.Configuration{})

	_, err := c.Prepare(query)
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	connMock.AssertExpectations(t)
}

func TestConnBegin(t *testing.T) {
	t.Parallel()

	txMock := new(baseTxMock)
	connMock := new(baseConnMock)
	connMock.On("Begin").Return(txMock, nil)

	c := sql.NewConnection(connMock, &sql.Configuration{})

	_, err := c.Begin()
	assert.NoError(t, err)

	txMock.AssertExpectations(t)
	connMock.AssertExpectations(t)
}

func TestConnBeginError(t *testing.T) {
	t.Parallel()

	connMock := new(baseConnMock)
	connMock.On("Begin").Return(nil, fmt.Errorf("test error"))

	c := sql.NewConnection(connMock, &sql.Configuration{})

	_, err := c.Begin()
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	connMock.AssertExpectations(t)
}

func TestConnResetSession(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	connMock := new(baseConnMock)
	connMock.On("ResetSession", ctx).Return(nil)

	c := sql.NewConnection(connMock, &sql.Configuration{})

	err := c.ResetSession(ctx)
	assert.NoError(t, err)

	connMock.AssertExpectations(t)
}

func TestConnResetSessionError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	connMock := new(baseConnMock)
	connMock.On("ResetSession", ctx).Return(fmt.Errorf("test error"))

	c := sql.NewConnection(connMock, &sql.Configuration{})

	err := c.ResetSession(ctx)
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	connMock.AssertExpectations(t)
}
