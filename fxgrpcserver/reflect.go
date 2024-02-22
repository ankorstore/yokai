package fxgrpcserver

import (
	"reflect"
)

func GetType(target any) string {
	return reflect.TypeOf(target).String()
}

func GetReturnType(target any) string {
	return reflect.TypeOf(target).Out(0).String()
}
