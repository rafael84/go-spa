package user

import (
	"net/http"

	"github.com/gotk/ctx"
)

func init() {
	ctx.Resource("/account/user/role", &Role{}, true)
}

type Role struct{}

func (r *Role) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	roles := []struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}{
		{Id: 0, Name: "Admin"},
		{Id: 1, Name: "User"},
		{Id: 2, Name: "Uploader"},
		{Id: 3, Name: "Read Only"},
	}
	return ctx.OK(rw, roles)
}
