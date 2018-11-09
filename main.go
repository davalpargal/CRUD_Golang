package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	DB *sql.DB
	Router *mux.Router
}

var app App

func main () {
	app = App{}

	app.ConnectToDb("testgolang")
	defer app.DB.Close()

	app.SetRouter()
	app.StartServer()
}

func (a *App) SetRouter() {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/users", a.AllUsersHandler)
}

func (a *App) StartServer() {
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
