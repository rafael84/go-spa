package storage

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"

	"github.com/rafael84/go-spa/backend/api"
	"github.com/rafael84/go-spa/backend/context"
	"github.com/rafael84/go-spa/backend/database"
)

func init() {
	api.AddPrivateEndpoint(
		&context.Endpoint{
			Path: "/storage/location",
			Handlers: context.MethodHandlers{
				"GET": LocationHandler,
			},
		},
	)
	api.AddPrivateEndpoint(
		&context.Endpoint{
			Path: "/storage/location/{id:[0-9]+}",
			Handlers: context.MethodHandlers{
				"GET": LocationHandler,
			},
		},
	)
}

func LocationHandler(sc *context.Context, rw http.ResponseWriter, req *http.Request) error {

	var locations []database.Entity
	var err error

	vars := mux.Vars(req)
	id, found := vars["id"]

	if found {
		locations, err = sc.DB.Filter(&Location{}, "id = $1", id)
	} else {
		locations, err = sc.DB.Filter(&Location{}, "")
	}
	if err != nil {
		log.Errorf("Query error: %v", err)
		return api.InternalServerError(rw, "Query error")
	}

	return api.OK(rw, locations)
}
