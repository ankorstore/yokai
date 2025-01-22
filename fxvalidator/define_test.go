package fxvalidator_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestAliasDefinition(t *testing.T) {
	t.Parallel()

	def := fxvalidator.NewAliasDefinition("foo", "bar,baz")

	assert.Equal(t, "foo", def.Alias())
	assert.Equal(t, "bar,baz", def.Tags())
}

func TestValidationDefinition(t *testing.T) {
	t.Parallel()

	fn := func(ctx context.Context, fl validator.FieldLevel) bool {
		return false
	}

	def := fxvalidator.NewValidationDefinition("foo", fn, true)
	defFn := def.Fn()

	assert.Equal(t, "foo", def.Tag())
	assert.IsType(t, new(validator.FuncCtx), &defFn)
	assert.True(t, def.CallEvenIfNull())
}

func TestStructValidationDefinition(t *testing.T) {
	t.Parallel()

	fn := func(ctx context.Context, sl validator.StructLevel) {}

	def := fxvalidator.NewStructValidationDefinition(fn, TestStruct{}, TestStructWithTestType{})
	defFn := def.Fn()

	assert.IsType(t, new(validator.StructLevelFuncCtx), &defFn)
	assert.Equal(t, []any{TestStruct{}, TestStructWithTestType{}}, def.Types())
}

func TestCustomTypeDefinition(t *testing.T) {
	t.Parallel()

	fn := func(field reflect.Value) interface{} {
		return nil
	}

	def := fxvalidator.NewCustomTypeDefinition(fn, TestType{})
	defFn := def.Fn()

	assert.IsType(t, new(validator.CustomTypeFunc), &defFn)
	assert.Equal(t, []any{TestType{}}, def.Types())
}
