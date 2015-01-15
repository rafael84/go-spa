package user

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gotk/ctx"

	"github.com/rafael84/go-spa/backend/base"
)

func init() {
	ctx.Resource("/account/user/signin", &SignInResource{}, true)
}

type SignInForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignInResource struct {
	*base.Resource
}

func (r *SignInResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	// decode request data
	var form SignInForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return ctx.BadRequest(rw, "Could not query user: %s", err) // how to translate this?
	}

	// validate email address
	if ok := regexp.MustCompile(base.EmailRegex).MatchString(form.Email); !ok {
		return ctx.BadRequest(rw, c.T("user.signin.invalid_email_address"))
	}

	// validate password length
	if len(form.Password) == 0 {
		return ctx.BadRequest(rw, c.T("user.signin.password_cannot_be_empty"))
	}

	// create new user service
	service := NewUserService(r.DB(c))

	// check user in database
	var user *User
	user, err = service.GetByEmail(form.Email)
	if err != nil {
		return ctx.BadRequest(rw, c.T("user.signin.invalid_email_or_password"))
	}

	// check user password
	if !user.Password.Valid(form.Password) {
		return ctx.BadRequest(rw, c.T("user.signin.invalid_email_or_password"))
	}

	// generate new token
	return tokenResponse(c, rw, newToken(user))
}
