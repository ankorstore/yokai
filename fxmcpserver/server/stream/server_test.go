package stream_test

import (
	"context"
	"github.com/ankorstore/yokai/fxmcpserver/server/stream"
	"testing"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestMCPStreamableHTTPServer(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("../../testdata/config"),
	)
	assert.NoError(t, err)

	lb := logtest.NewDefaultTestLogBuffer()
	lg, err := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(lb))
	assert.NoError(t, err)

	mcpSrv := server.NewMCPServer("test-server", "1.0.0")

	srv := stream.NewDefaultMCPStreamableHTTPServerFactory(cfg).Create(mcpSrv)

	assert.False(t, srv.Running())

	assert.Equal(
		t,
		map[string]any{
			"config": map[string]any{
				"address":             ":0",
				"stateless":           true,
				"base_path":           stream.DefaultBasePath,
				"keep_alive":          true,
				"keep_alive_interval": stream.DefaultKeepAliveInterval.Seconds(),
			},
			"status": map[string]any{
				"running": false,
			},
		},
		srv.Info(),
	)

	ctx := lg.WithContext(context.Background())

	//nolint:errcheck
	go srv.Start(ctx)

	time.Sleep(1 * time.Millisecond)

	assert.True(t, srv.Running())

	logtest.AssertHasLogRecord(t, lb, map[string]any{
		"level":   "info",
		"message": "starting MCP StreamableHTTP server on :0",
	})

	err = srv.Stop(ctx)
	assert.NoError(t, err)

	time.Sleep(1 * time.Millisecond)

	assert.False(t, srv.Running())

	logtest.AssertHasLogRecord(t, lb, map[string]any{
		"level":   "info",
		"message": "stopping MCP StreamableHTTP server",
	})
}
