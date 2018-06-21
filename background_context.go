package ctx

import (
	"time"
)

/*
Background Context
*/
type BackgroundContext interface {
	Set(key interface{}, value interface{})
	Get(key interface{}) interface{}

	// The fn function only gets called if there is a cache miss. Return error as nil to bypass health check.
	Persist(key interface{}, fn func() (interface{}, error)) interface{}
}

/*
Create new background context.

Avoid using this concurrently.
*/
func NewBackgroundContext() BackgroundContext {
	return &backgroundContext{
		maxAttempt: 5,
		timeout:    2 * time.Second,
		ctx:        map[interface{}]interface{}{},
	}
}

type backgroundContext struct {
	maxAttempt int
	timeout    time.Duration
	ctx        map[interface{}]interface{}
}

func (bc *backgroundContext) Set(key interface{}, value interface{}) {
	_, found := bc.ctx[key]
	panicOnFound(found)
	bc.ctx[key] = value
}

func (bc *backgroundContext) Get(key interface{}) interface{} { return bc.ctx[key] }

func (bc *backgroundContext) Persist(key interface{}, fn func() (interface{}, error)) interface{} {
	return persistWithHealthCheck(bc.maxAttempt, bc.timeout, bc.ctx, key, fn)
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
