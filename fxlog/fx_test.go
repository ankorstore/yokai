package fxlog_test

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/fx/fxtest"
)

func TestFxEventLoggerWithinFx(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TEST_LOG_LEVEL", "debug")
	t.Setenv("TEST_LOG_OUTPUT", "test")

	var buffer logtest.TestLogBuffer

	fxtest.New(
		t,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fx.WithLogger(fxlog.NewFxEventLogger),
		fx.Populate(&buffer),
	).RequireStart().RequireStop()

	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":       "info",
		"service":     "dev",
		"message":     "provided",
		"constructor": "github.com/ankorstore/yokai/fxlog.NewFxLogger()",
	})
}

func TestFxEventLogger(t *testing.T) {
	t.Parallel()

	someError := errors.New("some error")

	tests := []struct {
		name        string
		give        fxevent.Event
		wantMessage string
	}{
		{
			name: "OnStartExecuting",
			give: &fxevent.OnStartExecuting{
				FunctionName: "hook.onStart",
				CallerName:   "bytes.NewBuffer",
			},
			wantMessage: "{\"level\":\"info\",\"callee\":\"hook.onStart\",\"caller\":\"bytes.NewBuffer\",\"message\":\"OnStart hook executing\"}\n",
		},
		{
			name: "OnStopExecuting",
			give: &fxevent.OnStopExecuting{
				FunctionName: "hook.onStop1",
				CallerName:   "bytes.NewBuffer",
			},
			wantMessage: "{\"level\":\"info\",\"callee\":\"hook.onStop1\",\"caller\":\"bytes.NewBuffer\",\"message\":\"OnStop hook executing\"}\n",
		},
		{

			name: "OnStopExecutedError",
			give: &fxevent.OnStopExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Err:          fmt.Errorf("some error"),
			},
			wantMessage: "{\"level\":\"warn\",\"error\":\"some error\",\"callee\":\"hook.onStart1\",\"callee\":\"bytes.NewBuffer\",\"message\":\"OnStop hook failed\"}\n",
		},
		{
			name: "OnStopExecuted",
			give: &fxevent.OnStopExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Runtime:      time.Millisecond * 3,
			},
			wantMessage: "{\"level\":\"info\",\"callee\":\"hook.onStart1\",\"caller\":\"bytes.NewBuffer\",\"runtime\":\"3ms\",\"message\":\"OnStop hook executed\"}\n",
		},
		{

			name: "OnStartExecutedError",
			give: &fxevent.OnStartExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Err:          fmt.Errorf("some error"),
			},
			wantMessage: "{\"level\":\"warn\",\"error\":\"some error\",\"callee\":\"hook.onStart1\",\"caller\":\"bytes.NewBuffer\",\"message\":\"OnStart hook failed\"}\n",
		},
		{
			name: "OnStartExecuted",
			give: &fxevent.OnStartExecuted{
				FunctionName: "hook.onStart1",
				CallerName:   "bytes.NewBuffer",
				Runtime:      time.Millisecond * 3,
			},
			wantMessage: "{\"level\":\"info\",\"callee\":\"hook.onStart1\",\"caller\":\"bytes.NewBuffer\",\"runtime\":\"3ms\",\"message\":\"OnStart hook executed\"}\n",
		},
		{
			name:        "Supplied",
			give:        &fxevent.Supplied{TypeName: "*bytes.Buffer"},
			wantMessage: "{\"level\":\"info\",\"type\":\"*bytes.Buffer\",\"message\":\"supplied\"}\n",
		},
		{
			name:        "SuppliedError",
			give:        &fxevent.Supplied{TypeName: "*bytes.Buffer", Err: someError},
			wantMessage: "{\"level\":\"warn\",\"error\":\"some error\",\"type\":\"*bytes.Buffer\",\"message\":\"supplied\"}\n",
		},
		{
			name: "Provide",
			give: &fxevent.Provided{
				ConstructorName: "bytes.NewBuffer()",
				OutputTypeNames: []string{"*bytes.Buffer"},
			},
			wantMessage: "{\"level\":\"info\",\"type\":\"*bytes.Buffer\",\"constructor\":\"bytes.NewBuffer()\",\"message\":\"provided\"}\n",
		},
		{
			name:        "Provide with Error",
			give:        &fxevent.Provided{Err: someError},
			wantMessage: "{\"level\":\"error\",\"error\":\"some error\",\"message\":\"error encountered while applying options\"}\n",
		},
		{
			name:        "Invoked/Success",
			give:        &fxevent.Invoked{FunctionName: "bytes.NewBuffer()"},
			wantMessage: "{\"level\":\"info\",\"function\":\"bytes.NewBuffer()\",\"message\":\"invoked\"}\n",
		},
		{
			name:        "Invoked/Error",
			give:        &fxevent.Invoked{FunctionName: "bytes.NewBuffer()", Err: someError},
			wantMessage: "{\"level\":\"error\",\"error\":\"some error\",\"stack\":\"\",\"function\":\"bytes.NewBuffer()\",\"message\":\"invoke failed\"}\n",
		},
		{
			name:        "StartError",
			give:        &fxevent.Started{Err: someError},
			wantMessage: "{\"level\":\"error\",\"error\":\"some error\",\"message\":\"start failed\"}\n",
		},
		{
			name:        "Stopping",
			give:        &fxevent.Stopping{Signal: os.Interrupt},
			wantMessage: "{\"level\":\"info\",\"signal\":\"INTERRUPT\",\"message\":\"received signal\"}\n",
		},
		{
			name:        "Stopped",
			give:        &fxevent.Stopped{Err: someError},
			wantMessage: "{\"level\":\"error\",\"error\":\"some error\",\"message\":\"stop failed\"}\n",
		},
		{
			name:        "RollingBack",
			give:        &fxevent.RollingBack{StartErr: someError},
			wantMessage: "{\"level\":\"error\",\"error\":\"some error\",\"message\":\"start failed, rolling back\"}\n",
		},
		{
			name:        "RolledBackError",
			give:        &fxevent.RolledBack{Err: someError},
			wantMessage: "{\"level\":\"error\",\"error\":\"some error\",\"message\":\"rollback failed\"}\n",
		},
		{
			name:        "Started",
			give:        &fxevent.Started{},
			wantMessage: "{\"level\":\"info\",\"message\":\"started\"}\n",
		},
		{
			name:        "LoggerInitialized Error",
			give:        &fxevent.LoggerInitialized{Err: someError},
			wantMessage: "{\"level\":\"error\",\"error\":\"some error\",\"message\":\"custom logger initialization failed\"}\n",
		},
		{
			name:        "LoggerInitialized",
			give:        &fxevent.LoggerInitialized{ConstructorName: "bytes.NewBuffer()"},
			wantMessage: "{\"level\":\"info\",\"function\":\"bytes.NewBuffer()\",\"message\":\"initialized custom fxevent.Logger\"}\n",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			buf := bytes.NewBufferString("")
			fxlog.NewFxEventLogger(log.FromZerolog(zerolog.New(buf))).LogEvent(tt.give)

			assert.Equal(t, tt.wantMessage, buf.String())
		})
	}
}
