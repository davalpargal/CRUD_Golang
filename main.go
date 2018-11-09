package main

func main () {
	db, err := connectToDb()
	if err != nil {
		panic(err)
	}

	defer db.Close()

}
