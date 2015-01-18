package mediatype

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"

	"github.com/gorilla/mux"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"
)

func init() {
	ctx.Resource("/mediatype/{id:[0-9]+}", &Resource{}, false)
}

type Resource struct{}

func (r *Resource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	db := c.Vars["db"].(*pg.Session)

	mediatype, err := db.FindOne(&Model{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query media type id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_query_media_type"))
	}
	return ctx.OK(rw, mediatype)
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
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_parse_request"))
	}

	// get media type from database
	var entity pg.Entity
	entity, err = db.FindOne(&Model{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query media type id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_query_media_type"))
	}
	mediatype := entity.(*Model)

	// update the media type
	mediatype.Name = form.Name
	err = db.Update(mediatype)
	if err != nil {
		log.Errorf("Could not edit media type %s: %v", form.Name, err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_edit_media_type"))
	}

	return ctx.OK(rw, mediatype)
}

func (r *Resource) DELETE(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	db := c.Vars["db"].(*pg.Session)

	mediatype, err := db.FindOne(&Model{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query media type id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_query_media_type"))
	}
	err = db.Delete(mediatype)
	if err != nil {
		log.Errorf("Could not delete media type %s: %v", id, err)
		return ctx.InternalServerError(rw, c.T("mediatype.api.could_not_delete_media_type"))
	}
	return ctx.NoContent(rw)
}
