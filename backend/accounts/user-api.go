package accounts

import (
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"

	"github.com/rafael84/go-spa/backend/api"
	"github.com/rafael84/go-spa/backend/context"
	"github.com/rafael84/go-spa/backend/database"
	"github.com/rafael84/go-spa/backend/mail"
)

const (
	emailRegex   = `\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}\b`
	tokenExpTime = 10 // minutes
)

func init() {
	api.AddSimpleRoute("/accounts/user/resetPassword", ResetPasswordHandler)
	api.AddSimpleRoute("/accounts/user/signup", SignUpHandler)
	api.AddSimpleRoute("/accounts/user/signin", SignInHandler)
	api.AddSecureRoute("/accounts/user/me", MeHandler)
	api.AddSecureRoute("/accounts/user", UserHandler)
	api.AddSecureRoute("/accounts/token/renew", TokenRenewHandler)
}

// newToken generate a new JWT token.
// The expiration date is defined by `tokenExpTime`
func newToken(user *User) *jwt.Token {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["uid"] = user.Id.Int64
	token.Claims["user"] = user
	token.Claims["exp"] = time.Now().Add(time.Minute * tokenExpTime).Unix()
	return token
}

func tokenResponse(rw http.ResponseWriter, token *jwt.Token) error {
	tokenString, err := context.SignToken(token)
	if err != nil {
		return api.InternalServerError(rw, "Problem signin token")
	}
	return api.OK(rw, map[string]string{"token": tokenString})
}

// SignUpHandler handles the user registration logic
func SignUpHandler(c *context.Context, rw http.ResponseWriter, req *http.Request) error {
	// decode request data
	var form SignUpForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return api.BadRequest(rw, c.T("accounts.user.could_not_parse_request_data"))
	}

	// create new user service
	service := NewUserService(c.DB)
	// check whether the email address is already taken
	_, err = service.GetByEmail(form.Email)
	if err == nil {
		return api.BadRequest(rw, c.T("accounts.user.email_taken"))
	} else if err != database.ERecordNotFound {
		log.Errorf("Could not query user: %s", err)
		return api.InternalServerError(rw, "accounts.user.could_not_query_user")
	}

	// password validation
	if form.Password != form.Password {
		return api.BadRequest(rw, c.T("accounts.user.passwords_mismatch"))
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
		return api.InternalServerError(rw, "Could not create user: %s", err)
	}

	// return created user data
	return api.Created(rw, user)
}

func SignInHandler(c *context.Context, rw http.ResponseWriter, req *http.Request) error {

	// decode request data
	var form SignInForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return api.BadRequest(rw, "Could not query user: %s", err)
	}

	// validate email address
	if ok := regexp.MustCompile(emailRegex).MatchString(form.Email); !ok {
		return api.BadRequest(rw, "Invalid email address")
	}

	// validate password length
	if len(form.Password) == 0 {
		return api.BadRequest(rw, "Password cannot be empty")
	}

	// create new user service
	service := NewUserService(c.DB)

	// check user in database
	var user *User
	user, err = service.GetByEmail(form.Email)
	if err != nil {
		return api.BadRequest(rw, "Invalid email and/or password")
	}

	// check user password
	if user.Password != form.Password {
		return api.BadRequest(rw, "Invalid email and/or password")
	}

	// generate new token
	return tokenResponse(rw, newToken(user))
}

func TokenRenewHandler(sc *context.SecureContext, rw http.ResponseWriter, req *http.Request) error {

	// get user id from the current token
	userId, found := sc.Token.Claims["uid"]
	if !found {
		return api.BadRequest(rw, "Could not extract user from secure context")
	}

	// create new user service
	service := NewUserService(sc.DB)

	// check if user is still valid
	user, err := service.GetById(int(userId.(float64)))
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return api.InternalServerError(rw, "Could not query user.")
	}

	// generate new token
	return tokenResponse(rw, newToken(user))
}

func MeHandler(sc *context.SecureContext, rw http.ResponseWriter, req *http.Request) error {

	// get user id from current token
	userId, found := sc.Token.Claims["uid"]
	if !found {
		return api.BadRequest(rw, "Could not extract user from secure context")
	}

	// create new user service
	service := NewUserService(sc.DB)

	// query user data
	user, err := service.GetById(int(userId.(float64)))
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return api.InternalServerError(rw, "Could not query user.")
	}

	// return user data
	return api.OK(rw, user)
}

func UserHandler(sc *context.SecureContext, rw http.ResponseWriter, req *http.Request) error {
	return api.OK(rw, "To be implemented")
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

	go func() {
		mail.NewGmailAccount(
			os.Getenv("EMAIL_USERNAME"),
			os.Getenv("EMAIL_PASSWORD"),
		).Send(&mail.Message{
			From:    "Go-SPA",
			To:      []string{form.Email},
			Subject: "Reset Password",
			Body:    []byte("Access this link."),
		})
	}()

	return api.OK(rw, "Email sent")
}
