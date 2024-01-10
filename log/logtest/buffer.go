package logtest

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io"
	"strings"
	"sync"
)

// TestLogBuffer is the interface for test log buffers.
type TestLogBuffer interface {
	io.Writer
	Buffer() *bytes.Buffer
	Reset() TestLogBuffer
	Records() ([]*TestLogRecord, error)
	HasRecord(expectedAttributes map[string]interface{}) (bool, error)
	ContainRecord(expectedAttributes map[string]interface{}) (bool, error)
}

// DefaultTestLogBuffer is the default [TestLogBuffer] implementation.
type DefaultTestLogBuffer struct {
	buffer *bytes.Buffer
	lock   *sync.Mutex
}

// NewDefaultTestLogBuffer returns a [DefaultTestLogBuffer], implementing [TestLogBuffer].
func NewDefaultTestLogBuffer() TestLogBuffer {
	return &DefaultTestLogBuffer{
		buffer: &bytes.Buffer{},
		lock:   &sync.Mutex{},
	}
}

// Buffer returns the internal buffer.
func (b *DefaultTestLogBuffer) Buffer() *bytes.Buffer {
	return b.buffer
}

// Write writes into the internal buffer.
func (b *DefaultTestLogBuffer) Write(p []byte) (int, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	return b.buffer.Write(p)
}

// Reset resets the internal buffer.
func (b *DefaultTestLogBuffer) Reset() TestLogBuffer {
	b.lock.Lock()
	defer b.lock.Unlock()

	b.buffer.Reset()

	return b
}

// Records return the list of [TestLogRecord] from the internal buffer.
func (b *DefaultTestLogBuffer) Records() ([]*TestLogRecord, error) {
	b.lock.Lock()
	defer b.lock.Unlock()

	var records []*TestLogRecord

	clone := *b.buffer
	scanner := bufio.NewScanner(&clone)

	for scanner.Scan() {
		var attributes map[string]interface{}

		decoder := json.NewDecoder(strings.NewReader(scanner.Text()))
		decoder.UseNumber()

		if err := decoder.Decode(&attributes); err != nil {
			return nil, err
		}

		records = append(records, NewTestLogRecord(attributes))
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

// HasRecord return true if a log record from the internal buffer is exactly matching provided attributes.
func (b *DefaultTestLogBuffer) HasRecord(expectedAttributes map[string]interface{}) (bool, error) {
	records, err := b.Records()
	if err != nil {
		return false, err
	}

	for _, record := range records {
		if record.MatchAttributes(expectedAttributes) {
			return true, nil
		}
	}

	return false, nil
}

// ContainRecord return true if a log record from the internal buffer is partially matching provided attributes.
func (b *DefaultTestLogBuffer) ContainRecord(expectedAttributes map[string]interface{}) (bool, error) {
	records, err := b.Records()
	if err != nil {
		return false, err
	}

	for _, record := range records {
		if record.ContainAttributes(expectedAttributes) {
			return true, nil
		}
	}

	return false, nil
}
