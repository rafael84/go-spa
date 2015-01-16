package group

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"
)

func init() {
	ctx.Resource("/group", &Collection{}, false)
}

type Collection struct{}

func (r *Collection) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)

	groups, err := db.FindAll(&Model{}, "")
	if err != nil {
		log.Errorf("Could not query groups: %v", err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_query_groups"))
	}

	return ctx.OK(rw, groups)
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
		return ctx.BadRequest(rw, c.T("group.api.could_not_parse_request_data"))
	}

	// create new group
	grp := &Model{Name: form.Name}
	err = db.Create(grp)
	if err != nil {
		log.Errorf("Could not create group %s: %v", form.Name, err)
		return ctx.BadRequest(rw, c.T("group.api.could_not_create_group"))
	}

	return ctx.Created(rw, grp)
}
