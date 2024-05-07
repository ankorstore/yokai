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

type baseStmtMock struct {
	mock.Mock
}

func (m *baseStmtMock) Close() error {
	args := m.Called()

	return args.Error(0)
}

func (m *baseStmtMock) NumInput() int {
	args := m.Called()

	return args.Int(0)
}

func (m *baseStmtMock) Exec(args []driver.Value) (driver.Result, error) {
	callArgs := m.Called(args)

	if res, ok := callArgs.Get(0).(driver.Result); ok {
		return res, callArgs.Error(1)
	} else {
		return nil, callArgs.Error(1)
	}
}

func (m *baseStmtMock) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	callArgs := m.Called(ctx, args)

	if res, ok := callArgs.Get(0).(driver.Result); ok {
		return res, callArgs.Error(1)
	} else {
		return nil, callArgs.Error(1)
	}
}

func (m *baseStmtMock) Query(args []driver.Value) (driver.Rows, error) {
	callArgs := m.Called(args)

	if rows, ok := callArgs.Get(0).(driver.Rows); ok {
		return rows, callArgs.Error(1)
	} else {
		return nil, callArgs.Error(1)
	}
}

func (m *baseStmtMock) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	callArgs := m.Called(ctx, args)

	if rows, ok := callArgs.Get(0).(driver.Rows); ok {
		return rows, callArgs.Error(1)
	} else {
		return nil, callArgs.Error(1)
	}
}

type baseResultMock struct {
	mock.Mock
}

func (m *baseResultMock) LastInsertId() (int64, error) {
	args := m.Called()

	if lastInsertId, ok := args.Get(0).(int64); ok {
		return lastInsertId, args.Error(1)
	} else {
		return 0, args.Error(1)
	}
}

func (m *baseResultMock) RowsAffected() (int64, error) {
	args := m.Called()

	if rowsAffected, ok := args.Get(0).(int64); ok {
		return rowsAffected, args.Error(1)
	} else {
		return 0, args.Error(1)
	}
}

type baseRowsMock struct {
	mock.Mock
}

func (m *baseRowsMock) Columns() []string {
	args := m.Called()

	if cols, ok := args.Get(0).([]string); ok {
		return cols
	} else {
		return []string{}
	}
}

func (m *baseRowsMock) Close() error {
	args := m.Called()

	return args.Error(0)
}

func (m *baseRowsMock) Next(dest []driver.Value) error {
	args := m.Called(dest)

	return args.Error(0)
}

func TestStmtExec(t *testing.T) {
	t.Parallel()

	stmtMock := new(baseStmtMock)
	resultMock := new(baseResultMock)

	stmtMock.On("Exec", mock.Anything).Return(resultMock, nil)
	resultMock.On("LastInsertId").Return(int64(1), nil)
	resultMock.On("RowsAffected").Return(int64(1), nil)

	s := sql.NewStatement(stmtMock, nil, "INSERT INTO foo (bar) VALUES (?)", &sql.Configuration{})

	res, err := s.Exec([]driver.Value{"value"})
	assert.NoError(t, err)
	assert.Implements(t, (*driver.Result)(nil), res)

	stmtMock.AssertExpectations(t)
	resultMock.AssertExpectations(t)
}

func TestStmtExecError(t *testing.T) {
	t.Parallel()

	stmtMock := new(baseStmtMock)
	resultMock := new(baseResultMock)

	stmtMock.On("Exec", mock.Anything).Return(nil, fmt.Errorf("test error"))

	s := sql.NewStatement(stmtMock, nil, "INSERT INTO foo (bar) VALUES (?)", &sql.Configuration{})

	_, err := s.Exec([]driver.Value{"value"})
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	stmtMock.AssertExpectations(t)
	resultMock.AssertExpectations(t)
}

func TestStmtExecContextError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	stmtMock := new(baseStmtMock)
	resultMock := new(baseResultMock)

	stmtMock.On("ExecContext", ctx, mock.Anything).Return(nil, fmt.Errorf("test error"))

	s := sql.NewStatement(stmtMock, nil, "INSERT INTO foo (bar) VALUES (?)", &sql.Configuration{})

	_, err := s.ExecContext(ctx, []driver.NamedValue{})
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	stmtMock.AssertExpectations(t)
	resultMock.AssertExpectations(t)
}

func TestStmtQuery(t *testing.T) {
	t.Parallel()

	stmtMock := new(baseStmtMock)
	rowsMock := new(baseRowsMock)

	stmtMock.On("Query", mock.Anything).Return(rowsMock, nil)

	s := sql.NewStatement(stmtMock, nil, "SELECT * FROM foo", &sql.Configuration{})

	resultRows, err := s.Query([]driver.Value{})
	assert.NoError(t, err)
	assert.NotNil(t, resultRows)

	stmtMock.AssertExpectations(t)
	rowsMock.AssertExpectations(t)
}

func TestStmtQueryError(t *testing.T) {
	t.Parallel()

	stmtMock := new(baseStmtMock)
	rowsMock := new(baseRowsMock)

	stmtMock.On("Query", mock.Anything).Return(nil, fmt.Errorf("test error"))

	s := sql.NewStatement(stmtMock, nil, "SELECT * FROM foo", &sql.Configuration{})

	_, err := s.Query([]driver.Value{})
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	stmtMock.AssertExpectations(t)
	rowsMock.AssertExpectations(t)
}

func TestStmtQueryContextError(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	stmtMock := new(baseStmtMock)
	rowsMock := new(baseRowsMock)

	stmtMock.On("QueryContext", ctx, mock.Anything).Return(nil, fmt.Errorf("test error"))

	s := sql.NewStatement(stmtMock, nil, "SELECT * FROM foo", &sql.Configuration{})

	_, err := s.QueryContext(ctx, []driver.NamedValue{})
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	stmtMock.AssertExpectations(t)
	rowsMock.AssertExpectations(t)
}
