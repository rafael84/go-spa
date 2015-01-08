package account

import (
	"time"

	"github.com/guregu/null"

	"github.com/rafael84/go-spa/backend/database"
)

const (
	UserStateActive = iota
	UserStateInactive
)

type User struct {
	Id        null.Int       `db:"id,autofilled"           json:"id"`
	State     int            `db:"state"                   json:"state"`
	CreatedAt time.Time      `db:"created_at,autofilled"   json:"createdAt"`
	UpdatedAt time.Time      `db:"updated_at,autofilled"   json:"updatedAt"`
	Email     string         `db:"email"                   json:"email"`
	Password  string         `db:"password"                json:"-"`
	JsonData  database.JSONB `db:"json_data"               json:"jsonData,omitempty"`
}

type UserJsonData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (_ *User) Table() string {
	return "account.user"
}

func (u *User) DecodeJsonData() (*UserJsonData, error) {
	userJsonData := &UserJsonData{}
	err := u.JsonData.Get(userJsonData)
	if err != nil {
		return nil, err
	}
	return userJsonData, nil
}
