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

func main () {
	app := App{}

	app.connectToDb()
	defer app.DB.Close()

	app.startServer()
}

func (a *App) startServer() {
	a.Router = mux.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", a.Router))
}
