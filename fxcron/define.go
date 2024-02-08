package fxcron

import (
	"github.com/go-co-op/gocron/v2"
)

type CronJobDefinition interface {
	ReturnType() string
	Expression() string
	Options() []gocron.JobOption
}

type cronJobDefinition struct {
	returnType string
	expression string
	options    []gocron.JobOption
}

func NewCronJobDefinition(returnType string, expression string, options ...gocron.JobOption) CronJobDefinition {
	return &cronJobDefinition{
		returnType: returnType,
		expression: expression,
		options:    options,
	}
}

func (c *cronJobDefinition) ReturnType() string {
	return c.returnType
}

func (c *cronJobDefinition) Expression() string {
	return c.expression
}

func (c *cronJobDefinition) Options() []gocron.JobOption {
	return c.options
}
