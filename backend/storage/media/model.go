package media

import (
	"fmt"
	"time"

	"github.com/gotk/pg"
	"github.com/guregu/null"

	"github.com/rafael84/go-spa/backend/storage/location"
	"github.com/rafael84/go-spa/backend/storage/mediatype"
)

type Model struct {
	Id          null.Int  `db:"id,autofilled,pk"          json:"id"`
	State       int       `db:"state"                     json:"state"`
	CreatedAt   time.Time `db:"created_at,autofilled"     json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at,autofilled"     json:"updatedAt"`
	Name        string    `db:"name"                      json:"name"`
	MediatypeId int       `db:"media_type_id"             json:"mediatypeId"`
	LocationId  int       `db:"location_id"               json:"locationId"`
	Path        string    `db:"path"                      json:"-"`
	Data        pg.JSONB  `db:"json_data"                 json:"data,omitempty"`
}

func (_ *Model) Table() string { return "storage.media" }

type MediaData struct {
	URL string `json:"url"`
}

// EncodeData serializes the media data struct into model's Data field
// (BUG): this only works when the model has a valid Id
func (m *Model) EncodeData(loc *location.Model, mt *mediatype.Model) {
	m.Data.Encode(&MediaData{
		URL: fmt.Sprintf("%s/%s/%d", loc.StaticURL, mt.Name, m.Id.Int64),
	})
}
