package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/i18n"

	_ "github.com/rafael84/go-spa/backend/group"
	_ "github.com/rafael84/go-spa/backend/location"
	_ "github.com/rafael84/go-spa/backend/media"
	_ "github.com/rafael84/go-spa/backend/mediatype"
	"github.com/rafael84/go-spa/backend/middleware"
	_ "github.com/rafael84/go-spa/backend/reset"
	_ "github.com/rafael84/go-spa/backend/user"
)

const (
	pathPrefix   = "/api/v1"
	frontendPath = "../frontend"
	privKey      = "keys/app.rsa"     // openssl genrsa -out app.rsa 2048
	pubKey       = "keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
}

func mustLoadTranslations() {
	i18n.MustLoadTranslationFile("translations/en-us.all.json")
	i18n.MustLoadTranslationFile("translations/pt-br.all.json")
}

// mustLoadEnv loads a .env file with the environment settings
//
// The .env file must have the following structure:
//
// 		# Email settings
// 		export EMAIL_USERNAME=gospa@gmail.com
// 		export EMAIL_PASSWORD=******
//
// 		# Database
// 		DB_USER=gospa
// 		DB_NAME=gospa
// 		DB_PASSWORD=
// 		DB_HOST=127.0.0.1
// 		DB_PORT=5432
// 		DB_SSLMODE=disable
// 		export DB_CONN_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=${DB_SSLMODE}"
//
func mustLoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Unable to load environment settings: %s", err)
	}
}

func main() {
	mustLoadEnv()
	mustLoadTranslations()

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "3000"
	}
	listeningAddr := ":" + port

	router := mux.NewRouter().PathPrefix(pathPrefix).Subrouter()

	db, err := pg.NewSession(os.Getenv("DB_CONN_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}

	vars := map[string]interface{}{"db": db}

	err = ctx.Init(router, privKey, pubKey, vars)
	if err != nil {
		log.Fatalf("Unable to configure API: %s", err)
	}

	server := negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(middleware.WebLogger),
		negroni.NewStatic(http.Dir(frontendPath)),
	)
	server.UseHandler(router)

	log.Infof("Listening on address: %s", listeningAddr)
	log.Fatal(http.ListenAndServe(listeningAddr, server))
}
