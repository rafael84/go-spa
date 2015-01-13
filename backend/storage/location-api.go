package storage

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gotk/webctx"

	"github.com/rafael84/go-spa/backend/base"
)

func init() {
	webctx.Resource("/storage/location", &LocationResource{}, false)
	webctx.Resource("/storage/location/{id:[0-9]+}", &LocationItemResource{}, false)
}

type LocationResource struct {
	*base.Resource
}

func (r *LocationResource) GET(c *webctx.Context, rw http.ResponseWriter, req *http.Request) error {
	locations, err := r.DB(c).FindAll(&Location{}, "")
	if err != nil {
		log.Errorf("Query error: %v", err)
		return webctx.InternalServerError(rw, "Query error")
	}
	return webctx.OK(rw, locations)
}

type LocationItemResource struct {
	*base.Resource
}

func (r *LocationItemResource) GET(c *webctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	locations, err := r.DB(c).FindAll(&Location{}, "id = $1", id)
	if err != nil {
		log.Errorf("Query error: %v", err)
		return webctx.InternalServerError(rw, "Query error")
	}
	return webctx.OK(rw, locations)
}
