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

func clearDb() {
	clearAllQuery := `DELETE FROM USERS`
	_, _ = a.DB.Exec(clearAllQuery)
}

func TestGetAllUsers(t *testing.T) {
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
	body := []byte(`{}`)
	request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusBadRequest, response.Code)
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	if string(responseBody) != "Empty Payload" {
		t.Errorf("Expected Empty Payload Got %s", responseBody)
	}
}

func TestCreateUserWithIncorrectPayload(t *testing.T) {
	body := []byte(`{"foo":"bar"}`)
	request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusBadRequest {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusBadRequest, response.Code)
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	if string(responseBody) != "Bad Request" {
		t.Errorf("Expected Bad Request Got %s", responseBody)
	}
}

func TestCreateUserWithCorrectPayload(t *testing.T) {
	body := []byte(`{"username":"avd", "email":"avd@gojek.com"}`)
	request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusCreated {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusCreated, response.Code)
	}

	responseBody, _ := ioutil.ReadAll(response.Body)

	if string(responseBody) != "User Created" {
		t.Errorf("Expected User Created Got %s", responseBody)
	}
}

func TestGetUserWithValidUsername(t *testing.T) {
	clearDb()
	userJson := `{"username":"avd","email":"avd@gojek.com"}`
	body := []byte(userJson)
	request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	request, _ = http.NewRequest("GET", "/user/avd", nil)
	response = httptest.NewRecorder()
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
	body := []byte(userJson)
	request, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	response := httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	request, _ = http.NewRequest("DELETE", "/user/avd", nil)
	response = httptest.NewRecorder()
	a.Router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf("Expected Response code %d. Got %d\n", http.StatusOK, response.Code)
	}

	if responseBody := response.Body.String(); responseBody != "User Deleted" {
		t.Errorf("Expected User Deleted, got %s", responseBody)
	}
}
