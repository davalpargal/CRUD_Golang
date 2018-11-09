package main

import (
	"log"
	"net/http"
)

func main () {
	db, err := connectToDb()
	if err != nil {
		panic(err)
	}

	defer db.Close()

	startServer()
}

func startServer() {
	log.Fatal(http.ListenAndServe(":8080", nil))
}
