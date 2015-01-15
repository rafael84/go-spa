package user

import (
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/dgrijalva/jwt-go"
	"github.com/gotk/ctx"

	"github.com/rafael84/go-spa/backend/base"
)

const (
	tokenExpTime = 10 // minutes
)

func init() {
	ctx.Resource("/account/token/renew", &TokenRenewResource{}, false)
}

type TokenRenewResource struct {
	*base.Resource
}

func (r *TokenRenewResource) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {

	// get user id from the current token
	userId, found := c.Token.Claims["uid"]
	if !found {
		return ctx.BadRequest(rw, c.T("user.token.could_not_extract"))
	}

	// create new user service
	service := NewUserService(r.DB(c))

	// check if user is still valid
	user, err := service.GetById(int64(userId.(float64)))
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return ctx.InternalServerError(rw, c.T("user.token.could_not_query"))
	}

	// generate new token
	return tokenResponse(c, rw, newToken(user))

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

func tokenResponse(c *ctx.Context, rw http.ResponseWriter, token *jwt.Token) error {
	tokenString, err := ctx.SignToken(token)
	if err != nil {
		return ctx.InternalServerError(rw, c.T("user.token.problem_signing_token"))
	}
	return ctx.OK(rw, map[string]string{"token": tokenString})
}
