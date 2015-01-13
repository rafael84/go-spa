package storage

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gotk/ctx"

	"github.com/rafael84/go-spa/backend/base"
)

func init() {
	ctx.Resource("/storage/location", &LocationResource{}, false)
	ctx.Resource("/storage/location/{id:[0-9]+}", &LocationItemResource{}, false)
}

type LocationResource struct {
	*base.Resource
}

func (r *LocationResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	locations, err := r.DB(c).FindAll(&Location{}, "")
	if err != nil {
		log.Errorf("Query error: %v", err)
		return ctx.InternalServerError(rw, "Query error")
	}
	return ctx.OK(rw, locations)
}

type LocationItemResource struct {
	*base.Resource
}

func (r *LocationItemResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	locations, err := r.DB(c).FindAll(&Location{}, "id = $1", id)
	if err != nil {
		log.Errorf("Query error: %v", err)
		return ctx.InternalServerError(rw, "Query error")
	}
	return ctx.OK(rw, locations)
}
