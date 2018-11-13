package main

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a = App{}

func TestMain(m *testing.M) {
	a.ConnectToDb("testgolang")
	a.SetRouter()
	code := m.Run()
	os.Exit(code)
}

func insertUser(userJson string) (*httptest.ResponseRecorder, error) {
	body := []byte(userJson)
	request, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)
	return response, err
}

func clearDb() {
	clearAllQuery := `DELETE FROM USERS`
	_, _ = a.DB.Exec(clearAllQuery)
}

func TestGetAllUsersForEmptyDatabase(t *testing.T) {
	clearDb()

	req, _ := http.NewRequest("GET", "/users", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, req)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected empty result, got %s", body)
	}
}

func TestCreateUserWithEmptyPayload(t *testing.T) {
	clearDb()

	userJson := `{}`
	response, _ := insertUser(userJson)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusBadRequest, response.Code)
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	if string(responseBody) != "Empty Payload" {
		t.Errorf("Expected Empty Payload Got %s", responseBody)
	}
}

func TestCreateUserWithIncorrectPayload(t *testing.T) {
	clearDb()

	userJson := `{"foo":"bar"}`
	response, _ := insertUser(userJson)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusBadRequest, response.Code)
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	if string(responseBody) != "Bad Request" {
		t.Errorf("Expected Bad Request Got %s", responseBody)
	}
}

func TestCreateUserWithCorrectPayload(t *testing.T) {
	clearDb()

	userJson := `{"username":"avd", "email":"avd@gojek.com"}`
	response, _ := insertUser(userJson)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusCreated, response.Code)
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	if string(responseBody) != "User Created" {
		t.Errorf("Expected User Created Got %s", responseBody)
	}
}

func TestCreateUserWithDuplicatePayload(t *testing.T) {
	clearDb()

	userJson := `{"username":"avd", "email":"avd@gojek.com"}`
	insertUser(userJson)
	response, _ := insertUser(userJson)

	if response.Code != http.StatusConflict {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusConflict, response.Code)
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	if string(responseBody) != "Duplicate Username" {
		t.Errorf("Expected Duplicate Username Got %s", responseBody)
	}
}

func TestGetAllUsersForNonEmptyDatabase(t *testing.T) {
	clearDb()

	userJson := `{"username":"avd","email":"avd@gojek.com"}`
	insertUser(userJson)

	userJson = `{"username":"dav","email":"dav@gojek.com"}`
	insertUser(userJson)

	request, _ := http.NewRequest("GET", "/users", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	resultJson := `[{"username":"avd","email":"avd@gojek.com"},{"username":"dav","email":"dav@gojek.com"}]`
	if body := response.Body.String(); body != resultJson {
		t.Errorf("Expected empty result, got %s", body)
	}
}

func TestGetUserWithValidUsername(t *testing.T) {
	clearDb()
	userJson := `{"username":"avd","email":"avd@gojek.com"}`
	insertUser(userJson)

	request, _ := http.NewRequest("GET", "/user/avd", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	if responseBody := response.Body.String(); responseBody != userJson {
		t.Errorf("Expected %s, got %s", userJson, responseBody)
	}
}

func TestGetUserWithInvalidUsername(t *testing.T) {
	clearDb()

	request, _ := http.NewRequest("GET", "/user/avdx", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusNotFound, response.Code)
	}

	if responseBody := response.Body.String(); responseBody != "{}" {
		t.Errorf("Expected {}, got %s", responseBody)
	}
}

func TestDeleteUserWithValidUsername(t *testing.T) {
	clearDb()

	userJson := `{"username":"avd","email":"avd@gojek.com"}`
	insertUser(userJson)

	request, _ := http.NewRequest("DELETE", "/user/avd", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	if responseBody := response.Body.String(); responseBody != "User Deleted" {
		t.Errorf("Expected User Deleted, got %s", responseBody)
	}
}

func TestDeleteUserWithInvalidUsername(t *testing.T) {
	clearDb()

	request, _ := http.NewRequest("DELETE", "/user/avd", nil)
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusNotFound, response.Code)
	}

	if responseBody := response.Body.String(); responseBody != "User Not Found" {
		t.Errorf("Expected User Not Found, got %s", responseBody)
	}
}

func TestUpdateUserForValidUsername(t *testing.T) {
	clearDb()

	userInitialJson := `{"username":"avd","email":"avdinitial@gojek.com"}`
	userUpdateJson := `{"email":"avdfinal@gojek.com"}`
	insertUser(userInitialJson)

	body := []byte(userUpdateJson)
	request, _ := http.NewRequest("PATCH", "/user/avd", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	userFinalJson := `{"username":"avd","email":"avdfinal@gojek.com"}`

	if responseBody := response.Body.String(); responseBody != userFinalJson {
		t.Errorf("Expected %s, got %s", userFinalJson, responseBody)
	}
}

func TestUpdateUserForInvalidUsername(t *testing.T) {
	clearDb()
	userUpdateJson := `{"email":"avdfinal@gojek.com"}`

	body := []byte(userUpdateJson)
	request, _ := http.NewRequest("PATCH", "/user/avd", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusBadRequest, response.Code)
	}

	if responseBody := response.Body.String(); responseBody != "No such user exists" {
		t.Errorf("Expected no such user exists, got %s", responseBody)
	}
}

func TestUpdateUserForEmptyPayload(t *testing.T) {
	clearDb()

	userInitialJson := `{"username":"avd","email":"avdinitial@gojek.com"}`
	userUpdateJson := `{}`
	insertUser(userInitialJson)

	body := []byte(userUpdateJson)
	request, _ := http.NewRequest("PATCH", "/user/avd", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusBadRequest, response.Code)
	}

	if responseBody := response.Body.String(); responseBody != "Empty Payload" {
		t.Errorf("Expected Empty Payload, got %s", responseBody)
	}
}
