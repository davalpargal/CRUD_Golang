package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	duplicateDBEntry = `pq: duplicate key value violates unique constraint "users_username_key"`
)

type UpdatePayload struct {
	Email string  `json:email`
}


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

	created, err := createUser(a.DB, newUser)

	if !created {
		if err != nil && err.Error() == duplicateDBEntry {
			w.WriteHeader(http.StatusConflict)
			fmt.Fprint(w, "Duplicate Username")
			return
		}
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
	} else {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "User Not Found")
	}
}

func (a *App) UpdateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	username := params["username"]

	var updateData UpdatePayload

	payload := r.Body
	decoder := json.NewDecoder(payload)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&updateData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Incorrect Payload")
		return
	}

	if updateData.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Empty Payload")
		return
	}

	updated := updateEmailWithUsername(a.DB, username, updateData.Email)

	if updated {
		var updatedUser User
		updatedUser.Username = username
		updatedUser.Email = updateData.Email

		w.WriteHeader(http.StatusOK)
		jsonResponse, _ := json.Marshal(updatedUser)
		fmt.Fprintf(w, string(jsonResponse))
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "No such user exists")
	}
}
