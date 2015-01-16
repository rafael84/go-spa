package mediaupload

import (
	"errors"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"

	"github.com/rafael84/go-spa/backend/base"
	"github.com/rafael84/go-spa/backend/storage/location"
	"github.com/rafael84/go-spa/backend/storage/mediatype"
)

func MoveFile(location *location.Model, mediaType *mediatype.Model, srcPath string) (string, error) {

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
