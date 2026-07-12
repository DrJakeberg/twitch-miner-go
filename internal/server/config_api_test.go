package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestConfigServer(t *testing.T) *AnalyticsServer {
	t.Helper()
	return NewAnalyticsServer(":0", newTestLogger(t), nil, "")
}

func configRequest(t *testing.T, s *AnalyticsServer, method, path string, body any) *httptest.ResponseRecorder {
	t.Helper()
	var req *http.Request
	if body != nil {
		b, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rr := httptest.NewRecorder()
	s.srv.Handler.ServeHTTP(rr, req)
	return rr
}

// ── GET /api/config/schema ──────────────────────────────────────────────────

func TestConfigSchema_ReturnsSchema(t *testing.T) {
	s := newTestConfigServer(t)
	rr := configRequest(t, s, "GET", "/api/config/schema", nil)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out["strategies"] == nil {
		t.Error("expected strategies field")
	}
	if out["defaults"] == nil {
		t.Error("expected defaults field")
	}
}

func TestConfigSchema_ContainsAllStrategies(t *testing.T) {
	s := newTestConfigServer(t)
	rr := configRequest(t, s, "GET", "/api/config/schema", nil)
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	strategies, ok := out["strategies"].([]any)
	if !ok {
		t.Fatal("strategies is not an array")
	}
	required := []string{"SMART", "HIGH_ODDS", "MOST_VOTED"}
	for _, r := range required {
		found := false
		for _, s := range strategies {
			if s == r {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("missing required strategy: %s", r)
		}
	}
}

// ── POST /api/config/validate ───────────────────────────────────────────────

func TestConfigValidate_Valid(t *testing.T) {
	s := newTestConfigServer(t)
	rr := configRequest(t, s, "POST", "/api/config/validate", minimalAccountCfg("alice"))
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out["valid"] != true {
		t.Errorf("expected valid=true, got %v", out["valid"])
	}
}

func TestConfigValidate_MissingUsername(t *testing.T) {
	s := newTestConfigServer(t)
	rr := configRequest(t, s, "POST", "/api/config/validate", map[string]any{
		"streamers": []any{map[string]any{"username": "x"}},
	})
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out["valid"] != false {
		t.Error("expected valid=false")
	}
}

func TestConfigValidate_InvalidConfig(t *testing.T) {
	s := newTestConfigServer(t)
	rr := configRequest(t, s, "POST", "/api/config/validate", map[string]any{
		"username": "bad",
		// no streamers, followers, or watchers — fails validation
	})
	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d: %s", rr.Code, rr.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out["valid"] != false {
		t.Error("expected valid=false")
	}
}

func TestConfigValidate_InvalidJSON(t *testing.T) {
	s := newTestConfigServer(t)
	req := httptest.NewRequest("POST", "/api/config/validate", bytes.NewBufferString("{bad"))
	rr := httptest.NewRecorder()
	s.srv.Handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

// ── POST /api/config/generate ───────────────────────────────────────────────

func TestConfigGenerate_Valid(t *testing.T) {
	s := newTestConfigServer(t)
	rr := configRequest(t, s, "POST", "/api/config/generate", minimalAccountCfg("alice"))
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	if out["username"] != "alice" {
		t.Errorf("expected username alice, got %v", out["username"])
	}
	if out["filename"] != "alice.yaml" {
		t.Errorf("expected filename alice.yaml, got %v", out["filename"])
	}
	yamlContent, ok := out["yaml"].(string)
	if !ok || yamlContent == "" {
		t.Error("expected non-empty yaml content")
	}
	if !strings.Contains(yamlContent, "streamers:") {
		t.Error("expected yaml to contain streamers key")
	}
}

func TestConfigGenerate_MissingUsername(t *testing.T) {
	s := newTestConfigServer(t)
	rr := configRequest(t, s, "POST", "/api/config/generate", map[string]any{})
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestConfigGenerate_InvalidConfig(t *testing.T) {
	s := newTestConfigServer(t)
	rr := configRequest(t, s, "POST", "/api/config/generate", map[string]any{
		"username": "bad",
	})
	if rr.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d: %s", rr.Code, rr.Body.String())
	}
}

func TestConfigGenerate_InvalidJSON(t *testing.T) {
	s := newTestConfigServer(t)
	req := httptest.NewRequest("POST", "/api/config/generate", bytes.NewBufferString("{invalid"))
	rr := httptest.NewRecorder()
	s.srv.Handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestConfigGenerate_YAMLContainsAllFields(t *testing.T) {
	s := newTestConfigServer(t)
	cfg := map[string]any{
		"username":          "fulltest",
		"max_watch_streams": 3,
		"priority":          []any{"STREAK", "DROPS"},
		"streamers": []any{
			map[string]any{"username": "streamer1"},
			map[string]any{
				"username": "streamer2",
				"settings": map[string]any{
					"make_predictions": false,
				},
			},
		},
		"followers": map[string]any{
			"enabled": true,
			"order":   "DESC",
		},
	}
	rr := configRequest(t, s, "POST", "/api/config/generate", cfg)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rr.Code, rr.Body.String())
	}
	var out map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &out); err != nil {
		t.Fatal(err)
	}
	yamlContent := out["yaml"].(string)

	checks := []string{"max_watch_streams", "priority", "streamers", "followers"}
	for _, key := range checks {
		if !strings.Contains(yamlContent, key+":") {
			t.Errorf("expected yaml to contain %s", key)
		}
	}
}
