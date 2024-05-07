package sql_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type baseTxMock struct {
	mock.Mock
}

func (m *baseTxMock) Commit() error {
	args := m.Called()

	return args.Error(0)
}

func (m *baseTxMock) Rollback() error {
	args := m.Called()

	return args.Error(0)
}

func TestTxCommitError(t *testing.T) {
	t.Parallel()

	txMock := new(baseTxMock)

	txMock.On("Commit").Return(fmt.Errorf("test error"))

	tx := sql.NewTransaction(txMock, nil, &sql.Configuration{})

	err := tx.Commit()
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	txMock.AssertExpectations(t)
}

func TestTxRollbackError(t *testing.T) {
	t.Parallel()

	txMock := new(baseTxMock)

	txMock.On("Rollback").Return(fmt.Errorf("test error"))

	tx := sql.NewTransaction(txMock, nil, &sql.Configuration{})

	err := tx.Rollback()
	assert.Error(t, err)
	assert.Equal(t, "test error", err.Error())

	txMock.AssertExpectations(t)
}
