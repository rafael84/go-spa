package mediatype

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
	ctx.Resource("/mediatype", &MediaTypeResource{}, false)
	ctx.Resource("/mediatype/{id:[0-9]+}", &MediaTypeItemResource{}, false)
}

type MediaTypeResource struct {
	*base.Resource
}

func (r *MediaTypeResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	mediaTypes, err := r.DB(c).FindAll(&MediaType{}, "")
	if err != nil {
		log.Errorf("Query error: %v", err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_query_media_type"))
	}
	return ctx.OK(rw, mediaTypes)
}

func (r *MediaTypeResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {

	// decode request data
	var form = &struct {
		Name string `json:"name"`
	}{}
	err := json.NewDecoder(req.Body).Decode(form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_parse_request"))
	}

	// create new mediatype
	mediaType := &MediaType{
		Name: form.Name,
	}
	err = r.DB(c).Create(mediaType)
	if err != nil {
		log.Errorf("Could not create media type %s: %v", form.Name, err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_create_media_type"))
	}

	return ctx.Created(rw, mediaType)
}

type MediaTypeItemResource struct {
	*base.Resource
}

func (r *MediaTypeItemResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	mediaType, err := r.DB(c).FindOne(&MediaType{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query media type id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_query_media_type"))
	}
	return ctx.OK(rw, mediaType)
}

func (r *MediaTypeItemResource) PUT(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

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
	entity, err = r.DB(c).FindOne(&MediaType{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query media type id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_query_media_type"))
	}
	mediaType := entity.(*MediaType)

	// update the media type
	mediaType.Name = form.Name
	err = r.DB(c).Update(mediaType)
	if err != nil {
		log.Errorf("Could not edit media type %s: %v", form.Name, err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_edit_media_type"))
	}

	return ctx.OK(rw, mediaType)
}

func (r *MediaTypeItemResource) DELETE(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	mediaType, err := r.DB(c).FindOne(&MediaType{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query media type id %s: %v", id, err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_query_media_type"))
	}
	err = r.DB(c).Delete(mediaType)
	if err != nil {
		log.Errorf("Could not delete media type %s: %v", id, err)
		return ctx.InternalServerError(rw, c.T("mediatype.api.could_not_delete_media_type"))
	}
	return ctx.NoContent(rw)
}
