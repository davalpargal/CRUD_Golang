package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type App struct {
	DB     *sql.DB
	Router *mux.Router
}

var app App

func main() {
	app = App{}

	app.ConnectToDb("golang")
	defer app.DB.Close()

	app.SetRouter()
	app.StartServer()
}

func (a *App) SetRouter() {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/users", a.AllUsersHandler).Methods("GET")
	a.Router.HandleFunc("/users", a.CreateUserHandler).Methods("POST")
	a.Router.HandleFunc("/user/{username}", a.GetUserHandler).Methods("GET")
	a.Router.HandleFunc("/user/{username}", a.DeleteUserHandler).Methods("DELETE")
	a.Router.HandleFunc("/user/{username}", a.UpdateUserHandler).Methods("PATCH")
}

func (a *App) StartServer() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
