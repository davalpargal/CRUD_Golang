#CRUD APP

## Project Description
It is a CRUD app with username and email in Golang.

### Developer Steps
* Login into psql and create a db named Golang and a table named Users.

* Table consist of two columns username and email. Username is Unique field.

* Install lib/pq package by
```go get -u github.com/lib/pq```

*Install mux for Routing by
```go get -u github.com/gorilla/mux```

* Install the package by 
```go install``` inside the directory

Note: It will install the project as a binary with name as the project directory. See the workspace bin for the binary name. 

*Run the server using command
```go run [ProjectBinaryName]```

### Testing
* Create Database named testgolang and table users with same schema as above. 

* Run the tests using
``` go test ```

* For code coverage
``` go test -coverprofile coverage.out ```

* Convert to html using
```  go tool cover -html=coverage.out ```