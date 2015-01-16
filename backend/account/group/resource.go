package group

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"
)

func init() {
	ctx.Resource("/group/{id:[0-9]+}", &Resource{}, false)
}

type Resource struct{}

func (r *Resource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	db := c.Vars["db"].(*pg.Session)

	grp, err := db.FindOne(&Model{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_query_group"))
	}
	return ctx.OK(rw, grp)
}

func (r *Resource) PUT(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	db := c.Vars["db"].(*pg.Session)

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
	grp, err = db.FindOne(&Model{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_query_group"))
	}

	// update the group
	grp.(*Model).Name = form.Name
	err = db.Update(grp)
	if err != nil {
		log.Errorf("Could not edit group %s: %v", form.Name, err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_edit_group"))
	}

	return ctx.OK(rw, grp)
}

func (r *Resource) DELETE(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	db := c.Vars["db"].(*pg.Session)

	grp, err := db.FindOne(&Model{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query group id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_query_group"))
	}
	err = db.Delete(grp)
	if err != nil {
		log.Errorf("Could not delete group %s: %v", id, err)
		return ctx.InternalServerError(rw, c.T("group.api.could_not_delete_group"))
	}
	return ctx.NoContent(rw)
}
