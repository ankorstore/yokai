package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/httpserver/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPprofSymbolHandler(t *testing.T) {
	t.Parallel()

	httpServer := echo.New()
	httpServer.GET("/debug/pprof/symbol", handler.PprofSymbolHandler())

	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/symbol", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
