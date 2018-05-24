package ctx

import (
	"log"
	"time"
)

/*
Background Context
*/
type BackgroundContext interface {
	Set(name string, value interface{})
	Get(name string) interface{}
	Persist(name string, fn func() (interface{}, error)) interface{}
}

func NewBackgroundContext() BackgroundContext {
	return &backgroundContext{
		maxAttempt: 5,
		timeout:    1 * time.Second,
		ctx:        map[string]interface{}{},
	}
}

type backgroundContext struct {
	maxAttempt int
	timeout    time.Duration
	ctx        map[string]interface{}
}

func (bc *backgroundContext) Set(name string, value interface{}) { bc.ctx[name] = value }
func (bc *backgroundContext) Get(name string) interface{}        { return bc.ctx[name] }

func (bc *backgroundContext) Persist(name string, fn func() (interface{}, error)) interface{} {
	return persistWithHealthCheck(bc.maxAttempt, bc.timeout, bc.ctx, name, fn, func(err error) {
		log.Fatal(err)
	})
}
