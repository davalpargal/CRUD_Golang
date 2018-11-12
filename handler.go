package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func (a *App) AllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := GetAllUsers(a.DB)
	jsonResponse, _ := json.Marshal(users)
	fmt.Fprintf(w, string(jsonResponse))
}

func (a *App) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("Request Body could not be read")
		return
	}
	var newUser User

	json.Unmarshal(payload, &newUser)

	created := createUser(a.DB, newUser)

	if !created {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Empty Payload")
	}
}
