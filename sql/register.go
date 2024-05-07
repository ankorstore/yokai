package sql

import (
	"fmt"
)

const DriverRegistrationPrefix = "yokai"

var (
	GlobalDriverRegistry DriverRegistry
	GlobalDriverFactory  DriverFactory
)

func init() {
	GlobalDriverRegistry = NewDefaultDriverRegistry()
	GlobalDriverFactory = NewDefaultDriverFactory()
}

// Register registers a new Driver for a given name and an optional list of Hook.
func Register(name string, hooks ...Hook) (string, error) {
	registrationName := fmt.Sprintf("%s-%s", DriverRegistrationPrefix, name)

	if GlobalDriverRegistry.Has(registrationName) {
		return registrationName, nil
	}

	system := FetchSystem(name)
	if system == UnknownSystem {
		return "", fmt.Errorf("unsupported database system for driver %s", name)
	}

	driver, err := GlobalDriverFactory.Create(system, hooks...)
	if err != nil {
		return "", err
	}

	err = GlobalDriverRegistry.Add(registrationName, driver)
	if err != nil {
		return "", err
	}

	return registrationName, nil
}
