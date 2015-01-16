package location

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"
)

func init() {
	ctx.Resource("/location/{id:[0-9]+}", &Resource{}, false)
}

type Resource struct{}

func (lr *Resource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	db := c.Vars["db"].(*pg.Session)

	location, err := db.FindOne(&Model{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query location id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("location.api.could_not_query_location"))
	}
	return ctx.OK(rw, location)
}

func (lr *Resource) PUT(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	db := c.Vars["db"].(*pg.Session)

	// decode request data
	var form = &struct {
		Name       string `json:"name"`
		StaticURL  string `json:"staticURL"`
		StaticPath string `json:"staticPath"`
	}{}
	err := json.NewDecoder(req.Body).Decode(form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return ctx.BadRequest(rw, c.T("location.api.could_not_parse_request_data"))
	}

	// get location from database
	var entity pg.Entity
	entity, err = db.FindOne(&Model{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query location id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("location.api.could_not_query_location"))
	}
	location := entity.(*Model)

	// update the location
	location.Name = form.Name
	location.StaticURL = form.StaticURL
	location.StaticPath = form.StaticPath
	err = db.Update(location)
	if err != nil {
		log.Errorf("Could not edit location %s: %v", form.Name, err)
		return ctx.BadRequest(rw, c.T("location.api.could_not_edit_location"))
	}

	return ctx.OK(rw, location)
}

func (lr *Resource) DELETE(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	db := c.Vars["db"].(*pg.Session)

	location, err := db.FindOne(&Model{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query location id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("location.api.could_not_query_location"))
	}
	err = db.Delete(location)
	if err != nil {
		log.Errorf("Could not delete location %s: %v", id, err)
		return ctx.InternalServerError(rw, c.T("location.api.could_not_delete_location"))
	}
	return ctx.NoContent(rw)
}
