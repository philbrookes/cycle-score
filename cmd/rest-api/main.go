package main

import (
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"log"
	"github.com/philbrookes/cycle-score/pkg/config"
	"github.com/philbrookes/cycle-score/pkg/controller"
)

func main() {
	router := mux.NewRouter().StrictSlash(true)

	originsOk := handlers.AllowedOrigins(config.GetConfig().GetAllowedOrigins())
	methodsOk := handlers.AllowedMethods(config.GetConfig().GetAllowedMethods())

	controller.ConfigureScore(router.PathPrefix("/api/score").Subrouter())
	controller.ConfigureAuth(router.PathPrefix("/api/auth").Subrouter())

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("/usr/src/web/public/")))
	http.Handle("/", router)

	log.Fatal(http.ListenAndServe(config.GetConfig().GetPortListenerStr(), handlers.CORS(originsOk, methodsOk)(router)))
}