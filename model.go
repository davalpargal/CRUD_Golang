package main

import (
	"database/sql"
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

func createUser(db *sql.DB, newUser User) (created bool, err error) {
	if newUser.Username == "" && newUser.Email == "" {
		created = false
		err = nil
		return
	}
	query := `
	INSERT INTO USERS(USERNAME, EMAIL) VALUES($1, $2)`
	response, err := db.Exec(query, newUser.Username, newUser.Email)

	if err != nil {
		created = false
		return
	}

	if rowsChanged, _ := response.RowsAffected(); rowsChanged == 1 {
		created = true
	}
	return
}

func getUserWithUsername(db *sql.DB, username string) (user User, err error) {
	query := `SELECT * FROM USERS WHERE USERNAME = $1`

	row := db.QueryRow(query, username)

	err = row.Scan(&user.Username, &user.Email)
	return
}

func deleteUserWithUsername(db *sql.DB, username string) (deleted bool) {
	query := `DELETE FROM USERS WHERE USERNAME = $1`

	response, _ := db.Exec(query, username)

	count, err := response.RowsAffected()
	if err != nil {
		panic(err)
	}
	if count == 1 {
		deleted = true
	} else {
		deleted = false
	}
	return
}

func updateEmailWithUsername(db *sql.DB, username string, email interface{}) (updated bool) {
	query := `UPDATE USERS
SET EMAIL = $1
WHERE USERNAME = $2`

	response, _ := db.Exec(query, email, username)
	if rowsChanged, _ := response.RowsAffected(); rowsChanged == 1 {
		updated = true
	} else if rowsChanged == 0 {
		updated = false
	}
	return
}
