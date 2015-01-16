package resetpassword

import (
	"encoding/json"
	"net/http"

	"github.com/gotk/ctx"
	"github.com/gotk/pg"
)

type ValidKey struct {
	UserId int64  `json:"userId"`
	Key    string `json:"key"`
}

func init() {
	ctx.Resource("/account/reset-password/validate-key", &ValidateKey{}, true)
}

type ValidateKey struct{}

func (r *ValidateKey) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)

	// decode request data
	var form struct {
		Key string `json:"key"`
	}

	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return ctx.BadRequest(rw, c.T("reset.validate.unable_to_validate_key"))
	}

	resetToken, err := getToken(db, form.Key)
	if err != nil || !resetToken.Valid() {
		return ctx.BadRequest(rw, c.T("reset.validate.invalid_key"))
	}

	return ctx.OK(rw, ValidKey{resetToken.UserId, form.Key})
}
