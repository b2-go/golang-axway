package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/gin-gonic/gin"
	"log"
)

func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var m map[string]string
	m = make(map[string]string)
	err := json.Unmarshal(w.Body.Bytes(), &m)
	assert.Equal(t, nil, err)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "pong", m["message"])
}

func TestGet404Route(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/keyvalue/key", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 404, w.Code)
	assert.Equal(t, "not found", w.Body.String())
}

func TestGet200Route(t *testing.T) {
	router := setupRouter()
	testPostThenGet(t, "value", router)
}

func TestGetPutGet200Route(t *testing.T) {
	router := setupRouter()
	testPostThenGet(t, "value", router)
	testPostThenGet(t, "value2", router)
}

func testPostThenGet(t *testing.T, value string, router *gin.Engine) {
	wPost := httptest.NewRecorder()
	reqPost, _ := http.NewRequest("POST", "/v1/keyvalue?key=key&value=" + value, nil)
	router.ServeHTTP(wPost, reqPost)

	assert.Equal(t, 200, wPost.Code)
	assert.Equal(t, "ok: value set", wPost.Body.String())

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/keyvalue/key", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, value, w.Body.String())
}