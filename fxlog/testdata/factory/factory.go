package factory

import (
	"github.com/ankorstore/yokai/log"
)

type TestLoggerFactory struct{}

func NewTestLoggerFactory() log.LoggerFactory {
	return &TestLoggerFactory{}
}

func (f *TestLoggerFactory) Create(options ...log.LoggerOption) (*log.Logger, error) {
	return &log.Logger{}, nil
}
