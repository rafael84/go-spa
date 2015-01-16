package account

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/account/token"
	"github.com/rafael84/go-spa/backend/account/user"
	"github.com/rafael84/go-spa/backend/cfg"
)

func init() {
	ctx.Resource("/account/user/signin", &SignIn{}, true)
}

type SignIn struct{}

func (r *SignIn) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)

	// decode request data
	var form struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return ctx.BadRequest(rw, c.T("user.signin.could_not_query"))
	}

	// validate email address
	if ok := regexp.MustCompile(cfg.Email.Regex).MatchString(form.Email); !ok {
		return ctx.BadRequest(rw, c.T("user.signin.invalid_email_address"))
	}

	// validate password length
	if len(form.Password) == 0 {
		return ctx.BadRequest(rw, c.T("user.signin.password_cannot_be_empty"))
	}

	// check user in database
	var u *user.Model
	u, err = user.GetByEmail(db, form.Email)
	if err != nil {
		return ctx.BadRequest(rw, c.T("user.signin.invalid_email_or_password"))
	}

	// check user password
	if !u.Password.Valid(form.Password) {
		return ctx.BadRequest(rw, c.T("user.signin.invalid_email_or_password"))
	}

	// generate new token
	return token.Response(c, rw, token.New(u))
}
