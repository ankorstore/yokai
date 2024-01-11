package fxconfig_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxconfig/testdata/factory"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestModule(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	var cfg *config.Config

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fx.Populate(&cfg),
	).RequireStart().RequireStop()

	assert.Equal(t, "default-app", cfg.AppName())
	assert.Equal(t, config.AppEnvDev, cfg.AppEnv())
	assert.True(t, cfg.IsDevEnv())
	assert.False(t, cfg.IsTestEnv())
	assert.False(t, cfg.IsProdEnv())
	assert.False(t, cfg.AppDebug())
	assert.Equal(t, "0.1.0", cfg.AppVersion())

	assert.Equal(t, "default", cfg.GetString("config.values.string_value"))
	assert.Equal(t, 0, cfg.GetInt("config.values.int_value"))

	assert.Equal(t, "foo--baz", cfg.GetString("config.substitution"))
}

func TestModuleWithTestEnv(t *testing.T) {
	t.Setenv("APP_ENV", "test")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("BAR", "bar")

	var cfg *config.Config

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fx.Populate(&cfg),
	).RequireStart().RequireStop()

	assert.Equal(t, "test-app", cfg.AppName())
	assert.Equal(t, config.AppEnvTest, cfg.AppEnv())
	assert.False(t, cfg.IsDevEnv())
	assert.True(t, cfg.IsTestEnv())
	assert.False(t, cfg.IsProdEnv())
	assert.True(t, cfg.AppDebug())
	assert.Equal(t, "0.1.0", cfg.AppVersion())

	assert.Equal(t, "test", cfg.GetString("config.values.string_value"))
	assert.Equal(t, 0, cfg.GetInt("config.values.int_value"))

	assert.Equal(t, "foo-bar-baz", cfg.GetString("config.substitution"))
}

func TestModuleWithCustomEnv(t *testing.T) {
	t.Setenv("APP_ENV", "custom")
	t.Setenv("APP_CONFIG_PATH", "testdata/config")
	t.Setenv("BAR", "bar")

	var cfg *config.Config

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fx.Populate(&cfg),
	).RequireStart().RequireStop()

	assert.Equal(t, "custom-app", cfg.AppName())
	assert.Equal(t, "custom", cfg.AppEnv())
	assert.False(t, cfg.IsDevEnv())
	assert.False(t, cfg.IsTestEnv())
	assert.False(t, cfg.IsProdEnv())
	assert.True(t, cfg.AppDebug())
	assert.Equal(t, "0.1.0", cfg.AppVersion())

	assert.Equal(t, "custom", cfg.GetString("config.values.string_value"))
	assert.Equal(t, 0, cfg.GetInt("config.values.int_value"))

	assert.Equal(t, "foo-bar-baz", cfg.GetString("config.substitution"))
}

func TestModuleDecoration(t *testing.T) {
	var cfg *config.Config

	fxtest.New(
		t,
		fx.NopLogger,
		fxconfig.FxConfigModule,
		fx.Decorate(factory.NewTestConfigFactory),
		fx.Populate(&cfg),
	).RequireStart().RequireStop()

	assert.Equal(t, &config.Config{}, cfg)
}
