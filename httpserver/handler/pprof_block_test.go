package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/httpserver/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestPprofBlockHandler(t *testing.T) {
	t.Parallel()

	httpServer := echo.New()
	httpServer.GET("/debug/pprof/block", handler.PprofBlockHandler())

	req := httptest.NewRequest(http.MethodGet, "/debug/pprof/block", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
}
