package sse_test

import (
	"context"
	"testing"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestMCPSSEServer(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("../../testdata/config"),
	)
	assert.NoError(t, err)

	lb := logtest.NewDefaultTestLogBuffer()
	lg, err := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(lb))
	assert.NoError(t, err)

	mcpSrv := server.NewMCPServer("test-server", "1.0.0")

	srv := sse.NewDefaultMCPSSEServerFactory(cfg).Create(mcpSrv)

	assert.False(t, srv.Running())

	assert.Equal(
		t,
		map[string]any{
			"config": map[string]any{
				"address":             ":0",
				"base_url":            sse.DefaultBaseURL,
				"base_path":           sse.DefaultBasePath,
				"sse_endpoint":        sse.DefaultSSEEndpoint,
				"message_endpoint":    sse.DefaultMessageEndpoint,
				"keep_alive":          true,
				"keep_alive_interval": sse.DefaultKeepAliveInterval.Seconds(),
			},
			"status": map[string]any{
				"running": false,
			},
		},
		srv.Info(),
	)

	ctx := lg.WithContext(context.Background())

	go func(fCtx context.Context) {
		fErr := srv.Start(fCtx)
		assert.NoError(t, fErr)
	}(ctx)

	time.Sleep(1 * time.Millisecond)

	assert.True(t, srv.Running())

	logtest.AssertHasLogRecord(t, lb, map[string]any{
		"level":   "info",
		"message": "starting MCP SSE server on :0",
	})

	go func(fCtx context.Context) {
		fErr := srv.Stop(fCtx)
		assert.NoError(t, fErr)
	}(ctx)

	time.Sleep(1 * time.Millisecond)

	assert.False(t, srv.Running())

	logtest.AssertHasLogRecord(t, lb, map[string]any{
		"level":   "info",
		"message": "stopping MCP SSE server",
	})
}
