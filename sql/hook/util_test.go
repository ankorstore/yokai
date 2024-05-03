package hook_test

import (
	"testing"

	"github.com/ankorstore/yokai/sql/hook"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		list []string
		item string
		want bool
	}{
		{
			name: "contains at beginning",
			list: []string{"foo", "bar", "baz"},
			item: "foo",
			want: true,
		},
		{
			name: "contains at end",
			list: []string{"foo", "bar", "baz"},
			item: "baz",
			want: true,
		},
		{
			name: "contains in middle",
			list: []string{"foo", "bar", "baz"},
			item: "bar",
			want: true,
		},
		{
			name: "not contains",
			list: []string{"foo", "bar", "baz"},
			item: "invalid",
			want: false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			assert.Equal(t, test.want, hook.Contains(test.list, test.item))
		})
	}
}
