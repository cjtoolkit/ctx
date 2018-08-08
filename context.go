package ctx

import (
	"context"
	"net/http"
	"sync"
)

type regContext struct{}

/*
For storing anything related to user request
*/
type Context interface {
	Title() string
	SetTitle(title string)
	Data(key interface{}) interface{}
	SetData(key, value interface{})

	// The fn function only gets called if there is a cache miss.
	PersistData(key interface{}, fn func() interface{}) interface{}
	Dep(key interface{}) interface{}
	SetDep(key, value interface{})

	// The fn function only gets called if there is a cache miss.
	PersistDep(key interface{}, fn func() interface{}) interface{}

	Ctx() context.Context
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
}

/*
Create new context for user request, also saves context inside *http.Request
without disturbing the context of the user request.
*/
func NewContext(res http.ResponseWriter, req *http.Request) (*http.Request, Context) {
	ctx := &contextHolder{
		title: "Untitled",
		data:  &sync.Map{},
		dep:   &sync.Map{},
		res:   res,
	}

	req = req.WithContext(context.WithValue(req.Context(), regContext{}, ctx))
	ctx.req = req
	ctx.ctx = func() context.Context { return req.Context() }

	return req, ctx
}

/*
Create new context by context.
*/
func NewContextByContext(ctx context.Context) Context {
	ctxH := &contextHolder{
		title: "Untitled",
		data:  &sync.Map{},
		dep:   &sync.Map{},
	}
	ctxH.ctx = func() context.Context { return ctx }

	return ctxH
}

/*
Pulls out user context that was saved to the *http.Request.
*/
func GetContext(req *http.Request) Context {
	return req.Context().Value(regContext{}).(Context)
}

type contextHolder struct {
	rw    sync.RWMutex
	title string
	data  *sync.Map
	dep   *sync.Map
	ctx   func() context.Context
	req   *http.Request
	res   http.ResponseWriter
}

func (c *contextHolder) Title() string {
	c.rw.RLock()
	title := c.title
	c.rw.RUnlock()

	return title
}

func (c *contextHolder) SetTitle(title string) {
	c.rw.Lock()
	c.title = title
	c.rw.Unlock()
}

func (c *contextHolder) Data(key interface{}) interface{} {
	data, _ := c.data.Load(key)
	return data
}

func (c *contextHolder) SetData(key, value interface{}) {
	_, found := c.data.Load(key)
	panicOnFound(found)
	c.data.Store(key, value)
}

func (c *contextHolder) PersistData(key interface{}, fn func() interface{}) interface{} {
	return persist(c.data, key, fn)
}

func (c *contextHolder) Dep(key interface{}) interface{} {
	dep, _ := c.dep.Load(key)
	return dep
}

func (c *contextHolder) SetDep(key, value interface{}) {
	_, found := c.dep.Load(key)
	panicOnFound(found)
	c.dep.Store(key, value)
}

func (c *contextHolder) PersistDep(key interface{}, fn func() interface{}) interface{} {
	return persist(c.dep, key, fn)
}

func (c *contextHolder) Ctx() context.Context                { return c.ctx() }
func (c *contextHolder) Request() *http.Request              { return c.req }
func (c *contextHolder) ResponseWriter() http.ResponseWriter { return c.res }
