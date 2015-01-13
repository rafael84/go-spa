package user

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gotk/ctx"

	"github.com/rafael84/go-spa/backend/base"
)

const (
	emailRegex = `\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}\b`
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
		return ctx.BadRequest(rw, "Could not query user: %s", err)
	}

	// validate email address
	if ok := regexp.MustCompile(emailRegex).MatchString(form.Email); !ok {
		return ctx.BadRequest(rw, "Invalid email address")
	}

	// validate password length
	if len(form.Password) == 0 {
		return ctx.BadRequest(rw, "Password cannot be empty")
	}

	// create new user service
	service := NewUserService(r.DB(c))

	// check user in database
	var user *User
	user, err = service.GetByEmail(form.Email)
	if err != nil {
		return ctx.BadRequest(rw, "Invalid email and/or password")
	}

	// check user password
	if !user.Password.Valid(form.Password) {
		return ctx.BadRequest(rw, "Invalid email and/or password")
	}

	// generate new token
	return tokenResponse(rw, newToken(user))
}
