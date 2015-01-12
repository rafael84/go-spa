package account

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"regexp"

	log "github.com/Sirupsen/logrus"

	"github.com/rafael84/go-spa/backend/api"
	"github.com/rafael84/go-spa/backend/context"
	"github.com/rafael84/go-spa/backend/mail"
)

type ValidKey struct {
	UserId int64  `json:"userId"`
	Key    string `json:"key"`
}

func init() {
	api.AddEndpoint(
		&context.Endpoint{
			Public: true,
			Path:   "/account/reset-password",
			Handlers: context.MethodHandlers{
				"POST": ResetPasswordHandler,
			},
		},
	)
	api.AddEndpoint(
		&context.Endpoint{
			Public: true,
			Path:   "/account/reset-password/validate-key",
			Handlers: context.MethodHandlers{
				"POST": ValidateKeyHandler,
			},
		},
	)
	api.AddEndpoint(
		&context.Endpoint{
			Public: true,
			Path:   "/account/reset-password/complete",
			Handlers: context.MethodHandlers{
				"POST": CompleteHandler,
			},
		},
	)
}

func sendResetPasswordEmail(c *context.Context, user *User) {
	var body bytes.Buffer

	resetTokenService := NewResetTokenService(c.DB)

	resetToken, err := resetTokenService.Create(user.Id.NullInt64.Int64)
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
		To:      []string{user.Email},
		Subject: "Reset Password",
		Body:    body.Bytes(),
	})

	if err != nil {
		log.Errorf("Unable to send email: %s", err)
		return
	}

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

	go sendResetPasswordEmail(c, user)

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
	resetTokenService := NewResetTokenService(c.DB)
	resetToken, err := resetTokenService.GetByKey(form.ValidKey.Key)
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

	// invalidate token
	resetToken.State = ResetTokenInactive
	err = resetTokenService.Update(resetToken)
	if err != nil {
		log.Errorf("Unable to invalidate token: %s", err)
	}

	return api.OK(rw, user)
}
