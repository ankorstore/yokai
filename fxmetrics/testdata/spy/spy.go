package spy

import (
	"bytes"
	"fmt"
)

type SpyTB struct {
	failures int
	errors   *bytes.Buffer
	logs     *bytes.Buffer
}

func NewSpyTB() *SpyTB {
	return &SpyTB{0, &bytes.Buffer{}, &bytes.Buffer{}}
}

func (t *SpyTB) Failures() int {
	return t.failures
}

func (t *SpyTB) Errors() *bytes.Buffer {
	return t.errors
}

func (t *SpyTB) Logs() *bytes.Buffer {
	return t.logs
}

func (t *SpyTB) FailNow() {
	t.failures++
}

func (t *SpyTB) Errorf(format string, args ...interface{}) {
	fmt.Fprintf(t.errors, format, args...)
	t.errors.WriteRune('\n')
}

func (t *SpyTB) Logf(format string, args ...interface{}) {
	fmt.Fprintf(t.logs, format, args...)
	t.logs.WriteRune('\n')
}
