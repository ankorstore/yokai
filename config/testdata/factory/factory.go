package factory

import (
	"github.com/ankorstore/yokai/modules/config"
)

type TestConfigFactory struct{}

func NewTestConfigFactory() config.ConfigFactory {
	return &TestConfigFactory{}
}

func (f *TestConfigFactory) Create(options ...config.ConfigOption) (*config.Config, error) {
	return &config.Config{}, nil
}
