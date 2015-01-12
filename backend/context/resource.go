package context

import "net/http"

type HttpGetter interface {
	GET(c *Context, rw http.ResponseWriter, req *http.Request) error
}

type HttpPoster interface {
	POST(c *Context, rw http.ResponseWriter, req *http.Request) error
}

type HttpPutter interface {
	PUT(c *Context, rw http.ResponseWriter, req *http.Request) error
}

type HttpDeleter interface {
	DELETE(c *Context, rw http.ResponseWriter, req *http.Request) error
}

func Resource(path string, res interface{}, public bool) {
	handlers := make(MethodHandlers, 0)

	if _, ok := res.(HttpGetter); ok {
		handlers["GET"] = res.(HttpGetter).GET
	}

	if _, ok := res.(HttpPutter); ok {
		handlers["PUT"] = res.(HttpPutter).PUT
	}

	if _, ok := res.(HttpPoster); ok {
		handlers["POST"] = res.(HttpPoster).POST
	}

	if _, ok := res.(HttpDeleter); ok {
		handlers["DELETE"] = res.(HttpDeleter).DELETE
	}

	AddEndpoint(&Endpoint{
		Path:     path,
		Public:   public,
		Handlers: handlers,
	})
}

func (e *Endpoint) GET(handler ContextHandler) *Endpoint {
	e.Handlers["GET"] = handler
	return e
}

func (e *Endpoint) POST(handler ContextHandler) *Endpoint {
	e.Handlers["POST"] = handler
	return e
}

func (e *Endpoint) PUT(handler ContextHandler) *Endpoint {
	e.Handlers["PUT"] = handler
	return e
}

func (e *Endpoint) DELETE(handler ContextHandler) *Endpoint {
	e.Handlers["DELETE"] = handler
	return e
}
