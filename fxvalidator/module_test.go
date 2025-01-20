package fxvalidator_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type TestStruct struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Business string `validate:"oneof=brand retailer"`
}

type TestStructWithPrivate struct {
	private string `validate:"required,alpha"`
}

type TestType struct {
	Value string
}

type TestStructWithTestType struct {
	TestType TestType `validate:"required"`
}

func TestModule(t *testing.T) {
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	runTest := func(tb testing.TB, options ...fx.Option) *validator.Validate {
		tb.Helper()

		var validate *validator.Validate

		fxtest.New(
			t,
			fx.NopLogger,
			fxvalidator.FXValidatorModule,
			fxconfig.FxConfigModule,
			fx.Options(options...),
			fx.Populate(&validate),
		).RequireStart().RequireStop()

		return validate
	}

	t.Run("test validation success", func(t *testing.T) {
		validate := runTest(t)

		err := validate.Struct(TestStruct{
			Name:     "retailer",
			Email:    "retailer@example.com",
			Business: "retailer",
		})
		assert.NoError(t, err)

		err = validate.StructCtx(context.Background(), TestStruct{
			Name:     "retailer",
			Email:    "retailer@example.com",
			Business: "retailer",
		})
		assert.NoError(t, err)
	})

	t.Run("test validation error", func(t *testing.T) {
		validate := runTest(t)

		err := validate.Struct(TestStruct{
			Name:     "",
			Email:    "invalid",
			Business: "invalid",
		})
		assert.Error(t, err)

		var validationErrors validator.ValidationErrors
		ok := errors.As(err, &validationErrors)
		assert.True(t, ok)

		for _, vErr := range validationErrors {
			if vErr.StructField() == "Name" {
				assert.Equal(t, "Key: 'TestStruct.Name' Error:Field validation for 'Name' failed on the 'required' tag", vErr.Error())
			}

			if vErr.StructField() == "Email" {
				assert.Equal(t, "Key: 'TestStruct.Email' Error:Field validation for 'Email' failed on the 'email' tag", vErr.Error())
			}

			if vErr.StructField() == "Business" {
				assert.Equal(t, "Key: 'TestStruct.Business' Error:Field validation for 'Business' failed on the 'oneof' tag", vErr.Error())
			}
		}

		err = validate.StructCtx(context.Background(), TestStruct{
			Name:     "",
			Email:    "invalid",
			Business: "invalid",
		})
		assert.Error(t, err)

		ok = errors.As(err, &validationErrors)
		assert.True(t, ok)

		for _, vErr := range validationErrors {
			if vErr.StructField() == "Name" {
				assert.Equal(t, "Key: 'TestStruct.Name' Error:Field validation for 'Name' failed on the 'required' tag", vErr.Error())
			}

			if vErr.StructField() == "Email" {
				assert.Equal(t, "Key: 'TestStruct.Email' Error:Field validation for 'Email' failed on the 'email' tag", vErr.Error())
			}

			if vErr.StructField() == "Business" {
				assert.Equal(t, "Key: 'TestStruct.Business' Error:Field validation for 'Business' failed on the 'oneof' tag", vErr.Error())
			}
		}
	})

	t.Run("test validation success with private fields", func(t *testing.T) {
		t.Setenv("PRIVATE_FIELDS", "true")

		validate := runTest(t)

		err := validate.Struct(TestStructWithPrivate{
			private: "abc",
		})
		assert.NoError(t, err)

		err = validate.StructCtx(context.Background(), TestStructWithPrivate{
			private: "abc",
		})
		assert.NoError(t, err)
	})

	t.Run("test validation failure with private fields", func(t *testing.T) {
		t.Setenv("PRIVATE_FIELDS", "true")

		validate := runTest(t)

		err := validate.Struct(TestStructWithPrivate{
			private: "123",
		})
		assert.Error(t, err)

		var validationErrors validator.ValidationErrors
		ok := errors.As(err, &validationErrors)
		assert.True(t, ok)

		for _, vErr := range validationErrors {
			if vErr.StructField() == "private" {
				assert.Equal(t, "Key: 'TestStructWithPrivate.private' Error:Field validation for 'private' failed on the 'alpha' tag", vErr.Error())
			}
		}

		err = validate.StructCtx(context.Background(), TestStructWithPrivate{
			private: "123",
		})
		assert.Error(t, err)

		ok = errors.As(err, &validationErrors)
		assert.True(t, ok)

		for _, vErr := range validationErrors {
			if vErr.StructField() == "private" {
				assert.Equal(t, "Key: 'TestStructWithPrivate.private' Error:Field validation for 'private' failed on the 'alpha' tag", vErr.Error())
			}
		}
	})

	t.Run("test validation success with custom alias", func(t *testing.T) {
		validate := runTest(t, fxvalidator.AsAlias("test-alias", "required,alpha,max=10"))

		err := validate.Var("abcdefghi", "test-alias")
		assert.NoError(t, err)

		err = validate.VarCtx(context.Background(), "abcdefghi", "test-alias")
		assert.NoError(t, err)
	})

	t.Run("test validation error with custom alias", func(t *testing.T) {
		validate := runTest(t, fxvalidator.AsAlias("test-alias", "required,alpha,max=10"))

		err := validate.Var("invalid-1234", "test-alias")
		assert.Error(t, err)

		var validationError validator.ValidationErrors
		ok := errors.As(err, &validationError)
		assert.True(t, ok)
		assert.Contains(t, validationError.Error(), "failed on the 'test-alias' tag")

		err = validate.VarCtx(context.Background(), "invalid-1234", "test-alias")
		assert.Error(t, err)

		ok = errors.As(err, &validationError)
		assert.True(t, ok)
		assert.Contains(t, validationError.Error(), "failed on the 'test-alias' tag")
	})

	t.Run("test validation success with custom validation", func(t *testing.T) {
		fn := func(ctx context.Context, fl validator.FieldLevel) bool {
			return fl.Field().String() == "expected"
		}

		validate := runTest(t, fxvalidator.AsValidation("test-custom", fn, true))

		err := validate.Var("expected", "test-custom")
		assert.NoError(t, err)

		err = validate.VarCtx(context.Background(), "expected", "test-custom")
		assert.NoError(t, err)
	})

	t.Run("test validation error with custom validation", func(t *testing.T) {
		fn := func(ctx context.Context, fl validator.FieldLevel) bool {
			return fl.Field().String() == "expected"
		}

		validate := runTest(t, fxvalidator.AsValidation("test-custom", fn, true))

		err := validate.Var("invalid-1234", "test-custom")
		assert.Error(t, err)

		var validationError validator.ValidationErrors
		ok := errors.As(err, &validationError)
		assert.True(t, ok)
		assert.Contains(t, validationError.Error(), "failed on the 'test-custom' tag")

		err = validate.VarCtx(context.Background(), "invalid-1234", "test-custom")
		assert.Error(t, err)

		ok = errors.As(err, &validationError)
		assert.True(t, ok)
		assert.Contains(t, validationError.Error(), "failed on the 'test-custom' tag")
	})

	t.Run("test validation success with custom struct validation", func(t *testing.T) {
		fn := func(ctx context.Context, sl validator.StructLevel) {
			st, ok := sl.Current().Interface().(TestStruct)

			if ok {
				if st.Business == "retailer" && !strings.Contains(st.Email, "retailer") {
					sl.ReportError(st.Email, "Email", "Email", "invalid", "invalid retailer email")
				}
			}
		}

		validate := runTest(t, fxvalidator.AsStructValidation(fn, TestStruct{}))

		err := validate.StructCtx(context.Background(), TestStruct{
			Name:     "retailer",
			Email:    "retailer@example.com",
			Business: "retailer",
		})
		assert.NoError(t, err)
	})

	t.Run("test validation error with custom struct validation", func(t *testing.T) {
		fn := func(ctx context.Context, sl validator.StructLevel) {
			st, ok := sl.Current().Interface().(TestStruct)

			if ok {
				if st.Business == "retailer" && !strings.Contains(st.Email, "retailer") {
					sl.ReportError(st.Email, "Email", "Email", "invalid-retailer-email", "invalid retailer email")
				}
			}
		}

		validate := runTest(t, fxvalidator.AsStructValidation(fn, TestStruct{}))

		err := validate.StructCtx(context.Background(), TestStruct{
			Name:     "retailer",
			Email:    "invalid@example.com",
			Business: "retailer",
		})
		assert.Error(t, err)

		var validationError validator.ValidationErrors
		ok := errors.As(err, &validationError)
		assert.True(t, ok)
		assert.Contains(t, validationError.Error(), "failed on the 'invalid-retailer-email' tag")
	})

	t.Run("test validation success with custom type", func(t *testing.T) {
		fn := func(field reflect.Value) interface{} {
			if ct, ok := field.Interface().(TestType); ok {
				if ct.Value == "invalid" {
					return ""
				}

				return ct.Value
			}

			return ""
		}

		validate := runTest(t, fxvalidator.AsCustomType(fn, TestType{}))

		err := validate.StructCtx(context.Background(), TestStructWithTestType{
			TestType: TestType{
				Value: "valid",
			},
		})
		assert.NoError(t, err)
	})

	t.Run("test validation error with custom type", func(t *testing.T) {
		fn := func(field reflect.Value) interface{} {
			if ct, ok := field.Interface().(TestType); ok {
				if ct.Value == "invalid" {
					return ""
				}

				return ct.Value
			}

			return ""
		}

		validate := runTest(t, fxvalidator.AsCustomType(fn, TestType{}))

		err := validate.StructCtx(context.Background(), TestStructWithTestType{
			TestType: TestType{
				Value: "invalid",
			},
		})
		assert.Error(t, err)

		var validationError validator.ValidationErrors
		ok := errors.As(err, &validationError)
		assert.True(t, ok)
		assert.Equal(t, "Key: 'TestStructWithTestType.TestType' Error:Field validation for 'TestType' failed on the 'required' tag", validationError.Error())
	})
}
