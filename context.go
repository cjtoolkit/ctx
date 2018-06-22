package ctx

import (
	"context"
	"net/http"
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
	Request() *http.Request
	ResponseWriter() http.ResponseWriter
}

/*
Create new context for user request, also saves context inside *http.Request
without disturbing the context of the user request.

Avoid using this concurrently.
*/
func NewContext(res http.ResponseWriter, req *http.Request) (*http.Request, Context) {
	ctx := &contextHolder{
		title: "Untitled",
		data:  map[interface{}]interface{}{},
		dep:   map[interface{}]interface{}{},
		res:   res,
	}

	req = req.WithContext(context.WithValue(req.Context(), regContext{}, ctx))
	ctx.req = req

	return req, ctx
}

/*
Pulls out user context that was saved to the *http.Request.
*/
func GetContext(req *http.Request) Context {
	return req.Context().Value(regContext{}).(Context)
}

type contextHolder struct {
	title string
	data  map[interface{}]interface{}
	dep   map[interface{}]interface{}
	req   *http.Request
	res   http.ResponseWriter
}

func (c *contextHolder) Title() string                    { return c.title }
func (c *contextHolder) SetTitle(title string)            { c.title = title }
func (c *contextHolder) Data(key interface{}) interface{} { return c.data[key] }

func (c *contextHolder) SetData(key, value interface{}) {
	_, found := c.data[key]
	panicOnFound(found)
	c.data[key] = value
}

func (c *contextHolder) PersistData(key interface{}, fn func() interface{}) interface{} {
	return persist(c.data, key, fn)
}

func (c *contextHolder) Dep(key interface{}) interface{} { return c.dep[key] }

func (c *contextHolder) SetDep(key, value interface{}) {
	_, found := c.dep[key]
	panicOnFound(found)
	c.dep[key] = value
}

func (c *contextHolder) PersistDep(key interface{}, fn func() interface{}) interface{} {
	return persist(c.dep, key, fn)
}

func (c *contextHolder) Request() *http.Request              { return c.req }
func (c *contextHolder) ResponseWriter() http.ResponseWriter { return c.res }
