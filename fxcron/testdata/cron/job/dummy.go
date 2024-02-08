package job

import (
	"context"
)

type DummyCron struct{}

func NewDummyCron() *DummyCron {
	return &DummyCron{}
}

func (c *DummyCron) Name() string {
	return "dummy"
}

func (c *DummyCron) Run(ctx context.Context) error {
	return nil
}
