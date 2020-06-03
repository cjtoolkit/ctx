package ctx

import (
	goContext "context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/cjtoolkit/ctx/v2/internal"
)

type Context interface {
	// Set content to map, cannot use the same key value over and over again, otherwise it will panic.
	Set(key, value interface{})

	// Get content from map, return two values, the bool value indicated if it was found or not.
	Get(key interface{}) (interface{}, bool)

	// The fn function only gets called if there is a cache miss. Return error as nil to bypass health check.
	Persist(key interface{}, fn func() (interface{}, error)) interface{}
}

type contextBase struct {
	mutex  *sync.Mutex
	ctxMap map[interface{}]interface{}
}

/*
New Context, default max attempt is 1 and timeout is 0
*/
func NewContext(context goContext.Context) Context {
	return &contextBase{
		mutex: &sync.Mutex{},
		ctxMap: map[interface{}]interface{}{
			internal.GoContextKey{}: context,
		},
	}
}

/*
New Context with Map, default max attempt is 1 and timeout is 0
*/
func NewContextWithMap(m map[interface{}]interface{}) Context {
	return &contextBase{
		mutex:  &sync.Mutex{},
		ctxMap: m,
	}
}

/*
New Background Context, default max attempt is 5 and timeout is 2 seconds
*/
func NewBackgroundContext() Context {
	return &contextBase{
		mutex: &sync.Mutex{},
		ctxMap: map[interface{}]interface{}{
			internal.GoContextKey{}: goContext.Background(),
			settingsKey{}: &Settings{
				MaxAttempt: 5,
				Timeout:    2 * time.Second,
			},
		},
	}
}

func (c *contextBase) Set(key, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, found := c.get(key)
	internal.PanicOnFound(found)
	c.set(key, value)
}

func (c *contextBase) set(key, value interface{}) {
	c.ctxMap[key] = value
}

func (c *contextBase) Get(key interface{}) (interface{}, bool) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.get(key)
}

func (c *contextBase) get(key interface{}) (interface{}, bool) {
	v, found := c.ctxMap[key]
	return v, found
}

func (c *contextBase) Persist(key interface{}, fn func() (interface{}, error)) interface{} {
	c.mutex.Lock()
	v, found := c.get(key)
	if found {
		c.mutex.Unlock()
		return internal.CheckForLockOrReturnValue(v)
	}
	c.set(key, internal.Lock{})
	settings := c.settings()
	c.mutex.Unlock()

	attempt := 0

	for {
		v, err := fn()
		if nil == err {
			c.mutex.Lock()
			c.set(key, v)
			c.mutex.Unlock()
			return v
		}

		attempt++
		if attempt >= settings.MaxAttempt {
			log.Panic(err)
		}
		time.Sleep(settings.Timeout)
	}

	return nil
}

/*
Get Go Context
*/
func GoContext(context Context) goContext.Context {
	v, found := context.Get(internal.GoContextKey{})
	if !found {
		panic(errors.New("go context is not found"))
	}
	return v.(goContext.Context)
}

/*
Get Go Context

Deprecated: Use GoContext function instead. Will be removed in 3.0
*/
func GetGoContext(context Context) goContext.Context { return GoContext(context) }

/*
Clear Background Context, best used with a defer function after creating the new context.
*/
func ClearContext(context Context) {
	if context, ok := context.(*contextBase); ok {
		context.mutex.Lock()
		context.ctxMap = nil
		context.mutex.Unlock()
	}
}

/*
Set Health Check Max Attempt on Background Context

Deprecated: Use ContextSettings instead. Will be removed in 3.0
*/
func SetMaxAttempt(context Context, maxAttempt int) {
	ContextSettings(context).MaxAttempt = maxAttempt
}

/*
Set Health Check Time Out on Background Context

Deprecated: Use ContextSettings instead. Will be remove in 3.0
*/
func SetTimeout(context Context, timeout time.Duration) {
	ContextSettings(context).Timeout = timeout
}
