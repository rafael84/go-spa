package base

import (
	"github.com/gotk/webctx"

	"github.com/rafael84/go-spa/backend/database"
)

type Resource struct{}

func (r *Resource) DB(c *webctx.Context) *database.Session {
	return c.Vars["db"].(*database.Session)
}
