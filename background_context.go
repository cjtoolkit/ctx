package ctx

import (
	"time"
)

/*
Background Context
*/
type BackgroundContext interface {
	Set(name string, value interface{})
	Get(name string) interface{}

	// The fn function only gets called if there is a cache miss. Return error as nil to bypass health check.
	Persist(name string, fn func() (interface{}, error)) interface{}
}

/*
Create new background context.

Avoid using this concurrently.
*/
func NewBackgroundContext() BackgroundContext {
	return &backgroundContext{
		maxAttempt: 5,
		timeout:    2 * time.Second,
		ctx:        map[string]interface{}{},
	}
}

type backgroundContext struct {
	maxAttempt int
	timeout    time.Duration
	ctx        map[string]interface{}
}

func (bc *backgroundContext) Set(name string, value interface{}) {
	_, found := bc.ctx[name]
	panicOnFound(found)
	bc.ctx[name] = value
}

func (bc *backgroundContext) Get(name string) interface{} { return bc.ctx[name] }

func (bc *backgroundContext) Persist(name string, fn func() (interface{}, error)) interface{} {
	return persistWithHealthCheck(bc.maxAttempt, bc.timeout, bc.ctx, name, fn)
}

/*
Clear Background Context, best used with a defer function after creating the new context.
*/
func ClearBackgroundContext(context BackgroundContext) {
	context.(*backgroundContext).ctx = nil
}

/*
Set Health Check Max Attempt on Background Context
*/
func SetHealthCheckMaxAttemptOnBackgroundContext(context BackgroundContext, maxAttempt int) {
	context.(*backgroundContext).maxAttempt = maxAttempt
}

/*
Set Health Check Time Out on Background Context
*/
func SetHealthCheckTimeOutOnBackgroundContext(context BackgroundContext, timeout time.Duration) {
	context.(*backgroundContext).timeout = timeout
}
