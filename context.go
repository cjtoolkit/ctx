package ctx

import (
	"context"
	"net/http"
)

const contextName = "context-e0881a4717c598b16eb965d396f1aff6"

/*
For storing anything related to user request
*/
type Context interface {
	Title() string
	SetTitle(title string)
	Data(name string) interface{}
	SetData(name string, value interface{})

	// The fn function only gets called if there is a cache miss.
	PersistData(name string, fn func() interface{}) interface{}
	Dep(name string) interface{}
	SetDep(name string, value interface{})

	// The fn function only gets called if there is a cache miss.
	PersistDep(name string, fn func() interface{}) interface{}
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
		data:  map[string]interface{}{},
		dep:   map[string]interface{}{},
		res:   res,
	}

	req = req.WithContext(context.WithValue(req.Context(), contextName, ctx))
	ctx.req = req

	return req, ctx
}

/*
Pulls out user context that was saved to the *http.Request.
*/
func GetContext(req *http.Request) Context {
	return req.Context().Value(contextName).(Context)
}

type contextHolder struct {
	title string
	data  map[string]interface{}
	dep   map[string]interface{}
	req   *http.Request
	res   http.ResponseWriter
}

func (c *contextHolder) Title() string                { return c.title }
func (c *contextHolder) SetTitle(title string)        { c.title = title }
func (c *contextHolder) Data(name string) interface{} { return c.data[name] }

func (c *contextHolder) SetData(name string, value interface{}) {
	_, found := c.data[name]
	panicOnFound(found)
	c.data[name] = value
}

func (c *contextHolder) PersistData(name string, fn func() interface{}) interface{} {
	return persist(c.data, name, fn)
}

func (c *contextHolder) Dep(name string) interface{} { return c.dep[name] }

func (c *contextHolder) SetDep(name string, value interface{}) {
	_, found := c.dep[name]
	panicOnFound(found)
	c.dep[name] = value
}

func (c *contextHolder) PersistDep(name string, fn func() interface{}) interface{} {
	return persist(c.dep, name, fn)
}

func (c *contextHolder) Request() *http.Request              { return c.req }
func (c *contextHolder) ResponseWriter() http.ResponseWriter { return c.res }
