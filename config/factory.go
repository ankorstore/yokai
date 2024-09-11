package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// ConfigFactory is the interface for [Config] factories.
type ConfigFactory interface {
	Create(options ...ConfigOption) (*Config, error)
}

// DefaultConfigFactory is the default [ConfigFactory] implementation.
type DefaultConfigFactory struct{}

// NewDefaultConfigFactory returns a new [DefaultConfigFactory] instance.
func NewDefaultConfigFactory() *DefaultConfigFactory {
	return &DefaultConfigFactory{}
}

// Create returns a new [Config], and accepts a list of [ConfigOption].
// For example:
//
//	var cfg, _ = config.NewDefaultConfigFactory().Create()
//
// is equivalent to:
//
//	var cfg, _ = config.NewDefaultConfigFactory().Create(
//		config.WithFileName("config"),                      // config files base name
//		config.WithFilePaths(".", "./config", "./configs"), // config files lookup paths
//	)
func (f *DefaultConfigFactory) Create(options ...ConfigOption) (*Config, error) {
	// options
	appliedOptions := DefaultConfigOptions()
	for _, opt := range options {
		opt(&appliedOptions)
	}

	// viper
	v := viper.New()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	v.SetConfigName(appliedOptions.FileName)

	for _, path := range appliedOptions.FilePaths {
		v.AddConfigPath(path)
	}

	// defaults
	v.SetDefault("app.name", DefaultAppName)
	v.SetDefault("app.version", DefaultAppVersion)
	v.SetDefault("app.debug", false)

	// load
	if err := v.ReadInConfig(); err != nil {
		if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
			return nil, err
		}
	}

	// env overrides
	appEnv := os.Getenv("APP_ENV")
	if appEnv != "" {
		v.SetConfigName(fmt.Sprintf("%s.%s", appliedOptions.FileName, appEnv))
		if err := v.MergeInConfig(); err != nil {
			if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
				return nil, fmt.Errorf("could not merge config for env %s: %w", appEnv, err)
			}
		}
	}

	// env vars placeholders
	for _, key := range v.AllKeys() {
		val := v.GetString(key)
		if strings.Contains(val, "${") {
			v.Set(key, os.ExpandEnv(val))
		}
	}

	return NewConfig(v), nil
}
