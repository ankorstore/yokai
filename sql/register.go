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
func Register(driver string, hooks ...Hook) (string, error) {
	driverName := fmt.Sprintf("%s-%s", DriverRegistrationPrefix, driver)

	if GlobalDriverRegistry.Has(driverName) {
		return driverName, nil
	}

	system := FetchSystem(driver)
	if system == UnknownSystem {
		return "", fmt.Errorf("unsupported database system for driver %s", driver)
	}

	drv, err := GlobalDriverFactory.Create(system, hooks...)
	if err != nil {
		return "", err
	}

	err = GlobalDriverRegistry.Add(driverName, drv)
	if err != nil {
		return "", err
	}

	return driverName, nil
}

// RegisterNamed registers a named Driver for a given name and an optional list of Hook.
func RegisterNamed(driver string, name string, hooks ...Hook) (string, error) {
	driverName := fmt.Sprintf("%s-%s-%s", DriverRegistrationPrefix, driver, name)

	if GlobalDriverRegistry.Has(driverName) {
		return driverName, nil
	}

	system := FetchSystem(driver)
	if system == UnknownSystem {
		return "", fmt.Errorf("unsupported database system for driver %s", driver)
	}

	drv, err := GlobalDriverFactory.Create(system, hooks...)
	if err != nil {
		return "", err
	}

	err = GlobalDriverRegistry.Add(driverName, drv)
	if err != nil {
		return "", err
	}

	return driverName, nil
}
