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

func TestDebugVersionHandler(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("../testdata/config"),
	)
	assert.NoError(t, err)

	httpServer := echo.New()
	httpServer.GET("/debug/version", handler.DebugVersionHandler(cfg))

	req := httptest.NewRequest(http.MethodGet, "/debug/version", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(
		t,
		rec.Body.String(),
		`{"application":"test-app","version":"0.1.0"}`,
	)
}
