package routes

import (
	"github.com/eDyrr/onemoreday/cmd/server/handlers/web"
	"github.com/gorilla/mux"
)

func SetupWebRoutes(router *mux.Router) {

	router.HandleFunc("/", web.Register)
	router.HandleFunc("/login", web.Login)
}
