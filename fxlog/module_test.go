package fxlog_test

import (
	"io"
	"os"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxlog"
	"github.com/ankorstore/yokai/fxlog/testdata/factory"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModuleWithDebug(t *testing.T) {
	// even with error level configured, debug logs should be captured when APP_DEBUG=true
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_DEBUG", "true")
	t.Setenv("TEST_LOG_LEVEL", "error")
	t.Setenv("TEST_LOG_OUTPUT", "test")

	var buffer logtest.TestLogBuffer

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fx.Invoke(func(logger *log.Logger) {
			logger.Debug().Msg("test message")
		}),
		fx.Populate(&buffer),
	).RequireStart().RequireStop()

	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "debug",
		"service": "dev",
		"message": "test message",
	})
}

func TestModuleWithTestEnv(t *testing.T) {
	// test output writer should be used when APP_ENV=test, even if configured otherwise
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("APP_ENV", "test")
	t.Setenv("TEST_LOG_LEVEL", "debug")
	t.Setenv("TEST_LOG_OUTPUT", "stdout")

	var buffer logtest.TestLogBuffer

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fx.Invoke(func(logger *log.Logger) {
			logger.Debug().Msg("test message")
		}),
		fx.Populate(&buffer),
	).RequireStart().RequireStop()

	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "debug",
		"service": "test",
		"message": "test message",
	})
}

func TestModuleWithTestOutputWriter(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TEST_LOG_LEVEL", "debug")
	t.Setenv("TEST_LOG_OUTPUT", "test")

	prev := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	var buffer logtest.TestLogBuffer

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fx.Invoke(func(logger *log.Logger) {
			logger.Debug().Msg("test message")
		}),
		fx.Populate(&buffer),
	).RequireStart().RequireStop()

	// stdout should be empty
	err := w.Close()
	assert.NoError(t, err)
	out, _ := io.ReadAll(r)
	os.Stdout = prev
	assert.Empty(t, out)

	// test buffer should not be empty
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "debug",
		"service": "dev",
		"message": "test message",
	})
}

func TestModuleWithNoopOutputWriter(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TEST_LOG_LEVEL", "debug")
	t.Setenv("TEST_LOG_OUTPUT", "noop")

	prev := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	var buffer logtest.TestLogBuffer

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fx.Invoke(func(logger *log.Logger) {
			logger.Debug().Msg("test message")
		}),
		fx.Populate(&buffer),
	).RequireStart().RequireStop()

	// stdout should be empty
	err := w.Close()
	assert.NoError(t, err)
	out, _ := io.ReadAll(r)
	os.Stdout = prev
	assert.Empty(t, out)

	// test buffer should be empty
	hasRecord, _ := buffer.HasRecord(map[string]interface{}{
		"level":   "debug",
		"service": "dev",
		"message": "test message",
	})
	assert.False(t, hasRecord)
}

func TestModuleWithStdoutOutputWriter(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("TEST_LOG_LEVEL", "debug")
	t.Setenv("TEST_LOG_OUTPUT", "stdout")

	prev := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	var buffer logtest.TestLogBuffer

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fx.Invoke(func(logger *log.Logger) {
			logger.Debug().Msg("test message")
		}),
		fx.Populate(&buffer),
	).RequireStart().RequireStop()

	// stdout should not be empty
	err := w.Close()
	assert.NoError(t, err)
	out, _ := io.ReadAll(r)
	os.Stdout = prev
	assert.NotEmpty(t, out)
	assert.Contains(t, string(out), "test message")

	// test buffer should be empty
	hasRecord, _ := buffer.HasRecord(map[string]interface{}{
		"level":   "debug",
		"service": "dev",
		"message": "test message",
	})
	assert.False(t, hasRecord)
}

func TestModuleDecoration(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var logger *log.Logger

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fxlog.FxLogModule,
		fx.Decorate(factory.NewTestLoggerFactory),
		fx.Populate(&logger),
	).RequireStart().RequireStop()

	assert.Equal(t, &log.Logger{}, logger)
}
