package worker

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ankorstore/yokai/generate/uuid"
	"github.com/ankorstore/yokai/log"
)

const (
	LogRecordFieldWorkerName            = "worker"
	LogRecordFieldWorkerExecutionId     = "workerExecutionID"
	TraceSpanAttributeWorkerName        = "Worker"
	TraceSpanAttributeWorkerExecutionId = "WorkerExecutionID"
)

// WorkerPool is the [Worker] pool.
//
//nolint:containedctx
type WorkerPool struct {
	mutex             sync.Mutex
	waitGroup         sync.WaitGroup
	context           context.Context
	contextCancelFunc context.CancelFunc
	options           PoolOptions
	metrics           *WorkerMetrics
	generator         uuid.UuidGenerator
	registrations     map[string]*WorkerRegistration
	executions        map[string]*WorkerExecution
}

// NewWorkerPool returns a new [WorkerPool], with optional [WorkerPoolOption].
func NewWorkerPool(options ...WorkerPoolOption) *WorkerPool {
	poolOptions := DefaultWorkerPoolOptions()
	for _, opt := range options {
		opt(&poolOptions)
	}

	return &WorkerPool{
		options:       poolOptions,
		metrics:       poolOptions.Metrics,
		generator:     poolOptions.Generator,
		registrations: poolOptions.Registrations,
		executions:    make(map[string]*WorkerExecution),
	}
}

// Register registers a new [WorkerRegistration] onto the [WorkerPool].
func (p *WorkerPool) Register(registrations ...*WorkerRegistration) *WorkerPool {
	for _, registration := range registrations {
		p.registrations[registration.Worker().Name()] = registration
	}

	return p
}

// Start starts all [Worker] registered in the [WorkerPool].
func (p *WorkerPool) Start(ctx context.Context) error {
	p.context, p.contextCancelFunc = context.WithCancel(ctx)

	for _, worker := range p.registrations {
		//nolint:contextcheck
		p.startWorkerRegistration(p.context, worker)
	}

	return nil
}

// Stop gracefully stops all [Worker] registered in the [WorkerPool].
func (p *WorkerPool) Stop() error {
	p.contextCancelFunc()

	p.waitGroup.Wait()

	return nil
}

// Options returns the list of [PoolOptions] of the [WorkerPool].
func (p *WorkerPool) Options() PoolOptions {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.options
}

// Metrics returns the [WorkerPool] internal [WorkerMetrics].
func (p *WorkerPool) Metrics() *WorkerMetrics {
	return p.metrics
}

// Registrations returns the [WorkerPool] list of [WorkerRegistration].
func (p *WorkerPool) Registrations() map[string]*WorkerRegistration {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.registrations
}

// Registration returns the [WorkerRegistration] from the [WorkerPool] for a given worker name.
func (p *WorkerPool) Registration(name string) (*WorkerRegistration, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if registration, ok := p.registrations[name]; ok {
		return registration, nil
	}

	return nil, fmt.Errorf("registration for worker %s was not found", name)
}

// Executions returns the [WorkerPool] list of [WorkerExecution].
func (p *WorkerPool) Executions() map[string]*WorkerExecution {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return p.executions
}

// Execution returns the [WorkerExecution] from the [WorkerPool] for a given worker name.
func (p *WorkerPool) Execution(name string) (*WorkerExecution, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if execution, ok := p.executions[name]; ok {
		return execution, nil
	}

	return nil, fmt.Errorf("execution for worker %s was not found", name)
}

func (p *WorkerPool) startWorkerRegistration(ctx context.Context, registration *WorkerRegistration) {
	p.waitGroup.Add(1)

	execution := p.retrieveWorkerRegistrationExecution(registration)

	executionLogger := log.FromZerolog(
		log.CtxLogger(ctx).
			ToZerolog().
			With().
			Str(LogRecordFieldWorkerName, registration.Worker().Name()).
			Str(LogRecordFieldWorkerExecutionId, execution.Id()).
			Logger(),
	)

	executionCtx := context.WithValue(ctx, CtxWorkerNameKey{}, registration.Worker().Name())
	executionCtx = context.WithValue(executionCtx, CtxWorkerExecutionIdKey{}, execution.Id())
	executionCtx = executionLogger.WithContext(executionCtx)

	go func(ctx context.Context, workerExecution *WorkerExecution) {
		defer func() {
			p.waitGroup.Done()

			if r := recover(); r != nil {
				message := fmt.Sprintf(
					"stopping execution attempt %d/%d with recovered panic: %s",
					workerExecution.CurrentExecutionAttempt(),
					workerExecution.MaxExecutionsAttempts(),
					r,
				)

				executionLogger.Error().Msg(message)

				workerExecution.SetStatus(Error).AddEvent(message)

				p.metrics.IncrementWorkerExecutionError(registration.Worker().Name())

				if workerExecution.CurrentExecutionAttempt() < workerExecution.MaxExecutionsAttempts() {
					message = "restarting after panic recovery"

					executionLogger.Info().Msg(message)

					workerExecution.AddEvent(message).SetId(p.generator.Generate())

					p.metrics.IncrementWorkerExecutionRestart(registration.Worker().Name())

					p.startWorkerRegistration(ctx, registration)
				} else {
					message = "max execution attempts reached"

					executionLogger.Info().Msg(message)

					workerExecution.AddEvent(message)
				}
			}
		}()

		if workerExecution.CurrentExecutionAttempt() == 0 && workerExecution.DeferredStartThreshold() > 0 {
			message := fmt.Sprintf(
				"deferring execution attempt for %g seconds",
				workerExecution.DeferredStartThreshold(),
			)

			executionLogger.Info().Msg(message)

			workerExecution.SetStatus(Deferred).AddEvent(message)

			time.Sleep(time.Duration(workerExecution.DeferredStartThreshold()) * time.Second)
		}

		workerExecution.SetCurrentExecutionAttempt(workerExecution.CurrentExecutionAttempt() + 1)

		message := fmt.Sprintf(
			"starting execution attempt %d/%d",
			workerExecution.CurrentExecutionAttempt(),
			workerExecution.MaxExecutionsAttempts(),
		)

		executionLogger.Info().Msg(message)

		workerExecution.SetStatus(Running).AddEvent(message)

		p.metrics.IncrementWorkerExecutionStart(registration.Worker().Name())

		if err := registration.Worker().Run(ctx); err != nil {
			message = fmt.Sprintf(
				"stopping execution attempt %d/%d with error: %v",
				workerExecution.CurrentExecutionAttempt(),
				workerExecution.MaxExecutionsAttempts(),
				err.Error(),
			)

			executionLogger.Error().Err(err).Msg(message)

			workerExecution.SetStatus(Error).AddEvent(message)

			p.metrics.IncrementWorkerExecutionError(registration.Worker().Name())

			if workerExecution.CurrentExecutionAttempt() < workerExecution.MaxExecutionsAttempts() {
				message = "restarting after error"

				executionLogger.Info().Msg(message)

				workerExecution.AddEvent(message).SetId(p.generator.Generate())

				p.metrics.IncrementWorkerExecutionRestart(registration.Worker().Name())

				p.startWorkerRegistration(ctx, registration)
			} else {
				message = "max execution attempts reached"

				executionLogger.Info().Msg(message)

				workerExecution.AddEvent(message)
			}
		} else {
			message = fmt.Sprintf(
				"stopping execution attempt %d/%d with success",
				workerExecution.CurrentExecutionAttempt(),
				workerExecution.MaxExecutionsAttempts(),
			)

			executionLogger.Info().Msg(message)

			workerExecution.SetStatus(Success).AddEvent(message)

			p.metrics.IncrementWorkerExecutionSuccess(registration.Worker().Name())
		}
	}(executionCtx, execution)
}

func (p *WorkerPool) retrieveWorkerRegistrationExecution(registration *WorkerRegistration) *WorkerExecution {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.executions[registration.Worker().Name()]; !ok {
		executionOptions := DefaultWorkerExecutionOptions()

		executionOptions.DeferredStartThreshold = p.options.GlobalDeferredStartThreshold
		executionOptions.MaxExecutionsAttempts = p.options.GlobalMaxExecutionsAttempts

		for _, opt := range registration.Options() {
			opt(&executionOptions)
		}

		p.executions[registration.Worker().Name()] = NewWorkerExecution(
			p.generator.Generate(),
			registration.Worker().Name(),
			executionOptions,
		)
	}

	return p.executions[registration.Worker().Name()]
}
