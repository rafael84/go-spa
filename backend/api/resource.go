package api

import "github.com/rafael84/go-spa/backend/context"

type resource struct {
	context.Endpoint
}

func NewResource(path string) *resource {
	return &resource{context.Endpoint{
		Path:     path,
		Handlers: make(context.MethodHandlers, 0),
	}}
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
