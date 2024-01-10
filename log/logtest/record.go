package logtest

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ankorstore/yokai/log"
)

// TestLogRecord is a test log record, composed of attributes.
type TestLogRecord struct {
	attributes map[string]interface{}
}

// NewTestLogRecord returns a [TestLogRecord].
func NewTestLogRecord(attributes map[string]interface{}) *TestLogRecord {
	return &TestLogRecord{
		attributes: attributes,
	}
}

// Level returns the level of the [TestLogRecord].
func (r *TestLogRecord) Level() (string, error) {
	value, err := r.Attribute(log.Level)
	if err != nil {
		return "", err
	}

	if level, ok := value.(string); ok {
		return level, nil
	} else {
		return "", fmt.Errorf("cannot cast level as string")
	}
}

// Message returns the message of the [TestLogRecord].
func (r *TestLogRecord) Message() (string, error) {
	value, err := r.Attribute(log.Message)
	if err != nil {
		return "", err
	}

	if message, ok := value.(string); ok {
		return message, nil
	} else {
		return "", fmt.Errorf("cannot cast message as string")
	}
}

// Service returns the service name of the [TestLogRecord].
func (r *TestLogRecord) Service() (string, error) {
	value, err := r.Attribute(log.Service)
	if err != nil {
		return "", err
	}

	if service, ok := value.(string); ok {
		return service, nil
	} else {
		return "", fmt.Errorf("cannot cast service as string")
	}
}

// Time returns the time of the [TestLogRecord].
func (r *TestLogRecord) Time() (time.Time, error) {
	value, err := r.Attribute(log.Time)
	if err != nil {
		return time.Unix(0, 0), err
	}

	if t, ok := value.(int64); ok {
		return time.Unix(t, 0), nil
	} else {
		return time.Unix(0, 0), fmt.Errorf("cannot cast time as time.Time")
	}
}

// Attribute returns an attribute of the [TestLogRecord] by name.
func (r *TestLogRecord) Attribute(name string) (interface{}, error) {
	value, ok := r.attributes[name]
	if ok {
		return value, nil
	} else {
		return "", fmt.Errorf("attribute %s not found in record", name)
	}
}

// MatchAttributes returns true if the [TestLogRecord] exactly matches provided attributes.
//
//nolint:cyclop
func (r *TestLogRecord) MatchAttributes(expectedAttributes map[string]interface{}) bool {
	match := true

	if len(expectedAttributes) == 0 {
		return false
	}

	for expectedName, expectedValue := range expectedAttributes {
		value, ok := r.attributes[expectedName]
		if ok {
			switch av := value.(type) {
			case json.Number:
				switch ev := expectedValue.(type) {
				case int:
					avi, err := av.Int64()
					if err != nil {
						match = false

						break
					}
					match = match && int(avi) == ev
				case float64:
					avf, err := av.Float64()
					if err != nil {
						match = false

						break
					}
					match = match && avf == ev
				default:
					match = match && av.String() == fmt.Sprintf("%s", ev)
				}
			case bool:
				match = match && fmt.Sprintf("%v", av) == fmt.Sprintf("%v", expectedValue)
			default:
				match = match && fmt.Sprintf("%s", av) == fmt.Sprintf("%s", expectedValue)
			}
		} else {
			match = false

			break
		}
	}

	return match
}

// ContainAttributes returns true if the [TestLogRecord] partially matches provided attributes.
//
//nolint:cyclop
func (r *TestLogRecord) ContainAttributes(expectedAttributes map[string]interface{}) bool {
	match := true

	if len(expectedAttributes) == 0 {
		return false
	}

	for expectedName, expectedValue := range expectedAttributes {
		value, ok := r.attributes[expectedName]
		if ok {
			switch av := value.(type) {
			case json.Number:
				switch ev := expectedValue.(type) {
				case int:
					avi, err := av.Int64()
					if err != nil {
						match = false

						break
					}
					match = match && int(avi) == ev
				case float64:
					avf, err := av.Float64()
					if err != nil {
						match = false

						break
					}
					match = match && avf == ev
				default:
					match = match && strings.Contains(av.String(), fmt.Sprintf("%s", ev))
				}
			case bool:
				match = match && fmt.Sprintf("%v", av) == fmt.Sprintf("%v", expectedValue)
			default:
				match = match && strings.Contains(fmt.Sprintf("%s", av), fmt.Sprintf("%s", expectedValue))
			}
		} else {
			match = false

			break
		}
	}

	return match
}
