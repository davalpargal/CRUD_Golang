package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

func (a *App) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	user, err := getUserWithUsername(a.DB, username)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "{}")
		return
	}

	userJson, _ := json.Marshal(user)
	fmt.Fprint(w, string(userJson))
}

func (a *App) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]
	deleted := deleteUserWithUsername(a.DB, username)

	if deleted {
		fmt.Fprint(w, "User Deleted")
	}
}
