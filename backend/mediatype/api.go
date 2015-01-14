package mediatype

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gotk/ctx"

	"github.com/rafael84/go-spa/backend/base"
)

func init() {
	ctx.Resource("/storage/mediatype", &MediaTypeResource{}, false)
	ctx.Resource("/storage/mediatype/{id:[0-9]+}", &MediaTypeItemResource{}, false)
}

type MediaTypeResource struct {
	*base.Resource
}

func (r *MediaTypeResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	mediaTypes, err := r.DB(c).FindAll(&MediaType{}, "")
	if err != nil {
		log.Errorf("Query error: %v", err)
		return ctx.InternalServerError(rw, "Query error")
	}
	return ctx.OK(rw, mediaTypes)
}

type MediaTypeItemResource struct {
	*base.Resource
}

func (r *MediaTypeItemResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	mediaTypes, err := r.DB(c).FindAll(&MediaType{}, "id = $1", id)
	if err != nil {
		log.Errorf("Query error: %v", err)
		return ctx.InternalServerError(rw, "Query error")
	}
	return ctx.OK(rw, mediaTypes)
}
