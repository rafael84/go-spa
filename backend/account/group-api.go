package account

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/rafael84/go-spa/backend/base"
	"github.com/rafael84/go-spa/backend/context"
)

func init() {
	context.Resource("/account/group", &GroupResource{}, false)
	context.Resource("/account/group/{id:[0-9]+}", &GroupItemResource{}, false)
}

type GroupResource struct {
	*base.Resource
}

func (r *GroupResource) GET(c *context.Context, rw http.ResponseWriter, req *http.Request) error {
	groups, err := r.DB(c).Filter(&Group{}, "")
	if err != nil {
		log.Errorf("Could not query groups: %v", err)
		return context.BadRequest(rw, "Could not query groups")
	}
	return context.OK(rw, groups)
}

type GroupItemResource struct {
	*base.Resource
}

func (r *GroupItemResource) GET(c *context.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	group, err := r.DB(c).One(&Group{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return context.BadRequest(rw, "Could not query group")
	}
	return context.OK(rw, group)
}

func (r *GroupItemResource) DELETE(c *context.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	group, err := r.DB(c).One(&Group{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return context.BadRequest(rw, "Could not query group")
	}
	err = r.DB(c).Delete(group)
	if err != nil {
		log.Errorf("Could not delete group %s: %v", id, err)
		return context.InternalServerError(rw, "Could not delete user")
	}
	return context.NoContent(rw)
}
