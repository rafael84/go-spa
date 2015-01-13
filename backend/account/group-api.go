package account

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gotk/ctx"

	"github.com/rafael84/go-spa/backend/base"
)

func init() {
	ctx.Resource("/account/group", &GroupResource{}, false)
	ctx.Resource("/account/group/{id:[0-9]+}", &GroupItemResource{}, false)
}

type GroupResource struct {
	*base.Resource
}

func (r *GroupResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	groups, err := r.DB(c).FindAll(&Group{}, "")
	if err != nil {
		log.Errorf("Could not query groups: %v", err)
		return ctx.BadRequest(rw, "Could not query groups")
	}
	return ctx.OK(rw, groups)
}

type GroupItemResource struct {
	*base.Resource
}

func (r *GroupItemResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	group, err := r.DB(c).FindOne(&Group{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return ctx.BadRequest(rw, "Could not query group")
	}
	return ctx.OK(rw, group)
}

func (r *GroupItemResource) DELETE(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	group, err := r.DB(c).FindOne(&Group{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return ctx.BadRequest(rw, "Could not query group")
	}
	err = r.DB(c).Delete(group)
	if err != nil {
		log.Errorf("Could not delete group %s: %v", id, err)
		return ctx.InternalServerError(rw, "Could not delete user")
	}
	return ctx.NoContent(rw)
}
