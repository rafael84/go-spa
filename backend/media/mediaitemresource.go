package media

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/gotk/ctx"

	"github.com/rafael84/go-spa/backend/base"
)

func init() {
	ctx.Resource("/media/{id:[0-9]+}", &MediaItemResource{}, false)
}

type MediaItemResource struct {
	*base.Resource
}

func (r *MediaItemResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := r.DB(c)
	vars := mux.Vars(req)
	id := vars["id"]

	media, err := db.FindOne(&Media{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query media id %s: %v", id, err)
		return ctx.BadRequest(rw, "Could not query media")
	}
	return ctx.OK(rw, media)
}

func (r *MediaItemResource) PUT(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := r.DB(c)
	vars := mux.Vars(req)
	id := vars["id"]

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

	// get media from database
	entity, err := db.FindOne(&Media{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query media id %s: %v", id, err)
		return ctx.BadRequest(rw, "Could not query media")
	}
	media := entity.(*Media)

	// update the media
	media.Name = form.Name
	media.LocationId = form.LocationId
	media.MediaTypeId = form.MediaTypeId
	media.Path = dstPath
	err = db.Update(media)
	if err != nil {
		log.Errorf("Could not edit media %s: %v", form.Name, err)
		return ctx.BadRequest(rw, "Could not edit media")
	}

	return ctx.OK(rw, media)
}

func (r *MediaItemResource) DELETE(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := r.DB(c)
	vars := mux.Vars(req)
	id := vars["id"]

	media, err := db.FindOne(&Media{}, "id = $1", id)
	if err != nil {
		log.Errorf("Could not query media id %s: %v", id, err)
		return ctx.BadRequest(rw, "Could not query media")
	}
	err = db.Delete(media)
	if err != nil {
		log.Errorf("Could not delete media %s: %v", id, err)
		return ctx.InternalServerError(rw, "Could not delete media")
	}
	return ctx.NoContent(rw)
}
