package sql

// Configuration is the SQL components (driver, connector, connection, etc) configuration.
type Configuration struct {
	system System
	hooks  []Hook
}

// NewConfiguration returns a new Configuration.
func NewConfiguration(system System, hooks ...Hook) *Configuration {
	return &Configuration{
		system: system,
		hooks:  hooks,
	}
}

// System returns the Configuration System.
func (c *Configuration) System() System {
	return c.system
}

// Hooks returns the Configuration list of Hook.
func (c *Configuration) Hooks() []Hook {
	return c.hooks
}
