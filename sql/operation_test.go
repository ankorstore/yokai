package sql_test

import (
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/stretchr/testify/assert"
)

func TestContainsOperation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		operations []sql.Operation
		operation  sql.Operation
		want       bool
	}{
		{
			name: "contains at beginning of list",
			operations: []sql.Operation{
				sql.ConnectionPingOperation,
				sql.ConnectionResetSessionOperation,
				sql.ConnectionCloseOperation,
			},
			operation: sql.ConnectionPingOperation,
			want:      true,
		},
		{
			name: "contains at end of list",
			operations: []sql.Operation{
				sql.ConnectionPingOperation,
				sql.ConnectionResetSessionOperation,
				sql.ConnectionCloseOperation,
			},
			operation: sql.ConnectionCloseOperation,
			want:      true,
		},
		{
			name: "contains in middle of list",
			operations: []sql.Operation{
				sql.ConnectionPingOperation,
				sql.ConnectionResetSessionOperation,
				sql.ConnectionCloseOperation,
			},
			operation: sql.ConnectionResetSessionOperation,
			want:      true,
		},
		{
			name: "contains in single item list",
			operations: []sql.Operation{
				sql.ConnectionPingOperation,
			},
			operation: sql.ConnectionPingOperation,
			want:      true,
		},
		{
			name: "not contains in list",
			operations: []sql.Operation{
				sql.ConnectionPingOperation,
				sql.ConnectionResetSessionOperation,
				sql.ConnectionCloseOperation,
			},
			operation: sql.TransactionCommitOperation,
			want:      false,
		},
		{
			name:       "not contains in empty list",
			operations: []sql.Operation{},
			operation:  sql.TransactionCommitOperation,
			want:       false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.want, sql.ContainsOperation(test.operations, test.operation))
		})
	}
}
