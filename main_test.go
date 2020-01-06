package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jinzhu/gorm"
)

func TestLoginRoute(t *testing.T) {
	//--- Setup ---//
	db := setupDB(t)
	defer db.Close()

	handlersUser := HandlersUserApi(db)
	engine := router(handlersUser)

	//--- Setup ---//
	data := url.Values{}

	data.Add("username", "user")
	data.Add("password", "password")

	// Mock payload
	payload := strings.NewReader("")

	// Mock httpWriter and normal request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/login", payload)
	if err != nil {
		t.Log(err)
	}

	engine.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, fmt.Sprintf("Logged in %s", data["username"]), w.Body.String())
}

func setupDB(t *testing.T) *gorm.DB {

	conn := &database{
		host:     "localhost",
		port:     5432,
		user:     "userapp",
		password: "password",
		dbName:   "postgres",
	}

	// Opens the connection to the db
	db, err := gorm.Open("postgres", conn.String())
	if err != nil {
		t.Log(err)
	}

	return db
}
