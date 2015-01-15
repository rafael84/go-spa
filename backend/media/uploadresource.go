package media

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"
	"github.com/rafael84/go-spa/backend/base"
	"github.com/rafael84/go-spa/backend/location"
	"github.com/rafael84/go-spa/backend/mediatype"
)

func init() {
	ctx.Resource("/media/upload", &MediaUploadResource{}, false)
}

type MediaUploadResource struct {
	*base.Resource
}

func (r *MediaUploadResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	reader, err := req.MultipartReader()
	if err != nil {
		return ctx.BadRequest(rw, c.T("media.uploadresource.could_not_upload_file"))
	}
	var tempFile *os.File
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		tempFile, err = ioutil.TempFile(os.TempDir(), "spa")
		if err != nil {
			return ctx.InternalServerError(rw, c.T("media.uploadresource.could_not_create_temp_file"))
		}
		defer tempFile.Close()

		_, err = io.Copy(tempFile, part)
		if err != nil {
			break
		}
	}
	return ctx.Created(rw, tempFile.Name())
}

func getLocation(db *pg.Session, id int) (*location.Location, error) {
	entity, err := db.FindOne(&location.Location{}, "id=$1", id)
	if err != nil {
		return nil, err
	}
	return entity.(*location.Location), nil
}

func getMediaType(db *pg.Session, id int) (*mediatype.MediaType, error) {
	entity, err := db.FindOne(&mediatype.MediaType{}, "id=$1", id)
	if err != nil {
		return nil, err
	}
	return entity.(*mediatype.MediaType), nil
}

func moveUploadedFile(location *location.Location, mediaType *mediatype.MediaType, srcPath string) (string, error) {

	// create directories if necessary
	dir := fmt.Sprintf("/var/%s/%s", location.StaticPath, mediaType.Name)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		log.Errorf("Unable to create directory: %s", err)
		return "", errors.New("Could not process uploaded file")
	}

	// generate filename randomly
	filename, err := base.Random(16)
	if err != nil {
		log.Errorf("Unable to generate filename: %s", err)
		return "", errors.New("Could not process uploaded file")
	}
	dstPath := fmt.Sprintf("%s/%s", dir, filename)

	// move file to its destination
	err = os.Rename(srcPath, dstPath)
	if err != nil {
		log.Errorf("Could not move file %s", err)
		return "", errors.New("Could not process uploaded file")
	}

	return dstPath, nil
}
