package token

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/account/user"
)

const (
	tokenExpTime = 10 // minutes
)

func init() {
	ctx.Resource("/account/token/renew", &Renew{}, false)
}

type Renew struct{}

func (r *Renew) POST(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)

	// get user id from the current token
	userId, found := c.Token.Claims["uid"]
	if !found {
		return ctx.BadRequest(rw, c.T("user.token.could_not_extract"))
	}

	// check if user is still valid
	user, err := user.GetById(db, int64(userId.(float64)))
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return ctx.InternalServerError(rw, c.T("user.token.could_not_query"))
	}

	// generate new token
	return Response(c, rw, New(user))

}
