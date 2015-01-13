package user

import (
	"encoding/json"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gotk/ctx"
	"github.com/guregu/null"

	"github.com/rafael84/go-spa/backend/base"
)

func init() {
	ctx.Resource("/account/user/me", &MeResource{}, false)
}

type MeResource struct {
	*base.Resource
}

type MeForm struct {
	Id       null.Int     `json:"id"`
	Email    string       `json:"email"`
	JsonData UserJsonData `json:"jsonData,omitempty"`
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
