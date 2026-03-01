package updater

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckForUpdate_NewerAvailable(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ghRelease{TagName: "v1.2.0"})
	}))
	defer srv.Close()

	info, err := checkWithURL(context.Background(), "1.0.0", srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !info.Available {
		t.Error("expected update to be available")
	}
	if info.Latest != "1.2.0" {
		t.Errorf("Latest = %q, want %q", info.Latest, "1.2.0")
	}
}

func TestCheckForUpdate_UpToDate(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ghRelease{TagName: "v1.0.0"})
	}))
	defer srv.Close()

	info, err := checkWithURL(context.Background(), "1.0.0", srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Available {
		t.Error("expected no update available")
	}
}

func TestCheckForUpdate_DevVersion(t *testing.T) {
	info, err := checkWithURL(context.Background(), "dev", "http://unused")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if info.Available {
		t.Error("expected no update for dev version")
	}
}

func TestCheckForUpdate_ServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	_, err := checkWithURL(context.Background(), "1.0.0", srv.URL)
	if err == nil {
		t.Error("expected error on server failure")
	}
}
