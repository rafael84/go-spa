package mediatype

import (
	"time"

	"github.com/gotk/pg"
	"github.com/guregu/null"
)

type MediaType struct {
	Id        null.Int  `db:"id"             json:"id"`
	State     int       `db:"state"          json:"state"`
	CreatedAt time.Time `db:"created_at"     json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at"     json:"updatedAt"`
	Name      string    `db:"name"           json:"name"`
	JsonData  pg.JSONB  `db:"json_data"      json:"jsonData,omitempty"`
}

func (_ *MediaType) Table() string { return "storage.media_type" }
