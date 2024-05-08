package sql

import (
	"database/sql/driver"
)

var (
	_ driver.Driver        = (*Driver)(nil)
	_ driver.DriverContext = (*Driver)(nil)
)

// Driver is a SQL driver wrapping a driver.Driver.
type Driver struct {
	base          driver.Driver
	configuration *Configuration
}

// NewDriver returns a new Driver.
func NewDriver(base driver.Driver, configuration *Configuration) *Driver {
	return &Driver{
		base:          base,
		configuration: configuration,
	}
}

// Base returns the base driver.Driver of the Driver.
func (d *Driver) Base() driver.Driver {
	return d.base
}

// Configuration returns the Configuration of the Driver.
func (d *Driver) Configuration() *Configuration {
	return d.configuration
}

// Open returns a new Connection.
func (d *Driver) Open(dsn string) (driver.Conn, error) {
	connection, err := d.base.Open(dsn)
	if err != nil {
		return nil, err
	}

	return NewConnection(connection, d.configuration), nil
}

// OpenConnector returns a new Connector.
func (d *Driver) OpenConnector(dsn string) (driver.Connector, error) {
	if engine, ok := d.base.(driver.DriverContext); ok {
		connector, err := engine.OpenConnector(dsn)
		if err != nil {
			return nil, err
		}

		return NewConnector(dsn, connector, d), nil
	}

	return NewConnector(dsn, nil, d), nil
}
