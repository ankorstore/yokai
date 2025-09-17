package fxclock_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/ankorstore/yokai/fxclock"
	"github.com/ankorstore/yokai/fxclock/testdata/service"
	"github.com/ankorstore/yokai/fxconfig"
	"github.com/jonboulle/clockwork"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
)

func TestFxClockworkClockModule(t *testing.T) {
	t.Setenv("APP_ENV", config.AppEnvDev)
	t.Setenv("APP_CONFIG_PATH", "testdata/config")

	runTest := func(tb testing.TB) (clockwork.Clock, *clockwork.FakeClock, *service.TestService) {
		tb.Helper()

		var clock clockwork.Clock
		var fakeClock *clockwork.FakeClock
		var srv *service.TestService

		app := fxtest.New(
			tb,
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxclock.FxClockModule,
			fx.Provide(service.NewTestService),
			fx.Populate(&clock, &fakeClock, &srv),
		)

		app.RequireStart().RequireStop()
		assert.NoError(tb, app.Err())

		return clock, fakeClock, srv
	}

	t.Run("normal mode", func(t *testing.T) {
		clock, fakeClock, srv := runTest(t)

		assert.NotNil(t, clock)
		assert.Implements(t, (*clockwork.Clock)(nil), clock)
		assert.Equal(t, "*clockwork.realClock", fmt.Sprintf("%T", clock))

		assert.Nil(t, fakeClock)

		assert.NotNil(t, srv)
	})

	t.Run("test mode with default time", func(t *testing.T) {
		t.Setenv("APP_ENV", config.AppEnvTest)

		clock, fakeClock, srv := runTest(t)
		assert.NotNil(t, clock)
		assert.Implements(t, (*clockwork.Clock)(nil), clock)
		assert.Equal(t, "*clockwork.FakeClock", fmt.Sprintf("%T", clock))

		assert.NotNil(t, fakeClock)
		assert.Implements(t, (*clockwork.Clock)(nil), fakeClock)
		assert.Equal(t, "*clockwork.FakeClock", fmt.Sprintf("%T", fakeClock))

		assert.NotNil(t, srv)

		startTime := srv.Now()
		fakeClock.Advance(10 * time.Minute)

		assert.Equal(t, startTime.Add(10*time.Minute), srv.Now())
	})

	t.Run("with test clock and fixed time", func(t *testing.T) {
		testTime := "2025-03-30T12:00:00Z"

		t.Setenv("APP_ENV", config.AppEnvTest)
		t.Setenv("MODULES_CLOCK_TEST_TIME", testTime)

		clock, fakeClock, srv := runTest(t)
		assert.NotNil(t, clock)
		assert.Implements(t, (*clockwork.Clock)(nil), clock)
		assert.Equal(t, "*clockwork.FakeClock", fmt.Sprintf("%T", clock))

		assert.NotNil(t, fakeClock)
		assert.Implements(t, (*clockwork.Clock)(nil), fakeClock)
		assert.Equal(t, "*clockwork.FakeClock", fmt.Sprintf("%T", fakeClock))

		assert.NotNil(t, srv)

		expectedTime, _ := time.Parse(time.RFC3339, testTime)
		assert.Equal(t, expectedTime, srv.Now())

		fakeClock.Advance(5 * time.Hour)
		assert.Equal(t, expectedTime.Add(5*time.Hour), srv.Now())

		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			srv.Sleep(3 * time.Second)
			wg.Done()
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := fakeClock.BlockUntilContext(ctx, 1)
		assert.NoError(t, err)

		fakeClock.Advance(10 * time.Second)
		wg.Wait()
	})

	t.Run("test mode with invalid time", func(t *testing.T) {
		testTime := "invalid"
		t.Setenv("APP_ENV", config.AppEnvTest)
		t.Setenv("MODULES_CLOCK_TEST_TIME", testTime)

		app := fx.New(
			fx.NopLogger,
			fxconfig.FxConfigModule,
			fxclock.FxClockModule,
			fx.Invoke(func(clock clockwork.Clock) {}),
		)

		err := app.Start(context.Background())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), fmt.Sprintf("cannot parse %q", testTime))
	})
}
