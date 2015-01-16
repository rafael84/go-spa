package signin

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
	ctx.Resource("/account/user/signin", &Resource{}, true)
}

type Resource struct{}

func (r *Resource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)

	// decode request data
	var form struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return ctx.BadRequest(rw, "Could not query user") // TODO: translate this
	}

	// validate email address
	if ok := regexp.MustCompile(cfg.Email.Regex).MatchString(form.Email); !ok {
		return ctx.BadRequest(rw, c.T("user.signin.invalid_email_address"))
	}

	// validate password length
	if len(form.Password) == 0 {
		return ctx.BadRequest(rw, c.T("user.signin.password_cannot_be_empty"))
	}

	// create new user service
	service := user.NewUserService(db)

	// check user in database
	var user *user.Model
	user, err = service.GetByEmail(form.Email)
	if err != nil {
		return ctx.BadRequest(rw, c.T("user.signin.invalid_email_or_password"))
	}

	// check user password
	if !user.Password.Valid(form.Password) {
		return ctx.BadRequest(rw, c.T("user.signin.invalid_email_or_password"))
	}

	// generate new token
	return token.Response(c, rw, token.New(user))
}
