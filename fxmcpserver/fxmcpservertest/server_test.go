package fxmcpservertest_test

import (
	"context"
	"testing"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxmcpserver/fxmcpservertest"
	"github.com/ankorstore/yokai/fxmcpserver/server/sse"
	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
	"github.com/ankorstore/yokai/log/logtest"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/sdk/trace"
)

func TestMCPSSETestServer(t *testing.T) {
	t.Parallel()

	cfg, err := config.NewDefaultConfigFactory().Create(
		config.WithFilePaths("../testdata/config"),
	)
	assert.NoError(t, err)

	gm := uuid.NewDefaultUuidGenerator()

	tp := trace.NewTracerProvider()

	lb := logtest.NewDefaultTestLogBuffer()
	lg, err := log.NewDefaultLoggerFactory().Create(log.WithOutputWriter(lb))
	assert.NoError(t, err)

	hdl := sse.NewDefaultMCPSSEServerContextHandler(gm, tp, lg)

	mcpSrv := server.NewMCPServer("test-server", "1.0.0")

	srv := fxmcpservertest.NewMCPSSETestServer(cfg, mcpSrv, hdl)
	defer srv.Close()

	cli, err := srv.StartClient(context.Background())
	assert.NoError(t, err)

	err = cli.Ping(context.Background())
	assert.NoError(t, err)

	err = cli.Close()
	assert.NoError(t, err)
}
