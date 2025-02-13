package fxhttpserver_test

import (
	"fmt"
	"testing"

	"github.com/ankorstore/yokai/fxhttpserver"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var testMiddleware = func(next echo.HandlerFunc) echo.HandlerFunc {
	return next
}

var testHandler = func(c echo.Context) error {
	return nil
}

var testErrorHandler = func(err error, c echo.Context) {}

type testMiddlewareDefinitionMock struct {
	mock.Mock
}

func (m *testMiddlewareDefinitionMock) Concrete() bool {
	args := m.Called()

	return args.Bool(0)
}

func (m *testMiddlewareDefinitionMock) Middleware() any {
	args := m.Called()

	return args.Get(0)
}

func (m *testMiddlewareDefinitionMock) Kind() fxhttpserver.MiddlewareKind {
	args := m.Called()

	//nolint:forcetypeassert
	return args.Get(0).(fxhttpserver.MiddlewareKind)
}

type testHandlerDefinitionMock struct {
	mock.Mock
}

func (m *testHandlerDefinitionMock) Concrete() bool {
	args := m.Called()

	return args.Bool(0)
}

func (m *testHandlerDefinitionMock) Method() string {
	args := m.Called()

	return args.String(0)
}

func (m *testHandlerDefinitionMock) Path() string {
	args := m.Called()

	return args.String(0)
}

func (m *testHandlerDefinitionMock) Handler() any {
	args := m.Called()

	return args.Get(0)
}

func (m *testHandlerDefinitionMock) Middlewares() []fxhttpserver.MiddlewareDefinition {
	args := m.Called()

	//nolint:forcetypeassert
	return args.Get(0).([]fxhttpserver.MiddlewareDefinition)
}

type testMiddlewareImplementation struct{}

func (m testMiddlewareImplementation) Handle() echo.MiddlewareFunc {
	return testMiddleware
}

type testHandlerImplementation struct{}

func (h testHandlerImplementation) Handle() echo.HandlerFunc {
	return testHandler
}

type testErrorHandlerImplementation struct{}

func (h testErrorHandlerImplementation) Handle() echo.HTTPErrorHandler { return testErrorHandler }

func TestNewFxHttpServerRegistry(t *testing.T) {
	t.Parallel()

	param := fxhttpserver.FxHttpServerRegistryParam{
		Middlewares: []fxhttpserver.Middleware{testMiddlewareImplementation{}},
		Handlers:    []fxhttpserver.Handler{testHandlerImplementation{}},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	assert.IsType(t, &fxhttpserver.HttpServerRegistry{}, registry)
}

func TestResolveMiddlewaresSuccessWithAnonFuncType(t *testing.T) {
	t.Parallel()

	param := fxhttpserver.FxHttpServerRegistryParam{
		Middlewares: []fxhttpserver.Middleware{
			testMiddlewareImplementation{},
		},
		MiddlewareDefinitions: []fxhttpserver.MiddlewareDefinition{
			fxhttpserver.NewMiddlewareDefinition(testMiddleware, fxhttpserver.GlobalUse),
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	resolvedMiddlewares, err := registry.ResolveMiddlewares()
	assert.NoError(t, err)

	assert.Len(t, resolvedMiddlewares, 1)
	assert.Equal(t, fxhttpserver.GlobalUse, resolvedMiddlewares[0].Kind())
	assert.Equal(t, testMiddleware(testHandler)(nil), resolvedMiddlewares[0].Middleware()(testHandler)(nil))
}

func TestResolveMiddlewaresFailureOnInvalidImplementation(t *testing.T) {
	t.Parallel()

	middlewareDefinitionMock := new(testMiddlewareDefinitionMock)
	middlewareDefinitionMock.On("Concrete").Return(true)
	middlewareDefinitionMock.On("Middleware").Return(nil)
	middlewareDefinitionMock.On("Kind").Return(fxhttpserver.GlobalUse)

	param := fxhttpserver.FxHttpServerRegistryParam{
		Middlewares: []fxhttpserver.Middleware{
			testMiddlewareImplementation{},
		},
		MiddlewareDefinitions: []fxhttpserver.MiddlewareDefinition{
			middlewareDefinitionMock,
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	_, err := registry.ResolveMiddlewares()
	assert.Error(t, err)
	assert.Equal(t, "cannot cast middleware definition as MiddlewareFunc", err.Error())
}

func TestResolveMiddlewaresFailureOnMissingImplementation(t *testing.T) {
	t.Parallel()

	middlewareDefinitionMock := new(testMiddlewareDefinitionMock)
	middlewareDefinitionMock.On("Concrete").Return(false)
	middlewareDefinitionMock.On("Middleware").Return("invalid")
	middlewareDefinitionMock.On("Kind").Return(fxhttpserver.GlobalUse)

	param := fxhttpserver.FxHttpServerRegistryParam{
		Middlewares: []fxhttpserver.Middleware{},
		MiddlewareDefinitions: []fxhttpserver.MiddlewareDefinition{
			middlewareDefinitionMock,
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	_, err := registry.ResolveMiddlewares()
	assert.Error(t, err)
	assert.Equal(t, "cannot lookup registered middleware", err.Error())
}

func TestResolveHandlersSuccessWithAnonFuncType(t *testing.T) {
	t.Parallel()

	param := fxhttpserver.FxHttpServerRegistryParam{
		Handlers: []fxhttpserver.Handler{
			testHandlerImplementation{},
		},
		HandlerDefinitions: []fxhttpserver.HandlerDefinition{
			fxhttpserver.NewHandlerDefinition(
				"GET",
				"/path",
				testHandler,
				[]fxhttpserver.MiddlewareDefinition{
					fxhttpserver.NewMiddlewareDefinition(testMiddleware, fxhttpserver.GlobalUse),
				},
			),
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	resolvedHandlers, err := registry.ResolveHandlers()
	assert.NoError(t, err)

	assert.Len(t, resolvedHandlers, 1)
	assert.Equal(t, "GET", resolvedHandlers[0].Method())
	assert.Equal(t, "/path", resolvedHandlers[0].Path())
	assert.Equal(t, testHandler(nil), resolvedHandlers[0].Handler()(nil))
}

func TestResolveHandlersSuccessWithHandlerFuncType(t *testing.T) {
	t.Parallel()

	var h echo.HandlerFunc

	handlerDefinitionMock := new(testHandlerDefinitionMock)
	handlerDefinitionMock.On("Concrete").Return(true)
	handlerDefinitionMock.On("Method").Return("GET")
	handlerDefinitionMock.On("Path").Return("/path")
	handlerDefinitionMock.On("Handler").Return(h)
	handlerDefinitionMock.On("Middlewares").Return([]fxhttpserver.MiddlewareDefinition{})

	param := fxhttpserver.FxHttpServerRegistryParam{
		Handlers: []fxhttpserver.Handler{
			testHandlerImplementation{},
		},
		HandlerDefinitions: []fxhttpserver.HandlerDefinition{
			handlerDefinitionMock,
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	resolvedHandlers, err := registry.ResolveHandlers()
	assert.NoError(t, err)

	assert.Len(t, resolvedHandlers, 1)
	assert.Equal(t, "GET", resolvedHandlers[0].Method())
	assert.Equal(t, "/path", resolvedHandlers[0].Path())
}

func TestResolveHandlersFailureOnMissingHandlerImplementation(t *testing.T) {
	t.Parallel()

	handlerDefinitionMock := new(testHandlerDefinitionMock)
	handlerDefinitionMock.On("Concrete").Return(false)
	handlerDefinitionMock.On("Method").Return("GET")
	handlerDefinitionMock.On("Path").Return("/path")
	handlerDefinitionMock.On("Handler").Return("invalid")
	handlerDefinitionMock.On("Middlewares").Return([]fxhttpserver.MiddlewareDefinition{})

	param := fxhttpserver.FxHttpServerRegistryParam{
		Handlers: []fxhttpserver.Handler{},
		HandlerDefinitions: []fxhttpserver.HandlerDefinition{
			handlerDefinitionMock,
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	_, err := registry.ResolveHandlers()
	assert.Error(t, err)
	assert.Equal(t, "cannot lookup registered handler", err.Error())
}

func TestResolveHandlersFailureOnInvalidHandlerImplementation(t *testing.T) {
	t.Parallel()

	handlerDefinitionMock := new(testHandlerDefinitionMock)
	handlerDefinitionMock.On("Concrete").Return(true)
	handlerDefinitionMock.On("Method").Return("GET")
	handlerDefinitionMock.On("Path").Return("/path")
	handlerDefinitionMock.On("Handler").Return("invalid")
	handlerDefinitionMock.On("Middlewares").Return([]fxhttpserver.MiddlewareDefinition{})

	param := fxhttpserver.FxHttpServerRegistryParam{
		Handlers: []fxhttpserver.Handler{
			testHandlerImplementation{},
		},
		HandlerDefinitions: []fxhttpserver.HandlerDefinition{
			handlerDefinitionMock,
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	_, err := registry.ResolveHandlers()
	assert.Error(t, err)
	assert.Equal(t, "cannot cast handler definition as HandlerFunc", err.Error())
}

func TestResolveHandlersFailureOnInvalidMiddlewareImplementation(t *testing.T) {
	t.Parallel()

	middlewareDefinitionMock := new(testMiddlewareDefinitionMock)
	middlewareDefinitionMock.On("Concrete").Return(true)
	middlewareDefinitionMock.On("Middleware").Return("invalid")
	middlewareDefinitionMock.On("Kind").Return(fxhttpserver.GlobalUse)

	handlerDefinitionMock := new(testHandlerDefinitionMock)
	handlerDefinitionMock.On("Concrete").Return(true)
	handlerDefinitionMock.On("Method").Return("GET")
	handlerDefinitionMock.On("Path").Return("/path")
	handlerDefinitionMock.On("Handler").Return("invalid")
	handlerDefinitionMock.On("Middlewares").Return([]fxhttpserver.MiddlewareDefinition{middlewareDefinitionMock})

	param := fxhttpserver.FxHttpServerRegistryParam{
		Handlers: []fxhttpserver.Handler{
			testHandlerImplementation{},
		},
		HandlerDefinitions: []fxhttpserver.HandlerDefinition{
			handlerDefinitionMock,
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	_, err := registry.ResolveHandlers()
	assert.Error(t, err)
	assert.Equal(t, "cannot cast middleware definition as MiddlewareFunc", err.Error())
}

func TestResolveHandlersFailureOnMissingMiddlewareImplementation(t *testing.T) {
	t.Parallel()

	middlewareDefinitionMock := new(testMiddlewareDefinitionMock)
	middlewareDefinitionMock.On("Concrete").Return(false)
	middlewareDefinitionMock.On("Middleware").Return("invalid")
	middlewareDefinitionMock.On("Kind").Return(fxhttpserver.GlobalUse)

	handlerDefinitionMock := new(testHandlerDefinitionMock)
	handlerDefinitionMock.On("Concrete").Return(true)
	handlerDefinitionMock.On("Method").Return("GET")
	handlerDefinitionMock.On("Path").Return("/path")
	handlerDefinitionMock.On("Handler").Return("invalid")
	handlerDefinitionMock.On("Middlewares").Return([]fxhttpserver.MiddlewareDefinition{middlewareDefinitionMock})

	param := fxhttpserver.FxHttpServerRegistryParam{
		Handlers: []fxhttpserver.Handler{
			testHandlerImplementation{},
		},
		HandlerDefinitions: []fxhttpserver.HandlerDefinition{
			handlerDefinitionMock,
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	_, err := registry.ResolveHandlers()
	assert.Error(t, err)
	assert.Equal(t, "cannot lookup registered middleware", err.Error())
}

func TestResolveHandlersGroupsSuccess(t *testing.T) {
	t.Parallel()

	param := fxhttpserver.FxHttpServerRegistryParam{
		Handlers: []fxhttpserver.Handler{
			testHandlerImplementation{},
		},
		HandlersGroupDefinitions: []fxhttpserver.HandlersGroupDefinition{
			fxhttpserver.NewHandlersGroupDefinition(
				"/group",
				[]fxhttpserver.HandlerDefinition{
					fxhttpserver.NewHandlerDefinition(
						"GET",
						"/path",
						testHandler,
						[]fxhttpserver.MiddlewareDefinition{
							fxhttpserver.NewMiddlewareDefinition(testMiddleware, fxhttpserver.GlobalUse),
						},
					),
				},
				[]fxhttpserver.MiddlewareDefinition{},
			),
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	resolvedGroups, err := registry.ResolveHandlersGroups()
	assert.NoError(t, err)

	assert.Len(t, resolvedGroups, 1)
	assert.Equal(t, "/group", resolvedGroups[0].Prefix())
}

func TestResolveHandlersGroupFailureOnMissingGroupMiddlewareImplementation(t *testing.T) {
	t.Parallel()

	middlewareDefinitionMock := new(testMiddlewareDefinitionMock)
	middlewareDefinitionMock.On("Concrete").Return(false)
	middlewareDefinitionMock.On("Middleware").Return("invalid")
	middlewareDefinitionMock.On("Kind").Return(fxhttpserver.GlobalUse)

	handlerDefinitionMock := new(testHandlerDefinitionMock)
	handlerDefinitionMock.On("Concrete").Return(true)
	handlerDefinitionMock.On("Method").Return("GET")
	handlerDefinitionMock.On("Path").Return("/path")
	handlerDefinitionMock.On("Handler").Return("invalid")
	handlerDefinitionMock.On("Middlewares").Return([]fxhttpserver.MiddlewareDefinition{middlewareDefinitionMock})

	param := fxhttpserver.FxHttpServerRegistryParam{
		Handlers: []fxhttpserver.Handler{
			testHandlerImplementation{},
		},
		HandlersGroupDefinitions: []fxhttpserver.HandlersGroupDefinition{
			fxhttpserver.NewHandlersGroupDefinition(
				"/group",
				[]fxhttpserver.HandlerDefinition{
					handlerDefinitionMock,
				},
				[]fxhttpserver.MiddlewareDefinition{
					middlewareDefinitionMock,
				},
			),
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	_, err := registry.ResolveHandlersGroups()
	assert.Error(t, err)
	assert.Equal(t, "cannot lookup registered middleware", err.Error())
}

func TestResolveHandlersGroupFailureOnInvalidHandlerImplementation(t *testing.T) {
	t.Parallel()

	handlerDefinitionMock := new(testHandlerDefinitionMock)
	handlerDefinitionMock.On("Concrete").Return(true)
	handlerDefinitionMock.On("Method").Return("GET")
	handlerDefinitionMock.On("Path").Return("/path")
	handlerDefinitionMock.On("Handler").Return("invalid")
	handlerDefinitionMock.On("Middlewares").Return([]fxhttpserver.MiddlewareDefinition{})

	param := fxhttpserver.FxHttpServerRegistryParam{
		Handlers: []fxhttpserver.Handler{
			testHandlerImplementation{},
		},
		HandlersGroupDefinitions: []fxhttpserver.HandlersGroupDefinition{
			fxhttpserver.NewHandlersGroupDefinition(
				"/group",
				[]fxhttpserver.HandlerDefinition{
					handlerDefinitionMock,
				},
				[]fxhttpserver.MiddlewareDefinition{},
			),
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	_, err := registry.ResolveHandlersGroups()
	assert.Error(t, err)
	assert.Equal(t, "cannot cast handler definition as HandlerFunc", err.Error())
}

func TestResolveHandlersGroupFailureOnInvalidHandlerMiddlewareImplementation(t *testing.T) {
	t.Parallel()

	middlewareDefinitionMock := new(testMiddlewareDefinitionMock)
	middlewareDefinitionMock.On("Concrete").Return(true)
	middlewareDefinitionMock.On("Middleware").Return("invalid")
	middlewareDefinitionMock.On("Kind").Return(fxhttpserver.GlobalUse)

	handlerDefinitionMock := new(testHandlerDefinitionMock)
	handlerDefinitionMock.On("Concrete").Return(true)
	handlerDefinitionMock.On("Method").Return("GET")
	handlerDefinitionMock.On("Path").Return("/path")
	handlerDefinitionMock.On("Handler").Return("invalid")
	handlerDefinitionMock.On("Middlewares").Return([]fxhttpserver.MiddlewareDefinition{middlewareDefinitionMock})

	param := fxhttpserver.FxHttpServerRegistryParam{
		Handlers: []fxhttpserver.Handler{
			testHandlerImplementation{},
		},
		HandlersGroupDefinitions: []fxhttpserver.HandlersGroupDefinition{
			fxhttpserver.NewHandlersGroupDefinition(
				"/group",
				[]fxhttpserver.HandlerDefinition{
					handlerDefinitionMock,
				},
				[]fxhttpserver.MiddlewareDefinition{},
			),
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	_, err := registry.ResolveHandlersGroups()
	assert.Error(t, err)
	assert.Equal(t, "cannot cast middleware definition as MiddlewareFunc", err.Error())
}

func TestResolveErrorHandlerSuccess(t *testing.T) {
	t.Parallel()

	param := fxhttpserver.FxHttpServerRegistryParam{
		ErrorHandlers: []fxhttpserver.ErrorHandler{
			testErrorHandlerImplementation{},
		},
	}
	registry := fxhttpserver.NewFxHttpServerRegistry(param)

	resolvedErrorHandlers := registry.ResolveErrorHandlers()

	assert.Len(t, resolvedErrorHandlers, 1)

	assert.Equal(t, "echo.HTTPErrorHandler", fmt.Sprintf("%T", resolvedErrorHandlers[0].Handle()))
}
