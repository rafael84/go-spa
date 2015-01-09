package main

import (
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/i18n"

	_ "github.com/rafael84/go-spa/backend/account"
	"github.com/rafael84/go-spa/backend/api"
	"github.com/rafael84/go-spa/backend/database"
	"github.com/rafael84/go-spa/backend/middleware"
	_ "github.com/rafael84/go-spa/backend/storage"
)

const (
	pathPrefix   = "/api/v1"
	frontendPath = "../frontend"
)

func init() {
	log.SetLevel(log.DebugLevel)
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

	router := mux.NewRouter()

	db, err := database.NewSession(os.Getenv("DB_CONN_URL"))
	if err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}

	err = api.Configure(router, pathPrefix, db)
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
