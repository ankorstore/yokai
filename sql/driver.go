package sql

import (
	"database/sql/driver"
)

type Driver struct {
	base          driver.Driver
	configuration *Configuration
}

func NewDriver(base driver.Driver, configuration *Configuration) *Driver {
	return &Driver{
		base:          base,
		configuration: configuration,
	}
}

func (d *Driver) Configuration() *Configuration {
	return d.configuration
}

func (d *Driver) Open(dsn string) (driver.Conn, error) {
	connection, err := d.base.Open(dsn)
	if err != nil {
		return nil, err
	}

	return NewConnection(connection, d.configuration), nil
}

func (d *Driver) OpenConnector(dsn string) (driver.Connector, error) {
	if driverContext, ok := d.base.(driver.DriverContext); ok {
		connector, err := driverContext.OpenConnector(dsn)
		if err != nil {
			return nil, err
		}

		return NewConnector(dsn, connector, d), nil
	}

	return NewConnector(dsn, nil, d), nil
}
