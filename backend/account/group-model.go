package account

import (
	"time"

	"github.com/guregu/null"
	"github.com/rafael84/go-spa/backend/database"
)

type Group struct {
	Id        null.Int       `db:"id,pk"       json:"id"`
	State     int            `db:"state"       json:"state"`
	CreatedAt time.Time      `db:"created_at"  json:"createdAt"`
	UpdatedAt time.Time      `db:"updated_at"  json:"updatedAt"`
	Name      string         `db:"name"        json:"name"`
	JsonData  database.JSONB `db:"json_data"   json:"jsonData,omitempty"`
}

func (_ *Group) Table() string { return "account.group" }
