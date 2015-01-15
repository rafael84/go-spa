package media

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"

	"github.com/rafael84/go-spa/backend/base"
)

func init() {
	ctx.Resource("/media", &MediaResource{}, false)
}

type MediaResource struct {
	*base.Resource
}

func (r *MediaResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	medias, err := r.DB(c).FindAll(&Media{}, "")
	if err != nil {
		log.Errorf("Query error: %v", err)
		return ctx.BadRequest(rw, "Query error")
	}
	return ctx.OK(rw, medias)
}

func (r *MediaResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := r.DB(c)
	// decode request data
	var form = &MediaForm{}
	err := json.NewDecoder(req.Body).Decode(form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return ctx.BadRequest(rw, "Could not parse request data")
	}

	// get location from database
	location, err := getLocation(db, form.LocationId)
	if err != nil {
		log.Errorf("Could not locate the requested location: %s", err)
		return ctx.BadRequest(rw, "Could not locate the requested location")
	}

	// get media type from database
	mediaType, err := getMediaType(db, form.MediaTypeId)
	if err != nil {
		log.Errorf("Could not locate the requested media type: %s", err)
		return ctx.BadRequest(rw, "Could not locate the requested media type")
	}

	// move the uploaded file to the right place
	var dstPath string
	dstPath, err = moveUploadedFile(location, mediaType, form.Path)
	if err != nil {
		log.Errorf("Could not process the uploaded file: %s", err)
		return ctx.InternalServerError(rw, "Could not process the uploaded file")
	}

	// create new media
	media := &Media{
		Name:        form.Name,
		MediaTypeId: form.MediaTypeId,
		LocationId:  form.LocationId,
		Path:        dstPath,
	}
	err = db.Create(media)
	if err != nil {
		log.Errorf("Could not create media %s: %v", form.Name, err)
		return ctx.BadRequest(rw, "Could not create media")
	}
	return ctx.Created(rw, media)
}
