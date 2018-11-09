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
