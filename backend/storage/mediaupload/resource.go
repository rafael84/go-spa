package mediaupload

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gotk/ctx"
)

func init() {
	ctx.Resource("/media/upload", &Resource{}, false)
}

type Resource struct{}

func (r *Resource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
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
