package stdio_test

import (
	"context"
	"testing"
	"time"

	"github.com/ankorstore/yokai/fxmcpserver/server/stdio"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
)

func TestMCPStdioServer(t *testing.T) {
	t.Parallel()

	lb := logtest.NewDefaultTestLogBuffer()
	lg, err := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(lb))
	assert.NoError(t, err)

	mcpSrv := server.NewMCPServer("test-server", "1.0.0")

	srv := stdio.NewDefaultMCPStdioServerFactory().Create(mcpSrv)

	assert.False(t, srv.Running())

	assert.Equal(
		t,
		map[string]any{
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
		"message": "starting MCP Stdio server",
	})
}
