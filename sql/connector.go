package sql

import (
	"context"
	"database/sql/driver"
)

var _ driver.Connector = (*Connector)(nil)

// Connector is a SQL driver connector wrapping a driver.Connector.
type Connector struct {
	dsn    string
	base   driver.Connector
	driver *Driver
}

// NewConnector returns a new Connector.
func NewConnector(dsn string, base driver.Connector, driver *Driver) *Connector {
	return &Connector{
		dsn:    dsn,
		base:   base,
		driver: driver,
	}
}

// Connect returns a new driver.Conn.
func (c *Connector) Connect(ctx context.Context) (driver.Conn, error) {
	if c.base == nil {
		return c.driver.Open(c.dsn)
	}

	conn, err := c.base.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return NewConnection(conn, c.driver.Configuration()), nil
}

// Driver returns the Driver of the Connector.
func (c *Connector) Driver() driver.Driver {
	return c.driver
}
