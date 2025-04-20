package web

import (
	"fmt"
	"html/template"
	"net/http"
)

func Register(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/users/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Print("in the error section")
		return
	}
	tmpl.Execute(w, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/users/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func Home(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}
