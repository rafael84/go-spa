package media

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/storage/location"
	"github.com/rafael84/go-spa/backend/storage/mediatype"
	"github.com/rafael84/go-spa/backend/storage/mediaupload"
)

func init() {
	ctx.Resource("/media", &Collection{}, false)
}

type Collection struct{}

func (r *Collection) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)
	medias, err := db.FindAll(&Model{}, "")
	if err != nil {
		log.Errorf("Query error: %v", err)
		return ctx.BadRequest(rw, c.T("media.mediaresource.query_error"))
	}
	return ctx.OK(rw, medias)
}

func (r *Collection) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)

	// decode request data
	var form = &MediaForm{}
	err := json.NewDecoder(req.Body).Decode(form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return ctx.BadRequest(rw, c.T("media.mediaresource.could_not_parse_request_data"))
	}

	// get location from database
	loc, err := location.GetById(db, form.LocationId)
	if err != nil {
		log.Errorf("Could not locate the requested location: %s", err)
		return ctx.BadRequest(rw, c.T("media.mediaresource.could_not_locate_the_requested_location"))
	}

	// get media type from database
	mt, err := mediatype.GetById(db, form.MediaTypeId)
	if err != nil {
		log.Errorf("Could not locate the requested media type: %s", err)
		return ctx.BadRequest(rw, c.T("media.mediaresource.could_not_locate_the_requested_media_type"))
	}

	// move the uploaded file to the right place
	var dstPath string
	dstPath, err = mediaupload.MoveFile(loc, mt, form.Path)
	if err != nil {
		log.Errorf("Could not process the uploaded file: %s", err)
		return ctx.InternalServerError(rw, c.T("media.mediaresource.could_not_process_uploaded_file"))
	}

	// create new media
	media := &Model{
		Name:        form.Name,
		MediaTypeId: form.MediaTypeId,
		LocationId:  form.LocationId,
		Path:        dstPath,
	}
	err = db.Create(media)
	if err != nil {
		log.Errorf("Could not create media %s: %v", form.Name, err)
		return ctx.BadRequest(rw, c.T("media.mediaresource.could_not_create_media"))
	}
	return ctx.Created(rw, media)
}
