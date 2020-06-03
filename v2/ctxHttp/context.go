package ctxHttp

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/cjtoolkit/ctx/v2"
	"github.com/cjtoolkit/ctx/v2/internal"
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
	newReq := &http.Request{}
	_context := ctx.NewContextWithMap(map[interface{}]interface{}{
		internal.RequestKey{}:         newReq,
		internal.OriginalRequestKey{}: req,
		internal.ResponseKey{}:        res,
		internal.GoContextKey{}:       req.Context(),
		internal.TitleKey{}:           &title{Title: "Untitled"},
	})

	{
		reqWithContext := req.WithContext(context.WithValue(req.Context(), httpKey{}, _context))
		*newReq = *reqWithContext
	}

	return newReq
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

/*
Pulls out original request from context.
*/
func OriginalRequest(context ctx.Context) *http.Request {
	v, found := context.Get(internal.OriginalRequestKey{})
	if !found {
		panic(errors.New("request is not found"))
	}
	return v.(*http.Request)
}

/*
Pulls out response from context
*/
func Response(context ctx.Context) http.ResponseWriter {
	v, found := context.Get(internal.ResponseKey{})
	if !found {
		panic(errors.New("response is not found"))
	}
	return v.(http.ResponseWriter)
}

/*
Alias of Response
*/
func ResponseWriter(context ctx.Context) http.ResponseWriter { return Response(context) }

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

/*
Set title from inside context.
*/
func SetTitle(context ctx.Context, titleStr string) {
	v, found := context.Get(internal.TitleKey{})
	if !found {
		return
	}
	v.(*title).Title = titleStr
}

/*
Execute Http Handler
*/
func ExecuteHttpHandler(context ctx.Context, handler http.Handler) {
	handler.ServeHTTP(Response(context), Request(context))
}

/*
Execute Http Handler Function
*/
func ExecuteHttpHandlerFunc(context ctx.Context, handlerFunc http.HandlerFunc) {
	ExecuteHttpHandler(context, handlerFunc)
}

/*
Get Parsed Url Query
*/
func UrlQuery(context ctx.Context) url.Values {
	type key struct{}
	return context.Persist(key{}, func() (interface{}, error) {
		return url.ParseQuery(Request(context).URL.RawQuery)
	}).(url.Values)
}
