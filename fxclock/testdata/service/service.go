package service

import (
	"time"

	"github.com/jonboulle/clockwork"
)

type TestService struct {
	clock clockwork.Clock
}

func NewTestService(clock clockwork.Clock) *TestService {
	return &TestService{clock: clock}
}

func (s *TestService) Now() time.Time {
	return s.clock.Now()
}

func (s *TestService) Sleep(d time.Duration) {
	s.clock.Sleep(d)
}
