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
	StatementExecOperation            Operation = "statement:exec"
	StatementExecContextOperation     Operation = "statement:exec-context"
	StatementQueryOperation           Operation = "statement:query"
	StatementQueryContextOperation    Operation = "statement:query-context"
	TransactionCommitOperation        Operation = "transaction:commit"
	TransactionRollbackOperation      Operation = "transaction:rollback"
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
		StatementExecOperation,
		StatementExecContextOperation,
		StatementQueryOperation,
		StatementQueryContextOperation,
		TransactionCommitOperation,
		TransactionRollbackOperation:
		return o
	default:
		return UnknownOperation
	}
}

// ContainsOperation returns true if a given Operation item is contained id a list of Operation.
func ContainsOperation(list []Operation, item Operation) bool {
	for _, listItem := range list {
		if listItem == item {
			return true
		}
	}

	return false
}
