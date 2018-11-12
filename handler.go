package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func (a *App) AllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users := GetAllUsers(a.DB)
	jsonResponse, _ := json.Marshal(users)
	fmt.Fprintf(w, string(jsonResponse))
}

func (a *App) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	payload := r.Body

	var newUser User
	decoder := json.NewDecoder(payload)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&newUser); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Bad Request")
		return
	}

	created := createUser(a.DB, newUser)

	if !created {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Empty Payload")
	} else {
		w.WriteHeader(http.StatusCreated)
		fmt.Fprint(w, "User Created")
	}
}
