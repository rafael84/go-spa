package user

import (
	"time"

	"github.com/gotk/pg"
	"github.com/guregu/null"

	"github.com/rafael84/go-spa/backend/password"
)

const (
	UserStateActive = iota
	UserStateInactive
)

type Model struct {
	Id        null.Int        `db:"id,autofilled,pk"        json:"id"`
	State     int             `db:"state"                   json:"state"`
	CreatedAt time.Time       `db:"created_at,autofilled"   json:"createdAt"`
	UpdatedAt time.Time       `db:"updated_at,autofilled"   json:"updatedAt"`
	Email     string          `db:"email"                   json:"email"`
	Password  password.Salted `db:"password"                json:"-"`
	Role      int             `db:"role"                    json:"role"`
	JsonData  pg.JSONB        `db:"json_data"               json:"jsonData,omitempty"`
}

type UserJsonData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (_ *Model) Table() string {
	return "account.user"
}

func (u *Model) DecodeJsonData() (*UserJsonData, error) {
	userJsonData := &UserJsonData{}
	err := u.JsonData.Decode(userJsonData)
	if err != nil {
		return nil, err
	}
	return userJsonData, nil
}
