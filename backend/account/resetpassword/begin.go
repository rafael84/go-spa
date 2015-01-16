package resetpassword

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/account/user"
	"github.com/rafael84/go-spa/backend/cfg"
)

func init() {
	ctx.Resource("/account/reset-password", &Begin{}, true)
}

type Begin struct{}

func (r *Begin) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)

	// decode request data
	var form struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return ctx.BadRequest(rw, c.T("reset.begin.could_not_query"))
	}

	// validate email address
	if ok := regexp.MustCompile(cfg.Email.Regex).MatchString(form.Email); !ok {
		return ctx.BadRequest(rw, c.T("reset.begin.invalid_email_address"))
	}

	// create new user service
	userService := user.NewUserService(db)

	// get user from database
	var u *user.Model
	u, err = userService.GetByEmail(form.Email)
	if err != nil {
		return ctx.BadRequest(rw, c.T("reset.begin.user_not_found"))
	}

	go sendEmail(c, u)

	return ctx.OK(rw, c.T("reset.begin.email_sent"))
}
