package logtest

import "testing"

// AssertHasLogRecord allows to assert if a log record exactly matching provided attributes can be found.
func AssertHasLogRecord(tb testing.TB, testLogBuffer TestLogBuffer, expectedAttributes map[string]interface{}) bool {
	tb.Helper()

	hasRecord, err := testLogBuffer.HasRecord(expectedAttributes)
	if err != nil {
		tb.Errorf("error while asserting log record attributes match: %v", err)

		return false
	}

	if !hasRecord {
		tb.Errorf("cannot find log record with matching attributes %+v", expectedAttributes)

		return false
	}

	return true
}

// AssertHasNotLogRecord allows to assert if a log record exactly matching provided attributes cannot be found.
func AssertHasNotLogRecord(tb testing.TB, testLogBuffer TestLogBuffer, expectedAttributes map[string]interface{}) bool {
	tb.Helper()

	hasRecord, err := testLogBuffer.HasRecord(expectedAttributes)
	if err != nil {
		tb.Errorf("error while asserting log record attributes match: %v", err)

		return false
	}

	if hasRecord {
		tb.Errorf("can find log record with matching attributes %+v", expectedAttributes)

		return false
	}

	return true
}

// AssertContainLogRecord allows to assert if a log record partially matching provided attributes can be found.
func AssertContainLogRecord(tb testing.TB, testLogBuffer TestLogBuffer, expectedAttributes map[string]interface{}) bool {
	tb.Helper()

	containRecord, err := testLogBuffer.ContainRecord(expectedAttributes)
	if err != nil {
		tb.Errorf("error while asserting log record attributes contain: %v", err)

		return false
	}

	if !containRecord {
		tb.Errorf("cannot find log record with contained attributes %+v", expectedAttributes)

		return false
	}

	return true
}

// AssertContainNotLogRecord allows to assert if a log record partially matching provided attributes cannot be found.
func AssertContainNotLogRecord(tb testing.TB, testLogBuffer TestLogBuffer, expectedAttributes map[string]interface{}) bool {
	tb.Helper()

	containRecord, err := testLogBuffer.ContainRecord(expectedAttributes)
	if err != nil {
		tb.Errorf("error while asserting log record attributes contain: %v", err)

		return false
	}

	if containRecord {
		tb.Errorf("can find log record with contained attributes %+v", expectedAttributes)

		return false
	}

	return true
}
