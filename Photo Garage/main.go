package main

import (
	"net/http"

	"github.com/Hustle299/Project-0/controllers"
	"github.com/gorilla/mux"
)

func main() {

	r := mux.NewRouter()

	staticsC := controllers.NewStatic()
	r.Handle("/", staticsC.Home).Methods("GET")
	r.Handle("/contact", staticsC.Contact).Methods("GET")

	usersC := controllers.NewUsers()
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	http.ListenAndServe(":8080", r)
}
