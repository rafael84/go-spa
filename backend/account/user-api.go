package account

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/base"
)

const (
	emailRegex   = `\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}\b`
	tokenExpTime = 10 // minutes
)

func init() {
	ctx.Resource("/account/user/signup", &SignUpResource{}, true)
	ctx.Resource("/account/user/signin", &SignInResource{}, true)
	ctx.Resource("/account/user/me", &MeResource{}, false)
	ctx.Resource("/account/token/renew", &TokenRenewResource{}, false)
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
	tokenString, err := ctx.SignToken(token)
	if err != nil {
		return ctx.InternalServerError(rw, "Problem signing token")
	}
	return ctx.OK(rw, map[string]string{"token": tokenString})
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
		return ctx.BadRequest(rw, c.T("account.user.could_not_parse_request_data"))
	}

	// create new user service
	service := NewUserService(r.DB(c))
	// check whether the email address is already taken
	_, err = service.GetByEmail(form.Email)
	if err == nil {
		return ctx.BadRequest(rw, c.T("account.user.email_taken"))
	} else if err != pg.ERecordNotFound {
		log.Errorf("Could not query user: %s", err)
		return ctx.InternalServerError(rw, "account.user.could_not_query_user")
	}

	// password validation
	if form.Password != form.PasswordAgain {
		return ctx.BadRequest(rw, c.T("account.user.passwords_mismatch"))
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
		return ctx.InternalServerError(rw, "Could not create user: %s", err)
	}

	// return created user data
	return ctx.Created(rw, user)
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

type TokenRenewResource struct {
	*base.Resource
}

func (r *TokenRenewResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {

	// get user id from the current token
	userId, found := c.Token.Claims["uid"]
	if !found {
		return ctx.BadRequest(rw, "Could not extract user from context")
	}

	// create new user service
	service := NewUserService(r.DB(c))

	// check if user is still valid
	user, err := service.GetById(int64(userId.(float64)))
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return ctx.InternalServerError(rw, "Could not query user.")
	}

	// generate new token
	return tokenResponse(rw, newToken(user))

}

type MeResource struct {
	*base.Resource
}

func (r *MeResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	// get user id from current token
	userId, found := c.Token.Claims["uid"]
	if !found {
		return ctx.BadRequest(rw, "Could not extract user from context")
	}

	// create new user service
	service := NewUserService(r.DB(c))

	// query user data
	user, err := service.GetById(int64(userId.(float64)))
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return ctx.InternalServerError(rw, "Could not query user.")
	}

	// return user data
	return ctx.OK(rw, user)
}

func (r *MeResource) PUT(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {

	// decode request data
	var form MeForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return ctx.BadRequest(rw, "Could decode user profile data: %s", err)
	}

	// create new user service
	service := NewUserService(r.DB(c))

	// query user data
	user, err := service.GetById(form.Id.Int64)
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return ctx.InternalServerError(rw, "Could not query user.")
	}

	// get the json data from user
	jsonData, err := user.DecodeJsonData()
	if err != nil {
		return ctx.BadRequest(rw, "Could not decode json data")
	}

	// update the user
	user.Email = form.Email
	jsonData.FirstName = form.JsonData.FirstName
	jsonData.LastName = form.JsonData.LastName
	user.JsonData.Encode(jsonData)
	service.Update(user)

	// return user data
	return ctx.OK(rw, user)
}
