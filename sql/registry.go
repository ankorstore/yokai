package sql

import (
	"fmt"
	"sync"
)

type DriverRegistry interface {
	Has(name string) bool
	Add(name string, driver *Driver)
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

func (r *DefaultDriverRegistry) Add(name string, driver *Driver) {
	if r.Has(name) {
		return
	}

	r.mutex.RLock()
	defer r.mutex.RUnlock()

	r.drivers[name] = driver
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
