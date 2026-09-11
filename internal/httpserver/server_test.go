package httpserver

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

type checkFunc func(context.Context) error

func (f checkFunc) Ping(ctx context.Context) error { return f(ctx) }

var testAssets = fstest.MapFS{
	"templates/layout.html": {Data: []byte(`{{define "layout"}}<html>{{template "content" .}}</html>{{end}}`)},
	"templates/home.html":   {Data: []byte(`{{define "content"}}<h1>{{.Title}}</h1>{{end}}`)},
	"static/css/app.css":    {Data: []byte(`body { color: white; }`)},
}

func TestHomeAndHealthRoutes(t *testing.T) {
	handler, err := NewHandler(checkFunc(func(context.Context) error { return nil }), testAssets)
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}

	for _, test := range []struct {
		path       string
		wantStatus int
		wantBody   string
	}{
		{"/", http.StatusOK, "Hortons Diary"},
		{"/health/live", http.StatusOK, ""},
		{"/health/ready", http.StatusOK, ""},
		{"/static/css/app.css", http.StatusOK, "color: white"},
	} {
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, test.path, nil))
		if recorder.Code != test.wantStatus {
			t.Errorf("GET %s status = %d, want %d", test.path, recorder.Code, test.wantStatus)
		}
		if !strings.Contains(recorder.Body.String(), test.wantBody) {
			t.Errorf("GET %s body = %q, want it to contain %q", test.path, recorder.Body.String(), test.wantBody)
		}
	}
}

func TestReadinessFailsWhenDatabaseIsUnavailable(t *testing.T) {
	handler, err := NewHandler(checkFunc(func(context.Context) error { return errors.New("unavailable") }), testAssets)
	if err != nil {
		t.Fatalf("NewHandler: %v", err)
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/health/ready", nil))
	if recorder.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusServiceUnavailable)
	}
}
