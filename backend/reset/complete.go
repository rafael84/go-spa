package reset

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"

	"github.com/rafael84/go-spa/backend/base"
	"github.com/rafael84/go-spa/backend/user"
)

func init() {
	ctx.Resource("/account/reset-password/complete", &CompleteResource{}, true)
}

type CompleteResource struct {
	*base.Resource
}

func (r *CompleteResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	type ChangePasswordForm struct {
		Password      string   `json:"password"`
		PasswordAgain string   `json:"passwordAgain"`
		ValidKey      ValidKey `json:"validKey"`
	}

	// decode request data
	var form ChangePasswordForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return ctx.BadRequest(rw, "Unable to change the password")
	}

	// validate the passwords
	if form.Password != form.PasswordAgain {
		return ctx.BadRequest(rw, "Passwords mismatch")
	}

	// validate the key again
	resetTokenService := NewResetTokenService(r.DB(c))
	resetToken, err := resetTokenService.GetByKey(form.ValidKey.Key)
	if err != nil || !resetToken.Valid() {
		return ctx.BadRequest(rw, "Invalid Key")
	}

	// get user from db
	userService := user.NewUserService(r.DB(c))
	u, err := userService.GetById(resetToken.UserId)
	if err != nil {
		return ctx.InternalServerError(rw, "User not found")
	}

	// encode user password
	err = u.Password.Encode(form.Password)
	if err != nil {
		return ctx.InternalServerError(rw, "Could not change user password")
	}

	// change user data in database
	err = userService.Update(u)
	if err != nil {
		return ctx.InternalServerError(rw, "Could not change user password")
	}

	// invalidate token
	resetToken.State = ResetTokenInactive
	err = resetTokenService.Update(resetToken)
	if err != nil {
		log.Errorf("Unable to invalidate token: %s", err)
	}

	return ctx.OK(rw, u)
}
