#CRUD APP

## Project Description
It is a CRUD app with username and email in Golang.

### Developer Steps
* Login into psql and create a db named Golang and a table named Users.

* Table consist of two columns username and email.

* Install lib/pq package by
```go get -u github.com/lib/pq```

*Install mux for Routing by
```go get -u github.com/gorilla/mux```

*Run the server using command
```go run main.go db_connection.go model.go handler.go```

### Testing
* Run the tests using
``` go test ```

* For code coverage
``` go test -coverprofile coverage.out ```

* Convert to html using
```  go tool cover -html=coverage.out ```