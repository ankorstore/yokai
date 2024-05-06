package sql_test

import (
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/stretchr/testify/assert"
)

func TestOperationAsString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		operation sql.Operation
		expected  string
	}{
		{
			sql.ConnectionBeginOperation,
			"connection:begin",
		},
		{
			sql.ConnectionBeginTxOperation,
			"connection:begin-tx",
		},
		{
			sql.ConnectionExecOperation,
			"connection:exec",
		},
		{
			sql.ConnectionExecContextOperation,
			"connection:exec-context",
		},
		{
			sql.ConnectionQueryOperation,
			"connection:query",
		},
		{
			sql.ConnectionQueryContextOperation,
			"connection:query-context",
		},
		{
			sql.ConnectionPrepareOperation,
			"connection:prepare",
		},
		{
			sql.ConnectionPrepareContextOperation,
			"connection:prepare-context",
		},
		{
			sql.ConnectionPingOperation,
			"connection:ping",
		},
		{
			sql.ConnectionResetSessionOperation,
			"connection:reset-session",
		},
		{
			sql.ConnectionCloseOperation,
			"connection:close",
		},
		{
			sql.TransactionCommitOperation,
			"transaction:commit",
		},
		{
			sql.TransactionRollbackOperation,
			"transaction:rollback",
		},
		{
			sql.StatementExecOperation,
			"statement:exec",
		},
		{
			sql.StatementQueryOperation,
			"statement:query",
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, test.operation.String())
	}
}

func TestFetchOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		operation string
		expected  sql.Operation
	}{
		{
			"connection:begin",
			sql.ConnectionBeginOperation,
		},
		{
			"connection:begin-tx",
			sql.ConnectionBeginTxOperation,
		},
		{
			"connection:exec",
			sql.ConnectionExecOperation,
		},
		{
			"connection:exec-context",
			sql.ConnectionExecContextOperation,
		},
		{
			"connection:query",
			sql.ConnectionQueryOperation,
		},
		{
			"connection:query-context",
			sql.ConnectionQueryContextOperation,
		},
		{
			"connection:prepare",
			sql.ConnectionPrepareOperation,
		},
		{
			"connection:prepare-context",
			sql.ConnectionPrepareContextOperation,
		},
		{
			"connection:ping",
			sql.ConnectionPingOperation,
		},
		{
			"connection:reset-session",
			sql.ConnectionResetSessionOperation,
		},
		{
			"connection:close",
			sql.ConnectionCloseOperation,
		},
		{
			"transaction:commit",
			sql.TransactionCommitOperation,
		},
		{
			"transaction:rollback",
			sql.TransactionRollbackOperation,
		},
		{
			"statement:exec",
			sql.StatementExecOperation,
		},
		{
			"statement:query",
			sql.StatementQueryOperation,
		},
		{
			"",
			sql.UnknownOperation,
		},
		{
			"unknown",
			sql.UnknownOperation,
		},
		{
			"invalid",
			sql.UnknownOperation,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, sql.FetchOperation(test.operation))
	}
}

func TestContainsOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		operations []sql.Operation
		operation  sql.Operation
		want       bool
	}{
		{
			// contains at beginning of list
			operations: []sql.Operation{
				sql.ConnectionPingOperation,
				sql.ConnectionResetSessionOperation,
				sql.ConnectionCloseOperation,
			},
			operation: sql.ConnectionPingOperation,
			want:      true,
		},
		{
			// contains at end of list
			operations: []sql.Operation{
				sql.ConnectionPingOperation,
				sql.ConnectionResetSessionOperation,
				sql.ConnectionCloseOperation,
			},
			operation: sql.ConnectionCloseOperation,
			want:      true,
		},
		{
			// contains in middle of list
			operations: []sql.Operation{
				sql.ConnectionPingOperation,
				sql.ConnectionResetSessionOperation,
				sql.ConnectionCloseOperation,
			},
			operation: sql.ConnectionResetSessionOperation,
			want:      true,
		},
		{
			// contains in single item list
			operations: []sql.Operation{
				sql.ConnectionPingOperation,
			},
			operation: sql.ConnectionPingOperation,
			want:      true,
		},
		{
			// not contains in list
			operations: []sql.Operation{
				sql.ConnectionPingOperation,
				sql.ConnectionResetSessionOperation,
				sql.ConnectionCloseOperation,
			},
			operation: sql.TransactionCommitOperation,
			want:      false,
		},
		{
			// not contains in empty list
			operations: []sql.Operation{},
			operation:  sql.TransactionCommitOperation,
			want:       false,
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.want, sql.ContainsOperation(test.operations, test.operation))
	}
}
