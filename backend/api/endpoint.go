package api

import (
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/rafael84/go-spa/backend/context"
	"github.com/rafael84/go-spa/backend/database"
)

const (
	privKey = "api/keys/app.rsa"     // openssl genrsa -out app.rsa 2048
	pubKey  = "api/keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

var (
	endpoints = []*context.Endpoint{}
)

func AddEndpoint(endpoint *context.Endpoint) {
	endpoints = append(endpoints, endpoint)
}

func Configure(router *mux.Router, pathPrefix string, db *database.Session) error {
	apiRouter := router.PathPrefix(pathPrefix).Subrouter()

	context.LoadSecureKeys(privKey, pubKey)
	ctx, err := context.New(apiRouter, db)
	if err != nil {
		return err
	}
	ctx.AddEndpoints(endpoints...)

	return nil
}
