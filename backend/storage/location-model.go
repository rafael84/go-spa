package storage

import (
	"time"

	"github.com/gotk/pg"
	"github.com/guregu/null"
)

type Location struct {
	Id         null.Int  `db:"id"             json:"id"`
	State      int       `db:"state"          json:"state"`
	CreatedAt  time.Time `db:"created_at"     json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at"     json:"updatedAt"`
	Name       string    `db:"name"           json:"name"`
	StaticUrl  string    `db:"static_url"     json:"staticUrl"`
	StaticPath string    `db:"static_path"    json:"staticPath"`
	JsonData   pg.JSONB  `db:"json_data"      json:"jsonData,omitempty"`
}

func (_ *Location) Table() string { return "storage.location" }
