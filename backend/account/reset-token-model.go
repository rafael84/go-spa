package account

import (
	"time"

	"github.com/guregu/null"
)

type ResetToken struct {
	Id        null.Int  `db:"id,autofilled"           json:"id"`
	State     int       `db:"state"                   json:"state"`
	CreatedAt time.Time `db:"created_at,autofilled"   json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at,autofilled"   json:"updatedAt"`
	Salt      string    `db:"salt"                    json:"-"`
}

func (_ *ResetToken) Table() string {
	return "account.reset_token"
}
