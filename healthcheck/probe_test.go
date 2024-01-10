package healthcheck_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ankorstore/yokai/healthcheck"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type CheckerProbeMock struct {
	mock.Mock
}

func (m *CheckerProbeMock) Name() string {
	args := m.Called()

	return args.String(0)
}

func (m *CheckerProbeMock) Check(ctx context.Context) *healthcheck.CheckerProbeResult {
	args := m.Called(ctx)

	//nolint:forcetypeassert
	return args.Get(0).(*healthcheck.CheckerProbeResult)
}

func TestNewCheckerProbeResult(t *testing.T) {
	t.Parallel()

	successResult := healthcheck.NewCheckerProbeResult(true, "success")

	assert.True(t, successResult.Success)
	assert.Equal(t, "success", successResult.Message)

	failureResult := healthcheck.NewCheckerProbeResult(false, "failure")

	assert.False(t, failureResult.Success)
	assert.Equal(t, "failure", failureResult.Message)
}

func TestCheckerProbeResultAsJson(t *testing.T) {
	t.Parallel()

	result := healthcheck.NewCheckerProbeResult(true, "success")

	data, err := json.Marshal(result)

	assert.Nil(t, err)
	assert.Equal(t, `{"success":true,"message":"success"}`, string(data))
}
