package user

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"
	"github.com/guregu/null"
)

func init() {
	ctx.Resource("/account/user/me", &ProfileResource{}, false)
}

type ProfileResource struct{}

func (r *ProfileResource) GET(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)

	// get user id from current token
	userId, found := c.Token.Claims["uid"]
	if !found {
		return ctx.BadRequest(rw, c.T("user.me.could_not_extract"))
	}

	// create new user service
	service := NewUserService(db)

	// query user data
	user, err := service.GetById(int64(userId.(float64)))
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return ctx.InternalServerError(rw, c.T("user.me.could_not_query"))
	}

	// return user data
	return ctx.OK(rw, user)
}

func (r *ProfileResource) PUT(c *ctx.Context, rw http.ResponseWriter, req *http.Request) error {
	db := c.Vars["db"].(*pg.Session)

	// decode request data
	var form struct {
		Id       null.Int     `json:"id"`
		Email    string       `json:"email"`
		JsonData UserJsonData `json:"jsonData,omitempty"`
	}

	err := json.NewDecoder(req.Body).Decode(&form)
	if err != nil {
		return ctx.BadRequest(rw, c.T("user.me.could_not_decode_profile_data"))
	}

	// create new user service
	service := NewUserService(db)

	// query user data
	user, err := service.GetById(form.Id.Int64)
	if err != nil {
		log.Errorf("Could not query user: %v", err)
		return ctx.InternalServerError(rw, c.T("user.me.could_not_query"))
	}

	// get the json data from user
	jsonData, err := user.DecodeJsonData()
	if err != nil {
		return ctx.BadRequest(rw, c.T("user.me.could_not_decode_json_data"))
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
