package config_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestConfigModule(t *testing.T) {
	createTestConfig := func(t *testing.T) *config.Config {
		t.Helper()

		var cfg *config.Config

		fxtest.New(
			t,
			fx.NopLogger,
			config.ConfigModule,
			fx.Populate(&cfg),
		).RequireStart().RequireStop()

		return cfg
	}

	t.Run("without config files and without env", func(t *testing.T) {
		cfg := createTestConfig(t)

		assert.Equal(t, config.DefaultAppName, cfg.AppName())
		assert.Equal(t, "", cfg.AppEnv())
		assert.Equal(t, config.DefaultAppVersion, cfg.AppVersion())
		assert.Equal(t, "", cfg.AppDescription())
		assert.False(t, cfg.AppDebug())

		assert.False(t, cfg.IsDevEnv())
		assert.False(t, cfg.IsTestEnv())
		assert.False(t, cfg.IsProdEnv())

		assert.Equal(t, "", cfg.GetString("config.values.string_value"))
		assert.Equal(t, 0, cfg.GetInt("config.values.int_value"))
		assert.Equal(t, "", cfg.GetString("config.substitution"))
	})

	t.Run("with config files and without env", func(t *testing.T) {
		t.Setenv("APP_CONFIG_PATH", "testdata/config/valid")

		cfg := createTestConfig(t)

		assert.Equal(t, "default-app", cfg.AppName())
		assert.Equal(t, config.AppEnvDev, cfg.AppEnv())
		assert.Equal(t, "default-version", cfg.AppVersion())
		assert.Equal(t, "default-description", cfg.AppDescription())
		assert.False(t, cfg.AppDebug())

		assert.False(t, cfg.IsDevEnv())
		assert.False(t, cfg.IsTestEnv())
		assert.False(t, cfg.IsProdEnv())

		assert.Equal(t, "default", cfg.GetString("config.values.string_value"))
		assert.Equal(t, 0, cfg.GetInt("config.values.int_value"))
		assert.Equal(t, "foo--baz", cfg.GetString("config.placeholder"))
		assert.Equal(t, "foo", cfg.GetString("config.substitution"))
	})

	t.Run("with config files and with test env", func(t *testing.T) {
		t.Setenv("APP_ENV", "test")
		t.Setenv("APP_CONFIG_PATH", "testdata/config/valid")

		cfg := createTestConfig(t)

		assert.Equal(t, "test-app", cfg.AppName())
		assert.Equal(t, config.AppEnvTest, cfg.AppEnv())
		assert.Equal(t, "default-version", cfg.AppVersion())
		assert.Equal(t, "default-description", cfg.AppDescription())
		assert.False(t, cfg.AppDebug())

		assert.False(t, cfg.IsDevEnv())
		assert.True(t, cfg.IsTestEnv())
		assert.False(t, cfg.IsProdEnv())

		assert.Equal(t, "default", cfg.GetString("config.values.string_value"))
		assert.Equal(t, 0, cfg.GetInt("config.values.int_value"))
		assert.Equal(t, "foo--baz", cfg.GetString("config.placeholder"))
		assert.Equal(t, "foo", cfg.GetString("config.substitution"))
	})
}
