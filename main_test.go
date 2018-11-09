package main_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"CRUD"
)
var a = main.App{}

func TestMain(m *testing.M) {
	a.ConnectToDb("testgolang")
	a.SetRouter()
	code := m.Run()
	os.Exit(code)
}

func TestGetAllUsers(t *testing.T) {
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