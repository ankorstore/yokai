package service

import (
	"github.com/ankorstore/yokai/config"
)

type TestServiceDependency struct {
	config *config.Config
}

func NewTestServiceDependency(config *config.Config) *TestServiceDependency {
	return &TestServiceDependency{
		config: config,
	}
}

func (s *TestServiceDependency) AppName() string {
	return s.config.AppName()
}
