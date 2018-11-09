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
