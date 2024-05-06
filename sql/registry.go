package sql

import (
	"database/sql"
	"fmt"
	"sync"
)

type DriverRegistry interface {
	Has(name string) bool
	Add(name string, driver *Driver) error
	Get(name string) (*Driver, error)
}

type DefaultDriverRegistry struct {
	drivers map[string]*Driver
	mutex   sync.RWMutex
}

func NewDefaultDriverRegistry() *DefaultDriverRegistry {
	return &DefaultDriverRegistry{
		drivers: make(map[string]*Driver),
	}
}

func (r *DefaultDriverRegistry) Add(name string, driver *Driver) error {
	if r.Has(name) {
		return nil
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.drivers[name] = driver

	sql.Register(name, driver)
	if rec := recover(); rec != nil {
		return fmt.Errorf("cannot register driver %s: %v", name, rec)
	}

	return nil
}

func (r *DefaultDriverRegistry) Has(name string) bool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	_, ok := r.drivers[name]

	return ok
}

func (r *DefaultDriverRegistry) Get(name string) (*Driver, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	registeredDriver, ok := r.drivers[name]
	if !ok {
		return nil, fmt.Errorf("cannot find driver %s in driver registry", name)
	}

	return registeredDriver, nil
}
