package account

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gotk/webctx"

	"github.com/rafael84/go-spa/backend/base"
	"github.com/rafael84/go-spa/backend/database"
)

const (
	emailRegex   = `\b[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}\b`
	tokenExpTime = 10 // minutes
)

func init() {
	webctx.Resource("/account/user/signup", &SignUpResource{}, true)
	webctx.Resource("/account/user/signin", &SignInResource{}, true)
	webctx.Resource("/account/user/me", &MeResource{}, false)
	webctx.Resource("/account/token/renew", &TokenRenewResource{}, false)
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
	tokenString, err := webctx.SignToken(token)
	if err != nil {
		return webctx.InternalServerError(rw, "Problem signing token")
	}
	return webctx.OK(rw, map[string]string{"token": tokenString})
}

type SignUpResource struct {
	*base.Resource
}

func (r *SignUpResource) POST(c *webctx.Context, rw http.ResponseWriter, req *http.Request) error {
	// decode request data
	var form SignUpForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		log.Errorf("Could not parse request data: %s", err)
		return webctx.BadRequest(rw, c.T("account.user.could_not_parse_request_data"))
	}

	// create new user service
	service := NewUserService(r.DB(c))
	// check whether the email address is already taken
	_, err = service.GetByEmail(form.Email)
	if err == nil {
		return webctx.BadRequest(rw, c.T("account.user.email_taken"))
	} else if err != database.ERecordNotFound {
		log.Errorf("Could not query user: %s", err)
		return webctx.InternalServerError(rw, "account.user.could_not_query_user")
	}

	// password validation
	if form.Password != form.PasswordAgain {
		return webctx.BadRequest(rw, c.T("account.user.passwords_mismatch"))
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
		return webctx.InternalServerError(rw, "Could not create user: %s", err)
	}

	// return created user data
	return webctx.Created(rw, user)
}

type SignInResource struct {
	*base.Resource
}

func (r *SignInResource) POST(c *webctx.Context, rw http.ResponseWriter, req *http.Request) error {
	// decode request data
	var form SignInForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return webctx.BadRequest(rw, "Could not query user: %s", err)
	}

	// validate email address
	if ok := regexp.MustCompile(emailRegex).MatchString(form.Email); !ok {
		return webctx.BadRequest(rw, "Invalid email address")
	}

	// validate password length
	if len(form.Password) == 0 {
		return webctx.BadRequest(rw, "Password cannot be empty")
	}

	// create new user service
	service := NewUserService(r.DB(c))

	// check user in database
	var user *User
	user, err = service.GetByEmail(form.Email)
	if err != nil {
		return webctx.BadRequest(rw, "Invalid email and/or password")
	}

	// check user password
	if !user.Password.Valid(form.Password) {
		return webctx.BadRequest(rw, "Invalid email and/or password")
	}

	// generate new token
	return tokenResponse(rw, newToken(user))
}

type TokenRenewResource struct {
	*base.Resource
}

func (r *TokenRenewResource) POST(c *webctx.Context, rw http.ResponseWriter, req *http.Request) error {

	// get user id from the current token
	userId, found := c.Token.Claims["uid"]
	if !found {
		return webctx.BadRequest(rw, "Could not extract user from context")
	}

	// create new user service
	service := NewUserService(r.DB(c))

	// check if user is still valid
	user, err := service.GetById(int64(userId.(float64)))
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return webctx.InternalServerError(rw, "Could not query user.")
	}

	// generate new token
	return tokenResponse(rw, newToken(user))

}

type MeResource struct {
	*base.Resource
}

func (r *MeResource) GET(c *webctx.Context, rw http.ResponseWriter, req *http.Request) error {
	// get user id from current token
	userId, found := c.Token.Claims["uid"]
	if !found {
		return webctx.BadRequest(rw, "Could not extract user from context")
	}

	// create new user service
	service := NewUserService(r.DB(c))

	// query user data
	user, err := service.GetById(int64(userId.(float64)))
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return webctx.InternalServerError(rw, "Could not query user.")
	}

	// return user data
	return webctx.OK(rw, user)
}

func (r *MeResource) PUT(c *webctx.Context, rw http.ResponseWriter, req *http.Request) error {

	// decode request data
	var form MeForm
	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return webctx.BadRequest(rw, "Could decode user profile data: %s", err)
	}

	// create new user service
	service := NewUserService(r.DB(c))

	// query user data
	user, err := service.GetById(form.Id.Int64)
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return webctx.InternalServerError(rw, "Could not query user.")
	}

	// get the json data from user
	jsonData, err := user.DecodeJsonData()
	if err != nil {
		return webctx.BadRequest(rw, "Could not decode json data")
	}

	// update the user
	user.Email = form.Email
	jsonData.FirstName = form.JsonData.FirstName
	jsonData.LastName = form.JsonData.LastName
	user.JsonData.Set(jsonData)
	service.Update(user)

	// return user data
	return webctx.OK(rw, user)
}
