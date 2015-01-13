package base

import (
	"github.com/gotk/ctx"
	"github.com/gotk/pg"
)

type Resource struct{}

func (r *Resource) DB(c *ctx.Context) *pg.Session {
	return c.Vars["db"].(*pg.Session)
}
