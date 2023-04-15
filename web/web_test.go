package web

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIndex(t *testing.T) {
	router := NewWeb()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)
	router.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error("Status != 200")
	}
}

func TestNotifiersView(t *testing.T) {
	router := NewWeb()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/notifiers/view", nil)
	router.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error("Status != 200")
	}
}

func TestSchedulesView(t *testing.T) {
	router := NewWeb()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/schedules/view", nil)
	router.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error("Status != 200")
	}
}

func TestCreateWatchGet(t *testing.T) {
	router := NewWeb()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/watch/create", nil)
	router.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error("Status != 200")
	}
}

func TestCreateBackupView(t *testing.T) {
	router := NewWeb()
	w := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/backup/view", nil)
	router.router.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Error("Status != 200")
	}
}
