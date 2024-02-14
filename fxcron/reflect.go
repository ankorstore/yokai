package fxcron

import (
	"reflect"
)

// GetType returns the type of a target.
func GetType(target any) string {
	return reflect.TypeOf(target).String()
}

// GetReturnType returns the return type of a target.
func GetReturnType(target any) string {
	return reflect.TypeOf(target).Out(0).String()
}
