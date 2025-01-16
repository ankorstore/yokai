package fxvalidator_test

import (
	"errors"
	"testing"

	"github.com/ankorstore/yokai/fxconfig"
	"github.com/ankorstore/yokai/fxvalidator"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

type testStruct struct {
	Name     string `validate:"required"`
	Email    string `validate:"required,email"`
	Business string `validate:"oneof=brand retailer"`
}

type testPrivateStruct struct {
	private string `validate:"required,alpha"`
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

		err := validate.Struct(testStruct{
			Name:     "retailer",
			Email:    "retailer@example.com",
			Business: "retailer",
		})
		assert.NoError(t, err)
	})

	t.Run("test validation error", func(t *testing.T) {
		validate := runTest(t)

		err := validate.Struct(testStruct{
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
				assert.Equal(t, "Key: 'testStruct.Name' Error:Field validation for 'Name' failed on the 'required' tag", vErr.Error())
			}

			if vErr.StructField() == "Email" {
				assert.Equal(t, "Key: 'testStruct.Email' Error:Field validation for 'Email' failed on the 'email' tag", vErr.Error())
			}

			if vErr.StructField() == "Business" {
				assert.Equal(t, "Key: 'testStruct.Business' Error:Field validation for 'Business' failed on the 'oneof' tag", vErr.Error())
			}
		}
	})

	t.Run("test validation success with private fields", func(t *testing.T) {
		t.Setenv("PRIVATE_FIELDS", "true")

		validate := runTest(t)

		err := validate.Struct(testPrivateStruct{
			private: "abc",
		})
		assert.NoError(t, err)
	})

	t.Run("test validation failure with private fields", func(t *testing.T) {
		t.Setenv("PRIVATE_FIELDS", "true")

		validate := runTest(t)

		err := validate.Struct(testPrivateStruct{
			private: "123",
		})
		assert.Error(t, err)

		var validationErrors validator.ValidationErrors
		ok := errors.As(err, &validationErrors)
		assert.True(t, ok)

		for _, vErr := range validationErrors {
			if vErr.StructField() == "private" {
				assert.Equal(t, "Key: 'testPrivateStruct.private' Error:Field validation for 'private' failed on the 'alpha' tag", vErr.Error())
			}
		}
	})

	t.Run("test validation success with alias", func(t *testing.T) {
		validate := runTest(t, fxvalidator.AsAlias("test-alias", "required,alpha,max=10"))

		err := validate.Var("abcdefghi", "test-alias")
		assert.NoError(t, err)
	})

	t.Run("test validation error with alias", func(t *testing.T) {
		validate := runTest(t, fxvalidator.AsAlias("test-alias", "required,alpha,max=10"))

		err := validate.Var("abcdefg99", "test-alias")
		assert.NoError(t, err)
	})
}
