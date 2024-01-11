package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/generate/generatetest/uuid"
	"github.com/ankorstore/yokai/httpserver/middleware"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestRequestIdMiddlewareWithDefaults(t *testing.T) {
	t.Parallel()

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		return c.String(
			http.StatusOK,
			c.Request().Header.Get(echo.HeaderXRequestID),
		)
	}

	m := middleware.RequestIdMiddleware()
	h := m(handler)

	err := h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, rec.Header().Get(echo.HeaderXRequestID), rec.Body.String())
}

func TestRequestIdMiddlewareWithSkipper(t *testing.T) {
	t.Parallel()

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	m := middleware.RequestIdMiddlewareWithConfig(middleware.RequestIdMiddlewareConfig{
		Skipper: func(echo.Context) bool {
			return true
		},
	})
	h := m(handler)

	err := h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Empty(t, rec.Header().Get(echo.HeaderXRequestID))
}

func TestRequestIdMiddlewareWithRequestAlreadyContainingId(t *testing.T) {
	t.Parallel()

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(echo.HeaderXRequestID, "test-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		return c.String(
			http.StatusOK,
			c.Request().Header.Get(echo.HeaderXRequestID),
		)
	}

	m := middleware.RequestIdMiddleware()
	h := m(handler)

	err := h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "test-id", rec.Body.String())
	assert.Equal(t, "test-id", rec.Header().Get(echo.HeaderXRequestID))
}

func TestRequestIdMiddlewareWithCustomGenerator(t *testing.T) {
	t.Parallel()

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		return c.String(http.StatusOK, "test")
	}

	generator := uuid.NewTestUuidGenerator("generated-id")

	m := middleware.RequestIdMiddlewareWithConfig(middleware.RequestIdMiddlewareConfig{
		Generator: generator,
	})
	h := m(handler)

	err := h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "generated-id", rec.Header().Get(echo.HeaderXRequestID))
}

func TestRequestIdMiddlewareWithCustomIdHeader(t *testing.T) {
	t.Parallel()

	httpServer := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("custom-header", "custom-id")
	rec := httptest.NewRecorder()

	ctx := httpServer.NewContext(req, rec)
	handler := func(c echo.Context) error {
		return c.String(
			http.StatusOK,
			c.Request().Header.Get("custom-header"),
		)
	}

	m := middleware.RequestIdMiddlewareWithConfig(middleware.RequestIdMiddlewareConfig{
		RequestIdHeader: "custom-header",
	})
	h := m(handler)

	err := h(ctx)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "custom-id", rec.Body.String())
	assert.Equal(t, "custom-id", rec.Header().Get("custom-header"))
}
