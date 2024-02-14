package fxcron

import (
	"github.com/go-co-op/gocron/v2"
)

// CronJobDefinition is the interface for cron job definitions.
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

// NewCronJobDefinition returns a new [CronJobDefinition].
func NewCronJobDefinition(returnType string, expression string, options ...gocron.JobOption) CronJobDefinition {
	return &cronJobDefinition{
		returnType: returnType,
		expression: expression,
		options:    options,
	}
}

// ReturnType returns the definition return type.
func (c *cronJobDefinition) ReturnType() string {
	return c.returnType
}

// Expression returns the definition cron expression.
func (c *cronJobDefinition) Expression() string {
	return c.expression
}

// Options returns the definition cron job options.
func (c *cronJobDefinition) Options() []gocron.JobOption {
	return c.options
}
