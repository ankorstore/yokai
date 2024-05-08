package sql_test

import (
	"database/sql/driver"
	"testing"

	"github.com/ankorstore/yokai/sql"
	"github.com/stretchr/testify/assert"
)

func TestConvertNamedValuesToValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		nameValues []driver.NamedValue
		expected   []driver.Value
	}{
		// 2 items list
		{
			[]driver.NamedValue{
				{
					Name:    "foo",
					Ordinal: 1,
					Value:   "foo",
				},
				{
					Name:    "bar",
					Ordinal: 2,
					Value:   "bar",
				},
			},
			[]driver.Value{
				"foo",
				"bar",
			},
		},
		// 1 item list
		{
			[]driver.NamedValue{
				{
					Name:    "foo",
					Ordinal: 1,
					Value:   "foo",
				},
			},
			[]driver.Value{
				"foo",
			},
		},
		// empty list
		{
			[]driver.NamedValue{},
			[]driver.Value{},
		},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, sql.ConvertNamedValuesToValues(test.nameValues))
	}
}
