package mediatype

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"
)

func init() {
	ctx.Resource("/mediatype", &Collection{}, false)
}

type Collection struct{}

func (r *Collection) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)
	mediatypes, err := db.FindAll(&Model{}, "")
	if err != nil {
		log.Errorf("Query error: %v", err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_query_media_type"))
	}
	return ctx.OK(rw, mediatypes)
}

func (r *Collection) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
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

	// create new mediatype
	mediatype := &Model{
		Name: form.Name,
	}
	err = db.Create(mediatype)
	if err != nil {
		log.Errorf("Could not create media type %s: %v", form.Name, err)
		return ctx.BadRequest(rw, c.T("mediatype.api.could_not_create_media_type"))
	}

	return ctx.Created(rw, mediatype)
}
