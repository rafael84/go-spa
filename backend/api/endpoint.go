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
	privateEndpoints = []*context.Endpoint{}
)

func AddPrivateEndpoint(endpoint *context.Endpoint) {
	privateEndpoints = append(privateEndpoints, endpoint)
}

func AddResource(resource *resource) {
	privateEndpoints = append(privateEndpoints, &resource.Endpoint)
}

func Configure(router *mux.Router, pathPrefix string, db *database.Session) error {
	apiRouter := router.PathPrefix(pathPrefix).Subrouter()

	// add private endpoints
	context.LoadSecureKeys(privKey, pubKey)
	privateContext, err := context.New(apiRouter, db)
	if err != nil {
		return err
	}
	privateContext.AddEndpoints(privateEndpoints...)

	return nil
}
