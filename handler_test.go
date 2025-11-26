package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	req := httptest.NewRequest("GET", "/hello", nil)
	w := httptest.NewRecorder()

	HelloHandler(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf(`expected "200", got %d`, w.Code)
	}

	var resp HelloResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if resp.Message != "hello from prod-ready service" {
		t.Errorf(`expected "hello from prod-ready" service, got %s`, resp.Message)
	}

	t.Log("test passed")

}
