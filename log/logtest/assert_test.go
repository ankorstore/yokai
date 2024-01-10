package logtest_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/log/logtest"
	"github.com/stretchr/testify/assert"
)

func TestAssertHasLogRecord(t *testing.T) {
	t.Parallel()

	t.Run("test AssertHasLogRecord failure on error", func(t *testing.T) {
		t.Parallel()

		mt := new(testing.T)
		ce := fmt.Errorf("custom error")

		testLogBufferMock := new(TestLogBufferMock)
		testLogBufferMock.On("HasRecord").Return(false, ce)
		logtest.AssertHasLogRecord(mt, testLogBufferMock, map[string]interface{}{})

		assert.True(t, mt.Failed())
	})

	t.Run("test AssertHasLogRecord failure on missing attribute", func(t *testing.T) {
		t.Parallel()

		mt := new(testing.T)

		testLogBufferMock := new(TestLogBufferMock)
		testLogBufferMock.On("HasRecord").Return(false, nil)
		logtest.AssertHasLogRecord(mt, testLogBufferMock, map[string]interface{}{})

		assert.True(t, mt.Failed())
	})
}

func TestAssertNotHasLogRecord(t *testing.T) {
	t.Parallel()

	t.Run("test AssertHasNotLogRecord failure on error", func(t *testing.T) {
		t.Parallel()

		mt := new(testing.T)
		ce := fmt.Errorf("custom error")

		testLogBufferMock := new(TestLogBufferMock)
		testLogBufferMock.On("HasRecord").Return(false, ce)
		logtest.AssertHasNotLogRecord(mt, testLogBufferMock, map[string]interface{}{})

		assert.True(t, mt.Failed())
	})

	t.Run("test AssertHasNotLogRecord failure on missing attribute", func(t *testing.T) {
		t.Parallel()

		mt := new(testing.T)

		testLogBufferMock := new(TestLogBufferMock)
		testLogBufferMock.On("HasRecord").Return(true, nil)
		logtest.AssertHasNotLogRecord(mt, testLogBufferMock, map[string]interface{}{})

		assert.True(t, mt.Failed())
	})
}

func TestAssertContainLogRecord(t *testing.T) {
	t.Parallel()

	t.Run("test AssertContainLogRecord failure on error", func(t *testing.T) {
		t.Parallel()

		mt := new(testing.T)
		ce := fmt.Errorf("custom error")

		testLogBufferMock := new(TestLogBufferMock)
		testLogBufferMock.On("ContainRecord").Return(false, ce)
		logtest.AssertContainLogRecord(mt, testLogBufferMock, map[string]interface{}{})

		assert.True(t, mt.Failed())
	})

	t.Run("test AssertContainLogRecord failure on missing attribute", func(t *testing.T) {
		t.Parallel()

		mt := new(testing.T)

		testLogBufferMock := new(TestLogBufferMock)
		testLogBufferMock.On("ContainRecord").Return(false, nil)
		logtest.AssertContainLogRecord(mt, testLogBufferMock, map[string]interface{}{})

		assert.True(t, mt.Failed())
	})
}

func TestAssertContainNotLogRecord(t *testing.T) {
	t.Parallel()

	t.Run("test AssertContainNotLogRecord failure on error", func(t *testing.T) {
		t.Parallel()

		mt := new(testing.T)
		ce := fmt.Errorf("custom error")

		testLogBufferMock := new(TestLogBufferMock)
		testLogBufferMock.On("ContainRecord").Return(false, ce)
		logtest.AssertContainNotLogRecord(mt, testLogBufferMock, map[string]interface{}{})

		assert.True(t, mt.Failed())
	})

	t.Run("test AssertContainNotLogRecord failure on missing attribute", func(t *testing.T) {
		t.Parallel()

		mt := new(testing.T)

		testLogBufferMock := new(TestLogBufferMock)
		testLogBufferMock.On("ContainRecord").Return(true, nil)
		logtest.AssertContainNotLogRecord(mt, testLogBufferMock, map[string]interface{}{})

		assert.True(t, mt.Failed())
	})
}
