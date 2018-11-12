package main

import (
	"database/sql"
	"fmt"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

func GetAllUsers(db *sql.DB) (users []User) {
	query := `
	SELECT * FROM USERS`
	rows, err := db.Query(query)

	if err != nil {
		fmt.Println("query")
		panic(err)
	}

	users = make([]User, 0)

	for rows.Next() {
		user := User{}
		rows.Scan(&user.Username, &user.Email)
		users = append(users, user)
	}

	defer rows.Close()
	return
}

func createUser(db *sql.DB, newUser User) (created bool) {
	if newUser.Username == "" && newUser.Email == "" {
		return false
	}
	query := `
	INSERT INTO USERS(USERNAME, EMAIL) VALUES($1, $2)`
	response, err := db.Exec(query, newUser.Username, newUser.Email)

	if err != nil {
		created = false
	}

	if rowsChanged, _ := response.RowsAffected(); rowsChanged == 1 {
		created = true
	} else {
		created = false
	}
	return
}

func getUserWithUsername(db *sql.DB, username string) (user User, err error) {
	query := `SELECT * FROM USERS WHERE USERNAME = $1`

	row := db.QueryRow(query, username)

	err = row.Scan(&user.Username, &user.Email)
	return
}
