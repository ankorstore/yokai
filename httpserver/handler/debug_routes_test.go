package handler_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ankorstore/yokai/httpserver/handler"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestDebugRoutesHandler(t *testing.T) {
	t.Parallel()

	httpServer := echo.New()
	httpServer.GET("/debug/routes", handler.DebugRoutesHandler(httpServer))

	req := httptest.NewRequest(http.MethodGet, "/debug/routes", nil)
	rec := httptest.NewRecorder()
	httpServer.ServeHTTP(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Contains(t, rec.Body.String(), `"method":"GET","path":"/debug/routes"`)
}
