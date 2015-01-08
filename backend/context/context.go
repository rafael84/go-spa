package context

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/nicksnyder/go-i18n/i18n"

	"github.com/rafael84/go-spa/backend/database"
)

type Context struct {
	Router *mux.Router
	DB     *database.Session
	T      i18n.TranslateFunc
}

type Handler func(c *Context, rw http.ResponseWriter, req *http.Request) error

func (c *Context) updateT(req *http.Request) {
	acceptLang := req.Header.Get("Accept-Language")
	defaultLang := "en-US"
	c.T = i18n.MustTfunc(acceptLang, defaultLang)
}

func (c *Context) AddRoute(path string, handler Handler) {
	httpHandler := func(rw http.ResponseWriter, req *http.Request) {
		c.updateT(req)
		err := handler(c, rw, req)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
	}
	c.Router.HandleFunc(path, httpHandler)
}

func (c *Context) AddRoutes(routes map[string]Handler) {
	for path, handler := range routes {
		c.AddRoute(path, handler)
	}
}

func NewContext(router *mux.Router, db *database.Session) *Context {
	c := &Context{
		Router: router,
		DB:     db,
	}
	return c
}
