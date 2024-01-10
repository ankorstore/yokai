package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/httpserver/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDebugConfigHandler(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("../testdata/config"),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.GET("/debug/config", handler.DebugConfigHandler(cfg))

	req := httptest.NewRequest(http.MethodGet, "/debug/config", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(
		t,
		rec.Body.String(),
		`{"app":{"debug":true,"env":"test","name":"test-app","version":"0.1.0"},"config":{"some":"value"}}`,
	)
}
