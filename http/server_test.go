// Copyright 2018 The Moov Authors
// Use of this source code is governed by an Apache License
// license that can be found in the LICENSE file.

package http

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestHTTP__AddCORSHandler(t *testing.T) {
	router := mux.NewRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("OPTIONS", "https://api.moov.io/v1/auth/ping", nil)
	r.Header.Set("Origin", "https://moov.io")

	AddCORSHandler(router)
	router.ServeHTTP(w, r)
	w.Flush()

	if w.Code != 200 {
		t.Errorf("got %d", w.Code)
	}
	if v := w.Header().Get("Access-Control-Allow-Origin"); v != "https://moov.io" {
		t.Errorf("got %q", v)
	}
	headers := []string{
		"Access-Control-Allow-Methods",
		"Access-Control-Allow-Headers",
		"Access-Control-Allow-Credentials",
		"Content-Type",
	}
	for i := range headers {
		v := w.Header().Get(headers[i])
		if v == "" {
			t.Errorf("%s's value is an empty string", headers[i])
		}
	}
}

func TestHTTP__emptyOrigin(t *testing.T) {
	router := mux.NewRouter()
	w := httptest.NewRecorder()
	r := httptest.NewRequest("OPTIONS", "https://api.moov.io/v1/auth/ping", nil)
	r.Header.Set("Origin", "")

	AddCORSHandler(router)
	router.ServeHTTP(w, r)
	w.Flush()

	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d", w.Code)
	}
}

func TestHTTP__Problem(t *testing.T) {
	w := httptest.NewRecorder()
	Problem(w, errors.New("problem X"))
	w.Flush()

	// check http response
	if w.Code != http.StatusBadRequest {
		t.Errorf("got %d", w.Code)
	}
	v := w.Header().Get("Content-Type")
	if !strings.Contains(v, "application/json") {
		t.Errorf("got %s", v)
	}

	type resp struct {
		Error string `json:"error"`
	}
	var response resp
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Error(err)
	}
	if response.Error != "problem X" {
		t.Errorf("got %q", response.Error)
	}
}

func TestHTTP_InternalError(t *testing.T) {
	w := httptest.NewRecorder()
	where := InternalError(w, errors.New("problem Y"))
	w.Flush()

	if w.Code != http.StatusInternalServerError {
		t.Errorf("got %d", w.Code)
	}

	if !strings.HasPrefix(where, "server_test.go") { // This will always be this file's name
		t.Errorf("got %s", where)
	}
}
