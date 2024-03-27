package fxtrace

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ankorstore/yokai/config"
)

func BuildOtlpGrpcDialRetryPolicy(config *config.Config) string {
	// max attempts
	maxAttempts := 4
	maxAttemptsConfig := config.GetInt("modules.trace.processor.options.retry.max_attempts")
	if maxAttemptsConfig != 0 {
		maxAttempts = maxAttemptsConfig
	}

	// initial backoff
	initialBackoff := 0.1
	initialBackoffConfig := config.GetFloat64("modules.trace.processor.options.retry.initial_backoff")
	if initialBackoffConfig != 0 {
		initialBackoff = initialBackoffConfig
	}

	// max backoff
	maxBackoff := 1.0
	maxBackoffConfig := config.GetFloat64("modules.trace.processor.options.retry.max_backoff")
	if maxBackoffConfig != 0 {
		maxBackoff = maxBackoffConfig
	}

	// backoff multiplier
	backoffMultiplier := 2
	backoffMultiplierConfig := config.GetInt("modules.trace.processor.options.retry.backoff_multiplier")
	if backoffMultiplierConfig != 0 {
		backoffMultiplier = backoffMultiplierConfig
	}

	retryableStatusCodes := []string{"UNAVAILABLE"}
	retryableStatusCodesConfig := config.GetStringSlice("modules.trace.processor.options.retry.retryable_status_codes")
	if len(retryableStatusCodesConfig) > 0 {
		retryableStatusCodes = retryableStatusCodesConfig
	}

	retryPolicy := fmt.Sprintf(
		`{
            "methodConfig": [{
                "waitForReady": true,
                "retryPolicy": {
                    "MaxAttempts": %d,
                    "InitialBackoff": "%ss",
                    "MaxBackoff": "%ss",
                    "BackoffMultiplier": %d,
                    "RetryableStatusCodes": [ "%s" ]
                }
            }]
        }`,
		maxAttempts,
		strconv.FormatFloat(initialBackoff, 'f', -1, 64),
		strconv.FormatFloat(maxBackoff, 'f', -1, 64),
		backoffMultiplier,
		strings.Join(retryableStatusCodes, `", "`),
	)

	return retryPolicy
}
