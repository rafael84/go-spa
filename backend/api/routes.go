package api

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	simpleRoutes = map[string]context.Handler{}
	secureRoutes = map[string]context.SecureHandler{}
)

func AddSimpleRoute(endpoint string, handler context.Handler) {
	simpleRoutes[endpoint] = handler
}

func AddSecureRoute(endpoint string, handler context.SecureHandler) {
	secureRoutes[endpoint] = handler
}

func Configure(router *mux.Router, pathPrefix string, db *database.Session) error {
	apiRouter := router.PathPrefix(pathPrefix).Subrouter()

	// add simple routes
	simpleContext := context.NewContext(apiRouter, db)
	simpleContext.AddRoutes(simpleRoutes)

	// add secure routes
	context.LoadSecureKeys(privKey, pubKey)
	secureContext, err := context.NewSecureContext(simpleContext)
	if err != nil {
		return err
	}
	secureContext.AddRoutes(secureRoutes)

	return nil
}

func Success(rw http.ResponseWriter, statusCode int, response interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}
	rw.Write(jsonResponse)
	return nil
}

func Errorf(rw http.ResponseWriter, statusCode int, message string, args ...interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(statusCode)
	response := map[string]string{
		"error": fmt.Sprintf(message, args...),
	}
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return err
	}
	rw.Write(jsonResponse)
	return nil
}

func OK(rw http.ResponseWriter, response interface{}) error {
	return Success(rw, http.StatusOK, response)
}

func Created(rw http.ResponseWriter, response interface{}) error {
	return Success(rw, http.StatusCreated, response)
}

func BadRequest(rw http.ResponseWriter, message string, args ...interface{}) error {
	return Errorf(rw, http.StatusBadRequest, message, args...)
}

func InternalServerError(rw http.ResponseWriter, message string, args ...interface{}) error {
	return Errorf(rw, http.StatusInternalServerError, message, args...)
}
