package location

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"
)

func init() {
	ctx.Resource("/location", &Collection{}, false)
}

type Collection struct{}

func (lc *Collection) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)
	locations, err := db.FindAll(&Model{}, "")
	if err != nil {
		log.Errorf("Could not query locations: %v", err)
		return ctx.BadRequest(rw, c.T("location.api.could_not_query_locations"))
	}
	return ctx.OK(rw, locations)
}

func (lc *Collection) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
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

	// create new location
	loc := &Model{
		Name:       form.Name,
		StaticURL:  form.StaticURL,
		StaticPath: form.StaticPath,
	}
	err = db.Create(loc)
	if err != nil {
		log.Errorf("Could not create location %s: %v", form.Name, err)
		return ctx.BadRequest(rw, c.T("location.api.could_not_create_location"))
	}

	return ctx.Created(rw, loc)
}
