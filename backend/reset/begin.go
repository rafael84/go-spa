package reset

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/base"
	"github.com/rafael84/go-spa/backend/cfg"
	"github.com/rafael84/go-spa/backend/mail"
	"github.com/rafael84/go-spa/backend/user"
)

func init() {
	ctx.Resource("/account/reset-password", &ResetResource{}, true)
}

type ResetForm struct {
	Email string `json:"email"`
}

type ResetResource struct {
	*base.Resource
}

func (r *ResetResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	// decode request data
	var form ResetForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return ctx.BadRequest(rw, c.T("reset.begin.could_not_query"))
	}

	// validate email address
	if ok := regexp.MustCompile(base.EmailRegex).MatchString(form.Email); !ok {
		return ctx.BadRequest(rw, c.T("reset.begin.invalid_email_address"))
	}

	// create new user service
	userService := user.NewUserService(r.DB(c))

	// get user from database
	var u *user.User
	u, err = userService.GetByEmail(form.Email)
	if err != nil {
		return ctx.BadRequest(rw, c.T("reset.begin.user_not_found"))
	}

	go sendResetPasswordEmail(c, u)

	return ctx.OK(rw, c.T("reset.begin.email_sent"))
}

func sendResetPasswordEmail(c *ctx.Context, u *user.User) {
	var body bytes.Buffer

	resetTokenService := NewResetTokenService(c.Vars["db"].(*pg.Session))

	resetToken, err := resetTokenService.Create(u.Id.NullInt64.Int64)
	if err != nil {
		log.Errorf("Unable to create a new reset token: %s", err)
		return
	}

	body.WriteString("Access this link: ")
	body.WriteString(fmt.Sprintf("%s/#/reset-password/step2/", cfg.Server.BasePath()))
	body.WriteString(resetToken.Key)

	err = mail.NewGmailAccount(
		cfg.Email.Username,
		cfg.Email.Password,
	).Send(&mail.Message{
		From:    cfg.Email.From,
		To:      []string{u.Email},
		Subject: cfg.Email.Subject,
		Body:    body.Bytes(),
	})

	if err != nil {
		log.Errorf("Unable to send email: %s", err)
		return
	}

}
