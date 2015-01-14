package location

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
	ctx.Resource("/location", &LocationResource{}, false)
	ctx.Resource("/location/{id:[0-9]+}", &LocationItemResource{}, false)
}

type LocationResource struct {
	*base.Resource
}

func (r *LocationResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	locations, err := r.DB(c).FindAll(&Location{}, "")
	if err != nil {
		log.Errorf("Could not query locations: %v", err)
		return ctx.BadRequest(rw, "Could not query locations")
	}
	return ctx.OK(rw, locations)
}

func (r *LocationResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {

	// decode request data
	var form = &struct {
		Name       string `json:"name"`
		StaticURL  string `json:"staticURL"`
		StaticPath string `json:"staticPath"`
	}{}
	err := json.NewDecoder(req.Body).Decode(form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return ctx.BadRequest(rw, "Could not parse request data")
	}

	// create new location
	location := &Location{
		Name:       form.Name,
		StaticURL:  form.StaticURL,
		StaticPath: form.StaticPath,
	}
	err = r.DB(c).Create(location)
	if err != nil {
		log.Errorf("Could not create location %s: %v", form.Name, err)
		return ctx.BadRequest(rw, "Could not create location")
	}

	return ctx.Created(rw, location)
}

type LocationItemResource struct {
	*base.Resource
}

func (r *LocationItemResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	location, err := r.DB(c).FindOne(&Location{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query location id %s: %v", id, err)
		return ctx.BadRequest(rw, "Could not query location")
	}
	return ctx.OK(rw, location)
}

func (r *LocationItemResource) PUT(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	// decode request data
	var form = &struct {
		Name       string `json:"name"`
		StaticURL  string `json:"staticURL"`
		StaticPath string `json:"staticPath"`
	}{}
	err := json.NewDecoder(req.Body).Decode(form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return ctx.BadRequest(rw, "Could not parse request data")
	}

	// get location from database
	var entity pg.Entity
	entity, err = r.DB(c).FindOne(&Location{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query location id %s: %v", id, err)
		return ctx.BadRequest(rw, "Could not query location")
	}
	location := entity.(*Location)

	// update the location
	location.Name = form.Name
	location.StaticURL = form.StaticURL
	location.StaticPath = form.StaticPath
	err = r.DB(c).Update(location)
	if err != nil {
		log.Errorf("Could not edit location %s: %v", form.Name, err)
		return ctx.BadRequest(rw, "Could not edit location")
	}

	return ctx.OK(rw, location)
}

func (r *LocationItemResource) DELETE(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	vars := mux.Vars(req)
	id := vars["id"]

	location, err := r.DB(c).FindOne(&Location{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query location id %s: %v", id, err)
		return ctx.BadRequest(rw, "Could not query location")
	}
	err = r.DB(c).Delete(location)
	if err != nil {
		log.Errorf("Could not delete location %s: %v", id, err)
		return ctx.InternalServerError(rw, "Could not delete location")
	}
	return ctx.NoContent(rw)
}
