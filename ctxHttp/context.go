package ctxHttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/cjtoolkit/ctx"
	"github.com/cjtoolkit/ctx/internal"
)

type httpKey struct{}

type title struct {
	Title string
}

/*
Create new context for user request, also saves context inside *http.Request
without disturbing the context of the user request.
*/
func NewContext(req *http.Request, res http.ResponseWriter) *http.Request {
	_context := ctx.NewContextWithMap(map[interface{}]interface{}{
		internal.RequestKey{}:   req,
		internal.ResponseKey{}:  res,
		internal.GoContextKey{}: req.Context(),
		internal.TitleKey{}:     &title{Title: "Untitled"},
	})

	return req.WithContext(context.WithValue(req.Context(), httpKey{}, _context))
}

/*
Pulls out user context that was saved to the *http.Request.
*/
func Context(req *http.Request) ctx.Context {
	return req.Context().Value(httpKey{}).(ctx.Context)
}

/*
Pulls out request from context.
*/
func Request(context ctx.Context) *http.Request {
	v, found := context.Get(internal.RequestKey{})
	if !found {
		panic(errors.New("request is not found"))
	}
	return v.(*http.Request)
}

/**
Pulls out response from context
*/
func Response(context ctx.Context) http.ResponseWriter {
	v, found := context.Get(internal.ResponseKey{})
	if !found {
		panic(errors.New("response is not found"))
	}
	return v.(http.ResponseWriter)
}

/**
Pulls out title from context
*/
func Title(context ctx.Context) string {
	v, found := context.Get(internal.TitleKey{})
	if !found {
		return ""
	}
	return v.(*title).Title
}

/**
Set title from inside context.
*/
func SetTitle(context ctx.Context, titleStr string) {
	v, found := context.Get(internal.TitleKey{})
	if !found {
		return
	}
	v.(*title).Title = titleStr
}
