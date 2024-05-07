package sql

import (
	"database/sql/driver"
)

func ConvertNamedValuesToValues(namedValues []driver.NamedValue) []driver.Value {
	values := make([]driver.Value, len(namedValues))
	for name, value := range namedValues {
		values[name] = value.Value
	}

	return values
}
