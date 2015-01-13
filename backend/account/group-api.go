package account

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gotk/webctx"

	"github.com/rafael84/go-spa/backend/base"
)

func init() {
	webctx.Resource("/account/group", &GroupResource{}, false)
	webctx.Resource("/account/group/{id:[0-9]+}", &GroupItemResource{}, false)
}

type GroupResource struct {
	*base.Resource
}

func (r *GroupResource) GET(c *webctx.Context, rw http.ResponseWriter, req *http.Request) error {
	groups, err := r.DB(c).Filter(&Group{}, "")
	if err != nil {
		log.Errorf("Could not query groups: %v", err)
		return webctx.BadRequest(rw, "Could not query groups")
	}
	return webctx.OK(rw, groups)
}

type GroupItemResource struct {
	*base.Resource
}

func (r *GroupItemResource) GET(c *webctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	group, err := r.DB(c).One(&Group{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return webctx.BadRequest(rw, "Could not query group")
	}
	return webctx.OK(rw, group)
}

func (r *GroupItemResource) DELETE(c *webctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	group, err := r.DB(c).One(&Group{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return webctx.BadRequest(rw, "Could not query group")
	}
	err = r.DB(c).Delete(group)
	if err != nil {
		log.Errorf("Could not delete group %s: %v", id, err)
		return webctx.InternalServerError(rw, "Could not delete user")
	}
	return webctx.NoContent(rw)
}
