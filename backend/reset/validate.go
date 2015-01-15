package reset

import (
	"encoding/json"
	"net/http"

	"github.com/gotk/ctx"
	"github.com/rafael84/go-spa/backend/base"
)

type ValidKey struct {
	UserId int64  `json:"userId"`
	Key    string `json:"key"`
}

func init() {
	ctx.Resource("/account/reset-password/validate-key", &ValidateKeyResource{}, true)
}

type ValidateKeyForm struct {
	Key string `json:"key"`
}

type ValidateKeyResource struct {
	*base.Resource
}

func (r *ValidateKeyResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {

	// decode request data
	var form ValidateKeyForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return ctx.BadRequest(rw, c.T("reset.validate.unable_to_validate_key"))
	}

	service := NewResetTokenService(r.DB(c))

	resetToken, err := service.GetByKey(form.Key)
	if err != nil || !resetToken.Valid() {
		return ctx.BadRequest(rw, c.T("reset.validate.invalid_key"))
	}

	return ctx.OK(rw, ValidKey{resetToken.UserId, form.Key})
}
