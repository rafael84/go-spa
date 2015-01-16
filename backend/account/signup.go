package account

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/account/user"
)

func init() {
	ctx.Resource("/account/user/signup", &SignUp{}, true)
}

type SignUp struct{}

func (r *SignUp) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)

	// decode request data
	var form struct {
		FirstName     string `json:"firstName"`
		LastName      string `json:"lastName"`
		Email         string `json:"email"`
		Password      string `json:"password"`
		PasswordAgain string `json:"passwordAgain"`
		Role          int    `json:"role"`
	}
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return ctx.BadRequest(rw, c.T("user.signup.could_not_parse_request_data"))
	}

	// check whether the email address is already taken
	_, err = user.GetByEmail(db, form.Email)
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
	u, err := user.Create(
		db,
		form.Email,
		form.Password,
		form.Role,
		&user.UserJsonData{
			FirstName: form.FirstName,
			LastName:  form.LastName,
		},
	)
	if err != nil {
		return ctx.InternalServerError(rw, c.T("user.signup.could_not_create_user"))
	}

	// return created user data
	return ctx.Created(rw, u)
}
