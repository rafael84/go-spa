package location

import (
	"time"

	"github.com/gotk/pg"
	"github.com/guregu/null"
)

type Model struct {
	Id         null.Int  `db:"id,autofilled,pk"          json:"id"`
	State      int       `db:"state"                     json:"state"`
	CreatedAt  time.Time `db:"created_at,autofilled"     json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at,autofilled"     json:"updatedAt"`
	Name       string    `db:"name"                      json:"name"`
	StaticURL  string    `db:"static_url"                json:"staticURL"`
	StaticPath string    `db:"static_path"               json:"staticPath"`
	JsonData   pg.JSONB  `db:"json_data"                 json:"jsonData,omitempty"`
}

func (_ *Model) Table() string { return "storage.location" }
