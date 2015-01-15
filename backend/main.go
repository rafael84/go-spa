package main

import (
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/negroni"
	"github.com/gorilla/mux"
	"github.com/gotk/ctx"
	"github.com/gotk/pg"

	"github.com/rafael84/go-spa/backend/cfg"
	_ "github.com/rafael84/go-spa/backend/group"
	_ "github.com/rafael84/go-spa/backend/location"
	_ "github.com/rafael84/go-spa/backend/media"
	_ "github.com/rafael84/go-spa/backend/mediatype"
	"github.com/rafael84/go-spa/backend/middleware"
	_ "github.com/rafael84/go-spa/backend/reset"
	_ "github.com/rafael84/go-spa/backend/user"
)

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
}

func main() {
	cfg.MustLoad()
	listeningAddr := ":" + cfg.Server.Port

	router := mux.NewRouter().PathPrefix(cfg.Server.API.Prefix).Subrouter()

	db, err := pg.NewSession(cfg.DB.ConnectionURL())
	if err != nil {
		log.Fatalf("Unable to connect to database: %s", err)
	}

	vars := map[string]interface{}{"db": db}

	err = ctx.Init(router, cfg.Server.API.PrivKey, cfg.Server.API.PubKey, vars)
	if err != nil {
		log.Fatalf("Unable to configure API: %s", err)
	}

	server := negroni.New(
		negroni.NewRecovery(),
		negroni.HandlerFunc(middleware.WebLogger),
		negroni.NewStatic(http.Dir(cfg.Server.Frontend.Path)),
	)
	server.UseHandler(router)

	log.Infof("Listening on address: %s", listeningAddr)
	log.Fatal(http.ListenAndServe(listeningAddr, server))
}
