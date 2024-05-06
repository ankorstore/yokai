package sql

import (
	"context"
	"database/sql/driver"
)

// Connector is a SQL driver connector.
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

func (c *Connector) Driver() driver.Driver {
	return c.driver
}
