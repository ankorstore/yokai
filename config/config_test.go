package config_test

import (
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/stretchr/testify/assert"
)

func TestAppNameFromDefaultConfig(t *testing.T) {
	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "default-app", cfg.AppName())
}

func TestAppNameOverrideFromTestEnvConfig(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "test-app", cfg.AppName())
}

func TestAppNameOverrideFromCustomEnvConfig(t *testing.T) {
	t.Setenv("APP_ENV", "custom")

	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "custom-app", cfg.AppName())
}

func TestAppNameOverrideFromInvalidEnvConfig(t *testing.T) {
	t.Setenv("APP_ENV", "invalid")

	_, err := createTestConfig()

	assert.Error(t, err)
	assert.Contains(t, err.Error(), `could not load config file for env invalid: Config File "config.invalid" Not Found`)
}

func TestAppNameOverrideFromEnvVar(t *testing.T) {
	t.Setenv("APP_NAME", "env-app")

	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "env-app", cfg.AppName())
}

func TestAppEnvFromConfig(t *testing.T) {
	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, config.AppEnvDev, cfg.AppEnv())
	assert.False(t, cfg.IsProdEnv())
	assert.True(t, cfg.IsDevEnv())
	assert.False(t, cfg.IsTestEnv())
}

func TestAppEnvOverrideFromTestEnvConfig(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, config.AppEnvTest, cfg.AppEnv())
	assert.False(t, cfg.IsProdEnv())
	assert.False(t, cfg.IsDevEnv())
	assert.True(t, cfg.IsTestEnv())
}

func TestAppEnvOverrideFromCustomEnvConfig(t *testing.T) {
	t.Setenv("APP_ENV", "custom")

	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "custom", cfg.AppEnv())
	assert.False(t, cfg.IsProdEnv())
	assert.False(t, cfg.IsDevEnv())
	assert.False(t, cfg.IsTestEnv())
}

func TestAppDebugFromConfig(t *testing.T) {
	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.False(t, cfg.AppDebug())
}

func TestAppDebugOverrideFromTestEnvConfig(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.True(t, cfg.AppDebug())
}

func TestAppDebugOverrideFromEnvVar(t *testing.T) {
	t.Setenv("APP_DEBUG", "true")

	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.True(t, cfg.AppDebug())
}

func TestAppVersionFromConfig(t *testing.T) {
	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "0.1.0", cfg.AppVersion())
}

func TestAppVersionOverrideFromTestEnvConfig(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "0.1.0", cfg.AppVersion())
}

func TestAppVersionOverrideFromEnvVar(t *testing.T) {
	t.Setenv("APP_VERSION", "0.1.2")

	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "0.1.2", cfg.AppVersion())
}

func TestValuesFromConfig(t *testing.T) {
	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "default", cfg.GetString("config.values.string_value"))
	assert.Equal(t, 0, cfg.GetInt("config.values.int_value"))
}

func TestValuesOverrideFromTestEnvConfig(t *testing.T) {
	t.Setenv("APP_ENV", "test")

	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "test", cfg.GetString("config.values.string_value"))
	assert.Equal(t, 0, cfg.GetInt("config.values.int_value"))
}

func TestValuesOverrideFromCustomEnvConfig(t *testing.T) {
	t.Setenv("APP_ENV", "custom")

	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "custom", cfg.GetString("config.values.string_value"))
	assert.Equal(t, 0, cfg.GetInt("config.values.int_value"))
}

func TestValuesWithEnvVarsPlaceholder(t *testing.T) {
	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "foo--baz", cfg.GetString("config.placeholder"))

	t.Setenv("BAR", "bar")

	cfg, err = createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "foo-bar-baz", cfg.GetString("config.placeholder"))
}

func TestValuesWithEnvVarsSubstitution(t *testing.T) {
	cfg, err := createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "foo", cfg.GetString("config.substitution"))

	t.Setenv("CONFIG_SUBSTITUTION", "bar")

	cfg, err = createTestConfig()

	assert.Nil(t, err)
	assert.Equal(t, "bar", cfg.GetString("config.substitution"))
}

func createTestConfig() (*config.Config, error) {
	return config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("./testdata/config/valid"),
	)
}
