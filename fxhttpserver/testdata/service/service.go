package service

import (
	"github.com/ankorstore/yokai/config"
)

type TestService struct {
	config *config.Config
}

func NewTestService(config *config.Config) *TestService {
	return &TestService{
		config: config,
	}
}

func (s *TestService) GetAppName() string {
	return s.config.AppName()
}
