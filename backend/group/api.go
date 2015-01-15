package group

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/base"
)

func init() {
	ctx.Resource("/group", &GroupResource{}, false)
	ctx.Resource("/group/{id:[0-9]+}", &GroupItemResource{}, false)
}

type GroupResource struct {
	*base.Resource
}

func (r *GroupResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	groups, err := r.DB(c).FindAll(&Group{}, "")
	if err != nil {
		log.Errorf("Could not query groups: %v", err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_query_groups"))
	}
	return ctx.OK(rw, groups)
}

func (r *GroupResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {

	// decode request data
	var form = &struct {
		Name string `json:"name"`
	}{}
	err := json.NewDecoder(req.Body).Decode(form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_parse_request_data"))
	}

	// create new group
	grp := &Group{Name: form.Name}
	err = r.DB(c).Create(grp)
	if err != nil {
		log.Errorf("Could not create group %s: %v", form.Name, err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_create_group"))
	}

	return ctx.Created(rw, grp)
}

type GroupItemResource struct {
	*base.Resource
}

func (r *GroupItemResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	grp, err := r.DB(c).FindOne(&Group{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_query_group"))
	}
	return ctx.OK(rw, grp)
}

func (r *GroupItemResource) PUT(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {

	vars := mux.Vars(req)
	id := vars["id"]

	// decode request data
	var form = &struct {
		Name string `json:"name"`
	}{}
	err := json.NewDecoder(req.Body).Decode(form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_parse_request_data"))
	}

	// get group from database
	var grp pg.Entity
	grp, err = r.DB(c).FindOne(&Group{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_query_group"))
	}

	// update the group
	grp.(*Group).Name = form.Name
	err = r.DB(c).Update(grp)
	if err != nil {
		log.Errorf("Could not edit group %s: %v", form.Name, err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_edit_group"))
	}

	return ctx.OK(rw, grp)
}

func (r *GroupItemResource) DELETE(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	grp, err := r.DB(c).FindOne(&Group{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_query_group"))
	}
	err = r.DB(c).Delete(grp)
	if err != nil {
		log.Errorf("Could not delete group %s: %v", id, err)
		return ctx.InternalServerError(rw, c.T("group.api.could_not_delete_group"))
	}
	return ctx.NoContent(rw)
}
