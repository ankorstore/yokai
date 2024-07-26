package config

import (
	"os"

	"github.com/spf13/viper"
)

const (
	AppEnvProd = "prod" // prod environment
	AppEnvDev  = "dev"  // dev environment
	AppEnvTest = "test" // test environment

	DefaultAppName    = "app"     // default application name
	DefaultAppVersion = "unknown" // default application version
)

// Config allows to access the application configuration, and inherits of all [Viper] features.
//
// [Viper]: https://github.com/spf13/viper
type Config struct {
	*viper.Viper
}

// GetEnvVar returns the value of an env var.
func (c *Config) GetEnvVar(envVar string) string {
	return os.Getenv(envVar)
}

// AppName returns the configured application name (from config field app.name or env var APP_NAME).
func (c *Config) AppName() string {
	return c.GetString("app.name")
}

// AppDescription returns the configured application description (from config field app.description or env var APP_DESCRIPTION).
func (c *Config) AppDescription() string {
	return c.GetString("app.description")
}

// AppEnv returns the configured application environment (from config field app.env or env var APP_ENV).
func (c *Config) AppEnv() string {
	return c.GetString("app.env")
}

// AppVersion returns the configured application version (from config field app.version or env var APP_VERSION).
func (c *Config) AppVersion() string {
	return c.GetString("app.version")
}

// AppDebug returns if the application debug mode is enabled (from config field app.debug or env var APP_DEBUG).
func (c *Config) AppDebug() bool {
	return c.GetBool("app.debug")
}

// IsProdEnv returns if the application is running in prod mode.
func (c *Config) IsProdEnv() bool {
	return c.AppEnv() == AppEnvProd
}

// IsDevEnv returns if the application is running in dev mode.
func (c *Config) IsDevEnv() bool {
	return c.AppEnv() == AppEnvDev
}

// IsTestEnv returns if the application is running in test mode.
func (c *Config) IsTestEnv() bool {
	return c.AppEnv() == AppEnvTest
}
