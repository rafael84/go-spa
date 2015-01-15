package user

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/base"
)

func init() {
	ctx.Resource("/account/user/signup", &SignUpResource{}, true)
}

type SignUpForm struct {
	FirstName     string `json:"firstName"`
	LastName      string `json:"lastName"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	PasswordAgain string `json:"passwordAgain"`
}

type SignUpResource struct {
	*base.Resource
}

func (r *SignUpResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	// decode request data
	var form SignUpForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return ctx.BadRequest(rw, c.T("user.signup.could_not_parse_request_data"))
	}

	// create new user service
	service := NewUserService(r.DB(c))
	// check whether the email address is already taken
	_, err = service.GetByEmail(form.Email)
	if err == nil {
		return ctx.BadRequest(rw, c.T("user.signup.email_taken"))
	} else if err != pg.ERecordNotFound {
		log.Errorf("Could not query user: %s", err)
		return ctx.InternalServerError(rw, c.T("user.signup.could_not_query_user"))
	}

	// password validation
	if form.Password != form.PasswordAgain {
		return ctx.BadRequest(rw, c.T("user.signup.passwords_mismatch"))
	}

	// create new user
	user, err := service.Create(
		form.Email,
		form.Password,
		&UserJsonData{
			FirstName: form.FirstName,
			LastName:  form.LastName,
		},
	)
	if err != nil {
		return ctx.InternalServerError(rw, c.T("user.signup.could_not_create_user"))
	}

	// return created user data
	return ctx.Created(rw, user)
}
