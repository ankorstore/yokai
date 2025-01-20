package fxvalidator_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestAsAlias(t *testing.T) {
	t.Parallel()

	result := fxvalidator.AsAlias("alias", "required,alpha")

	assert.Equal(t, "fx.supplyOption", fmt.Sprintf("%T", result))
}

func TestAsValidation(t *testing.T) {
	t.Parallel()

	result := fxvalidator.AsValidation(
		"test-val",
		func(context.Context, validator.FieldLevel) bool {
			return false
		},
		true,
	)

	assert.Equal(t, "fx.supplyOption", fmt.Sprintf("%T", result))
}

func TestAsStructValidation(t *testing.T) {
	t.Parallel()

	result := fxvalidator.AsStructValidation(
		func(context.Context, validator.StructLevel) {},
		TestStruct{},
	)

	assert.Equal(t, "fx.supplyOption", fmt.Sprintf("%T", result))
}

func TestAsCustomType(t *testing.T) {
	t.Parallel()

	result := fxvalidator.AsCustomType(
		func(field reflect.Value) interface{} {
			return true
		},
		TestType{},
	)

	assert.Equal(t, "fx.supplyOption", fmt.Sprintf("%T", result))
}
