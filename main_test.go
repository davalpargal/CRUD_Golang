package main_test

import (
	"os"
	"testing"

	"CRUD"
)

func TestMain(m *testing.M) {
	var a = main.App{}
	a.ConnectToDb("testgolang")
	code := m.Run()
	os.Exit(code)
}