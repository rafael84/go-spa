package base

import (
	"github.com/gotk/pg"
	"github.com/gotk/webctx"
)

type Resource struct{}

func (r *Resource) DB(c *webctx.Context) *pg.Session {
	return c.Vars["db"].(*pg.Session)
}
