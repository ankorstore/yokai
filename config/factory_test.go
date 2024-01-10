package config_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfigFactory(t *testing.T) {
	factory := config.NewDefaultConfigFactory()

	assert.IsType(t, &config.DefaultConfigFactory{}, factory)
	assert.Implements(t, (*config.ConfigFactory)(nil), factory)
}

func TestCreateSuccess(t *testing.T) {
	factory := config.NewDefaultConfigFactory()

	cfg, err := factory.Create(config.WithFilePaths("./testdata/config/valid"))

	assert.Nil(t, err)
	assert.IsType(t, &config.Config{}, cfg)

	assert.Equal(t, "default-app", cfg.AppName())
	assert.Equal(t, config.AppEnvDev, cfg.AppEnv())
	assert.Equal(t, false, cfg.AppDebug())
	assert.Equal(t, "0.1.0", cfg.AppVersion())
}

func TestCreateFailureOnInvalidConfigEnv(t *testing.T) {
	factory := config.NewDefaultConfigFactory()

	t.Setenv("APP_ENV", "invalid")

	_, err := factory.Create(config.WithFilePaths("./testdata/config/valid"))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not load config file for env invalid")
}

func TestCreateFailureOnInvalidConfigPath(t *testing.T) {
	factory := config.NewDefaultConfigFactory()

	_, err := factory.Create(config.WithFilePaths("./invalid-path"))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), `Config File "config" Not Found`)
}

func TestCreateFailureOnInvalidConfigContent(t *testing.T) {
	factory := config.NewDefaultConfigFactory()

	t.Setenv("APP_ENV", "test")

	_, err := factory.Create(config.WithFilePaths("./testdata/config/invalid"))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not merge config for env test")
}
