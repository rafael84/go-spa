package media

import (
	"time"

	"github.com/gotk/pg"
	"github.com/guregu/null"
)

type Media struct {
	Id          null.Int  `db:"id,autofilled,pk"          json:"id"`
	State       int       `db:"state"                     json:"state"`
	CreatedAt   time.Time `db:"created_at,autofilled"     json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at,autofilled"     json:"updatedAt"`
	Name        string    `db:"name"                      json:"name"`
	MediaTypeId int       `db:"media_type_id"             json:"mediaTypeId"`
	LocationId  int       `db:"location_id"               json:"locationId"`
	Path        string    `db:"path"                      json:"path"`
	JsonData    pg.JSONB  `db:"json_data"                 json:"jsonData,omitempty"`
}

func (_ *Media) Table() string { return "storage.media" }
