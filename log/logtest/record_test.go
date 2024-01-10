package logtest_test

import (
	"testing"
	"time"

	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/stretchr/testify/assert"
)

func TestTestLogRecord(t *testing.T) {
	t.Parallel()

	t.Run("test common attributes success", func(t *testing.T) {
		t.Parallel()

		attributes := map[string]interface{}{
			log.Level:   "info",
			log.Message: "Hello, world!",
			log.Service: "test",
			log.Time:    time.Now().Unix(),
		}

		record := logtest.NewTestLogRecord(attributes)

		level, err := record.Level()
		assert.NoError(t, err)
		assert.Equal(t, "info", level)

		message, err := record.Message()
		assert.NoError(t, err)
		assert.Equal(t, "Hello, world!", message)

		service, err := record.Service()
		assert.NoError(t, err)
		assert.Equal(t, "test", service)

		logTime, err := record.Time()
		assert.NoError(t, err)
		//nolint:forcetypeassert
		assert.Equal(t, time.Unix(attributes[log.Time].(int64), 0), logTime)
	})

	t.Run("test common attributes failure on missing attributes", func(t *testing.T) {
		t.Parallel()
		emptyRecord := logtest.NewTestLogRecord(map[string]interface{}{})

		_, err := emptyRecord.Level()
		assert.Error(t, err)
		assert.Equal(t, "attribute level not found in record", err.Error())

		_, err = emptyRecord.Message()
		assert.Error(t, err)
		assert.Equal(t, "attribute message not found in record", err.Error())

		_, err = emptyRecord.Service()
		assert.Error(t, err)
		assert.Equal(t, "attribute service not found in record", err.Error())

		_, err = emptyRecord.Time()
		assert.Error(t, err)
		assert.Equal(t, "attribute time not found in record", err.Error())
	})

	t.Run("test common attributes failure on invalid attributes", func(t *testing.T) {
		t.Parallel()
		record := logtest.NewTestLogRecord(map[string]interface{}{
			"level":   0,
			"message": 0,
			"service": 0,
			"time":    "invalid",
		})

		_, err := record.Level()
		assert.Error(t, err)
		assert.Equal(t, "cannot cast level as string", err.Error())

		_, err = record.Message()
		assert.Error(t, err)
		assert.Equal(t, "cannot cast message as string", err.Error())

		_, err = record.Service()
		assert.Error(t, err)
		assert.Equal(t, "cannot cast service as string", err.Error())

		_, err = record.Time()
		assert.Error(t, err)
		assert.Equal(t, "cannot cast time as time.Time", err.Error())
	})

	t.Run("test Attribute", func(t *testing.T) {
		t.Parallel()

		attributes := map[string]interface{}{
			"attribute": "value",
		}
		record := logtest.NewTestLogRecord(attributes)

		value, err := record.Attribute("attribute")
		assert.NoError(t, err)
		assert.Equal(t, "value", value)

		_, err = record.Attribute("nonexistent")
		assert.Error(t, err)
	})

	t.Run("test MatchAttributes with attributes", func(t *testing.T) {
		t.Parallel()

		attributes := map[string]interface{}{
			"attribute": "value",
		}
		record := logtest.NewTestLogRecord(attributes)

		match := record.MatchAttributes(attributes)
		assert.True(t, match)

		nonMatchingAttributes := map[string]interface{}{
			"attribute": "otherValue",
		}
		match = record.MatchAttributes(nonMatchingAttributes)
		assert.False(t, match)

		nonExistentAttributes := map[string]interface{}{
			"nonexistent": "value",
		}
		match = record.MatchAttributes(nonExistentAttributes)
		assert.False(t, match)
	})

	t.Run("test ContainAttributes with attributes", func(t *testing.T) {
		t.Parallel()

		attributes := map[string]interface{}{
			"attribute": "some value",
		}
		record := logtest.NewTestLogRecord(attributes)

		match := record.ContainAttributes(attributes)
		assert.True(t, match)

		containedAttributes := map[string]interface{}{
			"attribute": "some v",
		}
		match = record.ContainAttributes(containedAttributes)
		assert.True(t, match)

		otherContainedAttributes := map[string]interface{}{
			"attribute": "me val",
		}
		match = record.ContainAttributes(otherContainedAttributes)
		assert.True(t, match)

		otherAgainContainedAttributes := map[string]interface{}{
			"attribute": "lue",
		}
		match = record.ContainAttributes(otherAgainContainedAttributes)
		assert.True(t, match)

		nonContainedAttributes := map[string]interface{}{
			"attribute": "otherValue",
		}
		match = record.ContainAttributes(nonContainedAttributes)
		assert.False(t, match)

		nonExistentAttributes := map[string]interface{}{
			"nonexistent": "value",
		}
		match = record.ContainAttributes(nonExistentAttributes)
		assert.False(t, match)
	})

	t.Run("test MatchAttributes without attributes", func(t *testing.T) {
		t.Parallel()

		record := logtest.NewTestLogRecord(map[string]interface{}{})
		match := record.MatchAttributes(map[string]interface{}{})
		assert.False(t, match)
	})

	t.Run("test ContainAttributes without attributes", func(t *testing.T) {
		t.Parallel()

		record := logtest.NewTestLogRecord(map[string]interface{}{})
		match := record.ContainAttributes(map[string]interface{}{})
		assert.False(t, match)
	})
}
