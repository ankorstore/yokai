package sql

// Configuration is the SQL components (driver, connector, etc) configuration.
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

// System returns the Configuration system.
func (c *Configuration) System() System {
	return c.system
}

// Hooks returns the Configuration list of hook.Hook.
func (c *Configuration) Hooks() []Hook {
	return c.hooks
}
