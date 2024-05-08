package sql

import (
	"database/sql"
	"fmt"
	"sync"
)

// DriverRegistry is the interface for Driver registries.
type DriverRegistry interface {
	Has(name string) bool
	Add(name string, driver *Driver) error
	Get(name string) (*Driver, error)
}

// DefaultDriverRegistry is the default DriverRegistry implementation.
type DefaultDriverRegistry struct {
	drivers map[string]*Driver
	mutex   sync.RWMutex
}

// NewDefaultDriverRegistry returns a new DefaultDriverRegistry.
func NewDefaultDriverRegistry() *DefaultDriverRegistry {
	return &DefaultDriverRegistry{
		drivers: make(map[string]*Driver),
	}
}

// Add adds and register a given Driver for a name.
func (r *DefaultDriverRegistry) Add(name string, driver *Driver) (err error) {
	if r.Has(name) {
		return nil
	}

	r.mutex.Lock()

	defer func() {
		r.mutex.Unlock()

		if rec := recover(); rec != nil {
			err = fmt.Errorf("cannot add driver %s: %v", name, rec)
		}
	}()

	r.drivers[name] = driver

	sql.Register(name, driver)

	return err
}

// Has returns true is a driver is already registered for a given name.
func (r *DefaultDriverRegistry) Has(name string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, ok := r.drivers[name]

	return ok
}

// Get returns a registered driver for a given name.
func (r *DefaultDriverRegistry) Get(name string) (*Driver, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	registeredDriver, ok := r.drivers[name]
	if !ok {
		return nil, fmt.Errorf("cannot find driver %s in driver registry", name)
	}

	return registeredDriver, nil
}
