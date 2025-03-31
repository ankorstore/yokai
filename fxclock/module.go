package fxclock

import (
	"time"

	"github.com/ankorstore/yokai/config"
	"github.com/jonboulle/clockwork"
	"go.uber.org/fx"
)

// ModuleName is the module name.
const ModuleName = "clock"

// FxClockModule is the [Fx] clockwork module.
//
// [Fx]: https://github.com/uber-go/fx
var FxClockModule = fx.Module(
	ModuleName,
	fx.Provide(
		NewFxClock,
	),
)

// FxClockParam allows injection of the required dependencies in [NewFxClock].
type FxClockParam struct {
	fx.In
	Config *config.Config
}

// NewFxClock returns a new [clockwork.Clock] instance.
func NewFxClock(p FxClockParam) (clockwork.Clock, *clockwork.FakeClock, error) {
	if p.Config.IsTestEnv() {
		testTimeCfg := p.Config.GetString("modules.clock.test.time")
		if testTimeCfg == "" {
			testClock := clockwork.NewFakeClock()

			return testClock, testClock, nil
		}

		testTime, err := time.Parse(time.RFC3339, testTimeCfg)
		if err != nil {
			return nil, nil, err
		}

		testClock := clockwork.NewFakeClockAt(testTime)

		return testClock, testClock, nil
	}

	return clockwork.NewRealClock(), nil, nil
}
