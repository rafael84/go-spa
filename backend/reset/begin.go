package reset

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/base"
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
		return ctx.BadRequest(rw, "Could not query user: %s", err)
	}

	// validate email address
	if ok := regexp.MustCompile(base.EmailRegex).MatchString(form.Email); !ok {
		return ctx.BadRequest(rw, "Invalid email address")
	}

	// create new user service
	userService := user.NewUserService(r.DB(c))

	// get user from database
	var u *user.User
	u, err = userService.GetByEmail(form.Email)
	if err != nil {
		return ctx.BadRequest(rw, "User not found")
	}

	go sendResetPasswordEmail(c, u)

	return ctx.OK(rw, "Email sent")
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
	body.WriteString("http://localhost:3000/#/reset-password/step2/")
	body.WriteString(resetToken.Key)

	err = mail.NewGmailAccount(
		os.Getenv("EMAIL_USERNAME"),
		os.Getenv("EMAIL_PASSWORD"),
	).Send(&mail.Message{
		From:    "Go-SPA",
		To:      []string{u.Email},
		Subject: "Reset Password",
		Body:    body.Bytes(),
	})

	if err != nil {
		log.Errorf("Unable to send email: %s", err)
		return
	}

}
