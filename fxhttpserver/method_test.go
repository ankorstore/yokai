package fxhttpserver_test

import (
	"testing"

	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/stretchr/testify/assert"
)

func TestExtractMethods(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input          string
		expectedOutput []string
		expectedError  string
	}{
		"with CONNECT": {
			input:          "CONNECT",
			expectedOutput: []string{"CONNECT"},
		},
		"with GET": {
			input:          "GET",
			expectedOutput: []string{"GET"},
		},
		"with POST": {
			input:          "POST",
			expectedOutput: []string{"POST"},
		},
		"with PUT": {
			input:          "PUT",
			expectedOutput: []string{"PUT"},
		},
		"with PATCH": {
			input:          "PATCH",
			expectedOutput: []string{"PATCH"},
		},
		"with DELETE": {
			input:          "DELETE",
			expectedOutput: []string{"DELETE"},
		},
		"with HEAD": {
			input:          "HEAD",
			expectedOutput: []string{"HEAD"},
		},
		"with OPTIONS": {
			input:          "OPTIONS",
			expectedOutput: []string{"OPTIONS"},
		},
		"with TRACE": {
			input:          "TRACE",
			expectedOutput: []string{"TRACE"},
		},
		"with PROPFIND": {
			input:          "PROPFIND",
			expectedOutput: []string{"PROPFIND"},
		},
		"with REPORT": {
			input:          "REPORT",
			expectedOutput: []string{"REPORT"},
		},
		"with get": {
			input:          "get",
			expectedOutput: []string{"GET"},
		},
		"with GET,POST": {
			input:          "GET,POST",
			expectedOutput: []string{"GET", "POST"},
		},
		"with  GET , POST ": {
			input:          " GET , POST ",
			expectedOutput: []string{"GET", "POST"},
		},
		"with get,post": {
			input:          "get,post",
			expectedOutput: []string{"GET", "POST"},
		},
		"with *": {
			input: "*",
			expectedOutput: []string{
				"CONNECT",
				"DELETE",
				"GET",
				"HEAD",
				"OPTIONS",
				"PATCH",
				"POST",
				"PUT",
				"TRACE",
				"PROPFIND",
				"REPORT",
			},
		},
		"with invalid": {
			input:          "invalid",
			expectedOutput: nil,
			expectedError:  `invalid HTTP method "INVALID"`,
		},
		"with get,invalid,POST": {
			input:          "get,invalid,POST",
			expectedOutput: nil,
			expectedError:  `invalid HTTP method "INVALID"`,
		},
	}

	for tn, tt := range tests {
		tt := tt

		t.Run(tn, func(t *testing.T) {
			t.Parallel()

			output, err := fxhttpserver.ExtractMethods(tt.input)
			if err != nil {
				assert.Equal(t, tt.expectedError, err.Error())
			}

			assert.Equal(t, tt.expectedOutput, output)
		})
	}
}
