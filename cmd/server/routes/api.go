package routes

import (
	"github.com/eDyrr/onemoreday/cmd/server/handlers/api"
	"github.com/gorilla/mux"
)

func SetupApiRoutes(router *mux.Router) {
	apiRouter := router.PathPrefix("/api").Subrouter()

	apiRouter.HandleFunc("/login", api.Login).Methods("POST")
	apiRouter.HandleFunc("/register", api.Register).Methods("POST")
}
