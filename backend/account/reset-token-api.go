package account

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"regexp"

	"github.com/rafael84/go-spa/backend/api"
	"github.com/rafael84/go-spa/backend/context"
	"github.com/rafael84/go-spa/backend/mail"
)

type ValidKey struct {
	UserId int64  `json:"userId"`
	Key    string `json:"key"`
}

func init() {
	api.AddSimpleRoute("/account/reset-password", ResetPasswordHandler)
	api.AddSimpleRoute("/account/reset-password/validate-key", ValidateKeyHandler)
	api.AddSimpleRoute("/account/reset-password/complete", CompleteHandler)
}

func ResetPasswordHandler(c *context.Context, rw http.ResponseWriter, req *http.Request) error {
	// decode request data
	var form ResetPasswordForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return api.BadRequest(rw, "Could not query user: %s", err)
	}

	// validate email address
	if ok := regexp.MustCompile(emailRegex).MatchString(form.Email); !ok {
		return api.BadRequest(rw, "Invalid email address")
	}

	// create new user service
	userService := NewUserService(c.DB)

	// get user from database
	var user *User
	user, err = userService.GetByEmail(form.Email)
	if err != nil {
		return api.BadRequest(rw, "User not found")
	}

	go func(user *User) {
		var body bytes.Buffer

		resetTokenService := NewResetTokenService(c.DB)
		resetToken, err := resetTokenService.Create(user.Id.NullInt64.Int64)

		if err != nil {
			return
		}

		body.WriteString("Access this link: ")
		body.WriteString("http://localhost:3000/#/reset-password/step2/")
		body.WriteString(resetToken.Key)

		mail.NewGmailAccount(
			os.Getenv("EMAIL_USERNAME"),
			os.Getenv("EMAIL_PASSWORD"),
		).Send(&mail.Message{
			From:    "Go-SPA",
			To:      []string{form.Email},
			Subject: "Reset Password",
			Body:    body.Bytes(),
		})
	}(user)

	return api.OK(rw, "Email sent")
}

func ValidateKeyHandler(c *context.Context, rw http.ResponseWriter, req *http.Request) error {
	type ValidateKeyForm struct {
		Key string `json:"key"`
	}

	// decode request data
	var form ValidateKeyForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return api.BadRequest(rw, "Unable to validate key")
	}

	service := NewResetTokenService(c.DB)

	resetToken, err := service.GetByKey(form.Key)
	if err != nil || !resetToken.Valid() {
		return api.BadRequest(rw, "Invalid Key")
	}

	return api.OK(rw, ValidKey{resetToken.UserId, form.Key})
}

func CompleteHandler(c *context.Context, rw http.ResponseWriter, req *http.Request) error {
	type ChangePasswordForm struct {
		Password      string   `json:"password"`
		PasswordAgain string   `json:"passwordAgain"`
		ValidKey      ValidKey `json:"validKey"`
	}

	// decode request data
	var form ChangePasswordForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return api.BadRequest(rw, "Unable to change the password")
	}

	// validate the passwords
	if form.Password != form.PasswordAgain {
		return api.BadRequest(rw, "Passwords mismatch")
	}

	// validate the key again
	service := NewResetTokenService(c.DB)
	resetToken, err := service.GetByKey(form.ValidKey.Key)
	if err != nil || !resetToken.Valid() {
		return api.BadRequest(rw, "Invalid Key")
	}

	// get user from db
	userService := NewUserService(c.DB)
	user, err := userService.GetById(resetToken.UserId)
	if err != nil {
		return api.InternalServerError(rw, "User not found")
	}

	// encode user password
	err = user.Password.Encode(form.Password)
	if err != nil {
		return api.InternalServerError(rw, "Could not change user password")
	}

	// change user data in database
	err = userService.Update(user)
	if err != nil {
		return api.InternalServerError(rw, "Could not change user password")
	}

	return api.OK(rw, "Password changed")
}