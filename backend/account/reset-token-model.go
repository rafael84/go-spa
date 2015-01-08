package account

import (
	"time"

	"github.com/guregu/null"
)

type ResetToken struct {
	Id         null.Int  `db:"id,autofilled"           json:"id"`
	State      int       `db:"state"                   json:"state"`
	CreatedAt  time.Time `db:"created_at,autofilled"   json:"createdAt"`
	UpdatedAt  time.Time `db:"updated_at,autofilled"   json:"updatedAt"`
	Key        string    `db:"key"                     json:"key"`
	Expiration time.Time `db:"expiration"              json:"expiration"`
	UserId     int64     `db:"user_id"                 json:"userId"`
}

func (_ *ResetToken) Table() string {
	return "account.reset_token"
}

func (r *ResetToken) Valid() bool {
	return r.Expiration.Before(time.Now())
}
