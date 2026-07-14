package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

/*
TESTING HTTP HANDLERS WITH httptest
===================================

net/http/httptest lets you exercise a handler in memory:

  httptest.NewRequest(method, target, body) -> a fake *http.Request
  httptest.NewRecorder()                    -> a ResponseWriter that records
                                               the status, headers, and body

You call the handler directly with these two and then assert on the recorder.
No server, no port, no network — fast and deterministic.
*/

func TestHello(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	Hello(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if body := rec.Body.String(); body != "Hello World!" {
		t.Errorf("body = %q, want %q", body, "Hello World!")
	}
}

func TestEcho(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("echo payload"))
	rec := httptest.NewRecorder()

	Echo(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("status = %d, want 200", rec.Code)
	}
	body, _ := io.ReadAll(rec.Body)
	if string(body) != "echo payload" {
		t.Errorf("echoed %q, want %q", body, "echo payload")
	}
}

func TestGreetJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/greet?name=Ada", nil)
	rec := httptest.NewRecorder()

	Greet(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}

	var resp map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["message"] != "Hello, Ada!" {
		t.Errorf("message = %q, want %q", resp["message"], "Hello, Ada!")
	}
}

func TestGreetMissingName(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/greet", nil)
	rec := httptest.NewRecorder()

	Greet(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}
