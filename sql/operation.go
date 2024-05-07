package sql

import "strings"

// Operation is an enum for the supported database operations.
type Operation string

const (
	UnknownOperation                  Operation = "unknown"
	ConnectionBeginOperation          Operation = "connection:begin"
	ConnectionBeginTxOperation        Operation = "connection:begin-tx"
	ConnectionExecOperation           Operation = "connection:exec"
	ConnectionExecContextOperation    Operation = "connection:exec-context"
	ConnectionQueryOperation          Operation = "connection:query"
	ConnectionQueryContextOperation   Operation = "connection:query-context"
	ConnectionPrepareOperation        Operation = "connection:prepare"
	ConnectionPrepareContextOperation Operation = "connection:prepare-context"
	ConnectionPingOperation           Operation = "connection:ping"
	ConnectionResetSessionOperation   Operation = "connection:reset-session"
	ConnectionCloseOperation          Operation = "connection:close"
	TransactionCommitOperation        Operation = "transaction:commit"
	TransactionRollbackOperation      Operation = "transaction:rollback"
	StatementExecOperation            Operation = "statement:exec"
	StatementExecContextOperation     Operation = "statement:exec-context"
	StatementQueryOperation           Operation = "statement:query"
	StatementQueryContextOperation    Operation = "statement:query-context"
)

// String returns a string representation of the Operation.
func (o Operation) String() string {
	return string(o)
}

// FetchOperation returns an Operation for a given name.
func FetchOperation(name string) Operation {
	//nolint:exhaustive
	switch o := Operation(strings.ToLower(name)); o {
	case ConnectionBeginOperation,
		ConnectionBeginTxOperation,
		ConnectionExecOperation,
		ConnectionExecContextOperation,
		ConnectionQueryOperation,
		ConnectionQueryContextOperation,
		ConnectionPrepareOperation,
		ConnectionPrepareContextOperation,
		ConnectionPingOperation,
		ConnectionResetSessionOperation,
		ConnectionCloseOperation,
		TransactionCommitOperation,
		TransactionRollbackOperation,
		StatementExecOperation,
		StatementExecContextOperation,
		StatementQueryOperation,
		StatementQueryContextOperation:
		return o
	default:
		return UnknownOperation
	}
}

func ContainsOperation(operations []Operation, operation Operation) bool {
	for _, operationsItem := range operations {
		if operationsItem == operation {
			return true
		}
	}

	return false
}
