package logtest_test

import (
	"bytes"
	"testing"

	"github.com/ankorstore/yokai/log/logtest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestLogBufferMock struct {
	mock.Mock
}

func (m *TestLogBufferMock) Buffer() *bytes.Buffer {
	args := m.Called()

	//nolint:forcetypeassert
	return args.Get(0).(*bytes.Buffer)
}

func (m *TestLogBufferMock) Write(p []byte) (int, error) {
	args := m.Called(p)

	return args.Int(0), args.Error(1)
}

func (m *TestLogBufferMock) Reset() logtest.TestLogBuffer {
	args := m.Called()

	//nolint:forcetypeassert
	return args.Get(0).(logtest.TestLogBuffer)
}

func (m *TestLogBufferMock) Records() ([]*logtest.TestLogRecord, error) {
	args := m.Called()

	//nolint:forcetypeassert
	return args.Get(0).([]*logtest.TestLogRecord), args.Error(1)
}

func (m *TestLogBufferMock) HasRecord(_ map[string]interface{}) (bool, error) {
	args := m.Called()

	//nolint:forcetypeassert
	return args.Get(0).(bool), args.Error(1)
}

func (m *TestLogBufferMock) ContainRecord(_ map[string]interface{}) (bool, error) {
	args := m.Called()

	//nolint:forcetypeassert
	return args.Get(0).(bool), args.Error(1)
}

func TestTestLogBuffer(t *testing.T) {
	t.Parallel()

	t.Run("test Buffer and ClearRecords", func(t *testing.T) {
		t.Parallel()

		buffer := logtest.NewDefaultTestLogBuffer()
		buffer.Reset()

		buffer.Buffer().WriteString("test")
		assert.NotEqual(t, 0, buffer.Buffer().Len())

		buffer.Reset()
		assert.Equal(t, 0, buffer.Buffer().Len())
	})

	t.Run("test Records and HasRecord", func(t *testing.T) {
		t.Parallel()

		buffer := logtest.NewDefaultTestLogBuffer()
		buffer.Reset()

		_, err := buffer.Write([]byte("{\"level\":\"info\",\"message\":\"first\",\"service\":\"test\",\"int\":1,\"float\":1.5,\"bool\":false,\"time\":1698312453}\n"))
		assert.NoError(t, err)

		_, err = buffer.Write([]byte("{\"level\":\"info\",\"message\":\"second\",\"service\":\"test\",\"int\":2,\"float\":2.5,\"bool\":true,\"time\":1698312454}\n"))
		assert.NoError(t, err)

		records, err := buffer.Records()
		assert.NoError(t, err)
		assert.Len(t, records, 2)

		match, err := buffer.HasRecord(map[string]interface{}{"message": "first", "int": 1, "float": 1.5, "bool": false, "time": 1698312453})
		assert.NoError(t, err)
		assert.True(t, match)

		match, err = buffer.HasRecord(map[string]interface{}{"message": "first", "int": "1", "float": "1.5", "bool": "false", "time": "1698312453"})
		assert.NoError(t, err)
		assert.True(t, match)

		match, err = buffer.HasRecord(map[string]interface{}{"message": "first", "int": "1", "float": "1.5", "bool": true, "time": "1698312453"})
		assert.NoError(t, err)
		assert.False(t, match)

		match, err = buffer.HasRecord(map[string]interface{}{"message": "second", "int": 2, "float": 2.5, "bool": true, "time": 1698312454})
		assert.NoError(t, err)
		assert.True(t, match)

		match, err = buffer.HasRecord(map[string]interface{}{"message": "second", "int": "2", "float": "2.5", "bool": "true", "time": "1698312454"})
		assert.NoError(t, err)
		assert.True(t, match)

		match, err = buffer.HasRecord(map[string]interface{}{"message": "second", "int": "2", "float": 1.5, "bool": "true", "time": "1698312454"})
		assert.NoError(t, err)
		assert.False(t, match)
	})

	t.Run("test Records and ContainRecord", func(t *testing.T) {
		t.Parallel()

		buffer := logtest.NewDefaultTestLogBuffer()
		buffer.Reset()

		_, err := buffer.Write([]byte("{\"level\":\"info\",\"message\":\"first\",\"service\":\"test\",\"int\":1,\"float\":1.5,\"bool\":false,\"time\":1698312453}\n"))
		assert.NoError(t, err)

		_, err = buffer.Write([]byte("{\"level\":\"info\",\"message\":\"second\",\"service\":\"test\",\"int\":2,\"float\":2.5,\"bool\":true,\"time\":1698312454}\n"))
		assert.NoError(t, err)

		records, err := buffer.Records()
		assert.NoError(t, err)
		assert.Len(t, records, 2)

		match, err := buffer.ContainRecord(map[string]interface{}{"message": "first", "int": 1, "float": 1.5, "bool": false, "time": 1698312453})
		assert.NoError(t, err)
		assert.True(t, match)

		match, err = buffer.ContainRecord(map[string]interface{}{"message": "first", "int": "1", "float": "1.5", "bool": "false", "time": "1698312453"})
		assert.NoError(t, err)
		assert.True(t, match)

		match, err = buffer.ContainRecord(map[string]interface{}{"message": "fir", "int": "1", "float": "1.5", "bool": "false", "time": "1698312453"})
		assert.NoError(t, err)
		assert.True(t, match)

		match, err = buffer.ContainRecord(map[string]interface{}{"message": "first", "int": "1", "float": "1.5", "bool": true, "time": "1698312453"})
		assert.NoError(t, err)
		assert.False(t, match)

		match, err = buffer.ContainRecord(map[string]interface{}{"message": "second", "int": 2, "float": 2.5, "bool": true, "time": 1698312454})
		assert.NoError(t, err)
		assert.True(t, match)

		match, err = buffer.ContainRecord(map[string]interface{}{"message": "second", "int": "2", "float": "2.5", "bool": "true", "time": "1698312454"})
		assert.NoError(t, err)
		assert.True(t, match)

		match, err = buffer.ContainRecord(map[string]interface{}{"message": "cond", "int": "2", "float": "2.5", "bool": "true", "time": "1698312454"})
		assert.NoError(t, err)
		assert.True(t, match)

		match, err = buffer.ContainRecord(map[string]interface{}{"message": "second", "int": "2", "float": 1.5, "bool": "true", "time": "1698312454"})
		assert.NoError(t, err)
		assert.False(t, match)
	})

	t.Run("test Records and HasRecord error on invalid log buffer content", func(t *testing.T) {
		t.Parallel()

		buffer := logtest.NewDefaultTestLogBuffer()
		buffer.Reset()

		_, err := buffer.Write([]byte("{{\n"))
		assert.NoError(t, err)

		_, err = buffer.Records()
		assert.Error(t, err)

		_, err = buffer.HasRecord(map[string]interface{}{"some": "value"})
		assert.Error(t, err)
	})

	t.Run("test Records and ContainRecord error on invalid log buffer content", func(t *testing.T) {
		t.Parallel()

		buffer := logtest.NewDefaultTestLogBuffer()
		buffer.Reset()

		_, err := buffer.Write([]byte("{{\n"))
		assert.NoError(t, err)

		_, err = buffer.Records()
		assert.Error(t, err)

		_, err = buffer.ContainRecord(map[string]interface{}{"some": "value"})
		assert.Error(t, err)
	})
}
