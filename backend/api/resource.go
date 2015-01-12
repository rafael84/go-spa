package api

import (
	"net/http"

	"github.com/rafael84/go-spa/backend/context"
)

type HttpGetter interface {
	GET(c *context.Context, rw http.ResponseWriter, req *http.Request) error
}

type HttpPoster interface {
	POST(c *context.Context, rw http.ResponseWriter, req *http.Request) error
}

type HttpPutter interface {
	PUT(c *context.Context, rw http.ResponseWriter, req *http.Request) error
}

type HttpDeleter interface {
	DELETE(c *context.Context, rw http.ResponseWriter, req *http.Request) error
}

type resource struct {
	context.Endpoint
}

func AddResource(resource *resource) {
	endpoints = append(endpoints, &resource.Endpoint)
}

func NewResource(path string) *resource {
	return &resource{context.Endpoint{
		Path:     path,
		Handlers: make(context.MethodHandlers, 0),
	}}
}

func Resource(path string, res interface{}) {
	handlers := make(context.MethodHandlers, 0)

	if _, ok := res.(HttpGetter); ok {
		handlers["GET"] = res.(HttpGetter).GET
	}

	if _, ok := res.(HttpPutter); ok {
		handlers["PUT"] = res.(HttpPutter).PUT
	}

	AddResource(&resource{context.Endpoint{
		Path:     path,
		Handlers: handlers,
	}})
}

func (r *resource) GET(handler context.ContextHandler) *resource {
	r.Handlers["GET"] = handler
	return r
}

func (r *resource) POST(handler context.ContextHandler) *resource {
	r.Handlers["POST"] = handler
	return r
}

func (r *resource) PUT(handler context.ContextHandler) *resource {
	r.Handlers["PUT"] = handler
	return r
}

func (r *resource) DELETE(handler context.ContextHandler) *resource {
	r.Handlers["DELETE"] = handler
	return r
}
