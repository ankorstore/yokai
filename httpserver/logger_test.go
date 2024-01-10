package httpserver_test

import (
	"io"
	"os"
	"os/exec"
	"testing"

	"github.com/ankorstore/yokai/httpserver"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/labstack/echo/v4"
	echologger "github.com/labstack/gommon/log"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
)

func TestNewEchoLogger(t *testing.T) {
	logger, err := log.NewDefaultLoggerFactory().Create()
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	assert.IsType(t, &httpserver.EchoLogger{}, echoLogger)
	assert.Implements(t, (*echo.Logger)(nil), echoLogger)
}

func TestLevelConversion(t *testing.T) {
	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithLevel(zerolog.TraceLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	assert.Equal(t, echologger.DEBUG, echoLogger.Level())
}

func TestPrefix(t *testing.T) {
	logger, err := log.NewDefaultLoggerFactory().Create()
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	echoLogger.SetPrefix("prefix")
	assert.Equal(t, "prefix", echoLogger.Prefix())
}

func TestToZerolog(t *testing.T) {
	logger, err := log.NewDefaultLoggerFactory().Create()
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	assert.IsType(t, &zerolog.Logger{}, echoLogger.ToZerolog())
}

func TestOutput(t *testing.T) {
	logger, err := log.NewDefaultLoggerFactory().Create()
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	buffer := logtest.NewDefaultTestLogBuffer()
	echoLogger.SetOutput(buffer)

	assert.Implements(t, (*io.Writer)(nil), echoLogger.Output())
}

func TestLevel(t *testing.T) {
	logger, err := log.NewDefaultLoggerFactory().Create()
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	echoLogger.SetLevel(echologger.DEBUG)
	assert.Equal(t, echologger.DEBUG, echoLogger.Level())

	echoLogger.SetLevel(echologger.INFO)
	assert.Equal(t, echologger.INFO, echoLogger.Level())

	echoLogger.SetLevel(echologger.WARN)
	assert.Equal(t, echologger.WARN, echoLogger.Level())

	echoLogger.SetLevel(echologger.ERROR)
	assert.Equal(t, echologger.ERROR, echoLogger.Level())

	echoLogger.SetLevel(echologger.OFF)
	assert.Equal(t, echologger.OFF, echoLogger.Level())
}

func TestHeader(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)
	echoLogger.SetHeader("test header")
	echoLogger.Info("test message")

	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "info",
		"header":  "test header",
		"message": "test message",
	})
}

func TestDebugLogging(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
		log.WithLevel(zerolog.DebugLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	echoLogger.Debug("test regular message")
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "debug",
		"message": "test regular message",
	})

	echoLogger.Debugf("test placeholder message: %s", "placeholder")
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "debug",
		"message": "test placeholder message: placeholder",
	})

	echoLogger.Debugj(echologger.JSON{"message": "test json message"})
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "debug",
		"message": "test json message",
	})
}

func TestInfoLogging(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
		log.WithLevel(zerolog.InfoLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	echoLogger.Info("test regular message")
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "info",
		"message": "test regular message",
	})

	echoLogger.Infof("test placeholder message: %s", "placeholder")
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "info",
		"message": "test placeholder message: placeholder",
	})

	echoLogger.Infoj(echologger.JSON{"message": "test json message"})
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "info",
		"message": "test json message",
	})
}

func TestWarnLogging(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
		log.WithLevel(zerolog.WarnLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	echoLogger.Warn("test regular message")
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "warn",
		"message": "test regular message",
	})

	echoLogger.Warnf("test placeholder message: %s", "placeholder")
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "warn",
		"message": "test placeholder message: placeholder",
	})

	echoLogger.Warnj(echologger.JSON{"message": "test json message"})
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "warn",
		"message": "test json message",
	})
}

func TestErrorLogging(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
		log.WithLevel(zerolog.ErrorLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	echoLogger.Error("test regular message")
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "error",
		"message": "test regular message",
	})

	echoLogger.Errorf("test placeholder message: %s", "placeholder")
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "error",
		"message": "test placeholder message: placeholder",
	})

	echoLogger.Errorj(echologger.JSON{"message": "test json message"})
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "error",
		"message": "test json message",
	})
}

func TestPrintLogging(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
		log.WithLevel(zerolog.FatalLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	echoLogger.Print("test regular message")
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "-",
		"message": "test regular message",
	})

	echoLogger.Printf("test placeholder message: %s", "placeholder")
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "-",
		"message": "test placeholder message: placeholder",
	})

	echoLogger.Printj(echologger.JSON{"message": "test json message"})
	logtest.AssertHasLogRecord(t, buffer, map[string]interface{}{
		"level":   "-",
		"message": "test json message",
	})
}

func TestFatalLogging(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
		log.WithLevel(zerolog.FatalLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	if os.Getenv("SHOULD_FATAL") == "1" {
		echoLogger.Fatal("message")
	}

	//nolint:gosec
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalLogging")
	cmd.Env = append(os.Environ(), "SHOULD_FATAL=1")

	err = cmd.Run()

	//nolint:errorlint
	if e, ok := err.(*exec.ExitError); ok && e.Success() {
		t.Errorf("test ran with err %v, want exit status 1", err)
	} else {
		return
	}
}

func TestFatalFLogging(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
		log.WithLevel(zerolog.FatalLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	if os.Getenv("SHOULD_FATAL_F") == "1" {
		echoLogger.Fatalf("test placeholder message: %s", "placeholder")
	}

	//nolint:gosec
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalFLogging")
	cmd.Env = append(os.Environ(), "SHOULD_FATAL_F=1")

	err = cmd.Run()

	//nolint:errorlint
	if e, ok := err.(*exec.ExitError); ok && e.Success() {
		t.Errorf("test ran with err %v, want exit status 1", err)
	} else {
		return
	}
}

func TestFatalJLogging(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
		log.WithLevel(zerolog.FatalLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	if os.Getenv("SHOULD_FATAL_J") == "1" {
		echoLogger.Fatalj(echologger.JSON{"message": "message"})
	}

	//nolint:gosec
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalJLogging")
	cmd.Env = append(os.Environ(), "SHOULD_FATAL_J=1")

	err = cmd.Run()

	//nolint:errorlint
	if e, ok := err.(*exec.ExitError); ok && e.Success() {
		t.Errorf("test ran with err %v, want exit status 1", err)
	} else {
		return
	}
}

func TestPanicLogging(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
		log.WithLevel(zerolog.FatalLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "test panic", r)
		} else {
			t.Errorf("logger did not panic")
		}
	}()

	echoLogger.Panic("test panic")
}

func TestPanicFLogging(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
		log.WithLevel(zerolog.FatalLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "test panic f", r)
		} else {
			t.Errorf("logger did not panic panic f")
		}
	}()

	echoLogger.Panicf("test %s", "panic f")
}

func TestPanicJLogging(t *testing.T) {
	buffer := logtest.NewDefaultTestLogBuffer()

	logger, err := log.NewDefaultLoggerFactory().Create(
		log.WithOutputWriter(buffer),
		log.WithLevel(zerolog.FatalLevel),
	)
	assert.NoError(t, err)

	echoLogger := httpserver.NewEchoLogger(logger)

	defer func() {
		if r := recover(); r != nil {
			assert.Equal(t, "", r)
		} else {
			t.Errorf("logger did not panic on panic j")
		}
	}()

	echoLogger.Panicj(echologger.JSON{"message": "test panic j"})
}
