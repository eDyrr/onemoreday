package main

import (
	"net/http"
	"path/filepath"

	"github.com/eDyrr/onemoreday/cmd/server/routes"
	"github.com/gorilla/mux"
)

func main() {
	// db, _ := database.New("./testDB.db")
	// defer db.Close()

	r := mux.NewRouter()

	staticPath := filepath.Join("static")
	fs := http.FileServer(http.Dir(staticPath))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// r.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
	// 	tmpl, err := template.ParseFiles("./internal/web/templates/users/login.html")
	// 	if err != nil {
	// 		http.Error(w, err.Error(), http.StatusInternalServerError)
	// 		fmt.Print("in the error section")
	// 		return
	// 	}
	// 	tmpl.Execute(w, nil)
	// })

	routes.SetupWebRoutes(r)
	routes.SetupApiRoutes(r)

	http.ListenAndServe(":8080", r)
}
