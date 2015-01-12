package storage

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/rafael84/go-spa/backend/base"
	"github.com/rafael84/go-spa/backend/context"
)

func init() {
	context.Resource("/storage/location", &LocationResource{}, false)
	context.Resource("/storage/location/{id:[0-9]+}", &LocationItemResource{}, false)
}

type LocationResource struct {
	*base.Resource
}

func (r *LocationResource) GET(c *context.Context, rw http.ResponseWriter, req *http.Request) error {
	locations, err := r.DB(c).Filter(&Location{}, "")
	if err != nil {
		log.Errorf("Query error: %v", err)
		return context.InternalServerError(rw, "Query error")
	}
	return context.OK(rw, locations)
}

type LocationItemResource struct {
	*base.Resource
}

func (r *LocationItemResource) GET(c *context.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	locations, err := r.DB(c).Filter(&Location{}, "id = $1", id)
	if err != nil {
		log.Errorf("Query error: %v", err)
		return context.InternalServerError(rw, "Query error")
	}
	return context.OK(rw, locations)
}
