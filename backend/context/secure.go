package context

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/auth0/go-jwt-middleware"
	"github.com/codegangsta/negroni"
	"github.com/dgrijalva/jwt-go"
)

type SecureContext struct {
	*Context

	Token      *jwt.Token
	middleware *jwtmiddleware.JWTMiddleware
}

type SecureHandler func(sc *SecureContext, rw http.ResponseWriter, req *http.Request) error

var (
	verifyKey []byte
	signKey   []byte
)

func SignToken(token *jwt.Token) (string, error) {
	return token.SignedString(signKey)
}

func LoadSecureKeys(privateKeyPath, publicKeyPath string) (err error) {
	signKey, err = ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return fmt.Errorf("Error reading private key")
	}
	verifyKey, err = ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return fmt.Errorf("Error reading public key")
	}
	return nil
}

func NewSecureContext(context *Context) (sc *SecureContext, err error) {
	sc = &SecureContext{Context: context}
	options := jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) { return verifyKey, nil },
	}
	sc.middleware = jwtmiddleware.New(options)
	return sc, nil
}

func (sc *SecureContext) AddRoute(path string, handler SecureHandler) {

	httpHandler := func(rw http.ResponseWriter, req *http.Request) {
		sc.Token, _ = jwt.ParseFromRequest(req, sc.middleware.Options.ValidationKeyGetter)
		sc.Context.updateT(req)
		err := handler(sc, rw, req)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusBadRequest)
		}
	}

	sc.Router.Handle(
		path, negroni.New(
			negroni.HandlerFunc(sc.middleware.HandlerWithNext),
			negroni.Wrap(http.HandlerFunc(httpHandler)),
		),
	)
}

func (sc *SecureContext) AddRoutes(routes map[string]SecureHandler) {
	for path, handler := range routes {
		sc.AddRoute(path, handler)
	}
}
