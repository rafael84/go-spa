package token

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gotk/ctx"

	"github.com/rafael84/go-spa/backend/account/user"
)

// New generate a new JWT token.
// The expiration date is defined by `tokenExpTime`
func New(user *user.Model) *jwt.Token {
	token := jwt.New(jwt.GetSigningMethod("RS256"))
	token.Claims["uid"] = user.Id.Int64
	token.Claims["user"] = user
	token.Claims["exp"] = time.Now().Add(time.Minute * tokenExpTime).Unix()
	return token
}

func Response(c *ctx.Context, rw http.ResponseWriter, token *jwt.Token) error {
	tokenString, err := ctx.SignToken(token)
	if err != nil {
		return ctx.InternalServerError(rw, c.T("user.token.problem_signing_token"))
	}
	return ctx.OK(rw, map[string]string{"token": tokenString})
}
