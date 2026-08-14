package web

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestFrontendAndAPIRoutes(t *testing.T) {
	oldMux := http.DefaultServeMux
	http.DefaultServeMux = http.NewServeMux()
	defer func() { http.DefaultServeMux = oldMux }()

	RegisterRoutes()

	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	http.DefaultServeMux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET / status = %d, want %d", rec.Code, http.StatusOK)
	}
	body := rec.Body.String()
	if !strings.Contains(body, `id="root"`) || !strings.Contains(body, "/assets/") {
		t.Fatalf("GET / did not serve built React index: %q", body)
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/state", nil)
	http.DefaultServeMux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/state status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !strings.Contains(rec.Body.String(), `"max":20`) {
		t.Fatalf("GET /api/state did not return HP JSON: %q", rec.Body.String())
	}

	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/defense", nil)
	http.DefaultServeMux.ServeHTTP(rec, req)
	if rec.Code != http.StatusFound {
		t.Fatalf("GET /defense status = %d, want %d", rec.Code, http.StatusFound)
	}

	assets, err := fs.Glob(frontendFS, "static/dist/assets/*.css")
	if err != nil || len(assets) == 0 {
		t.Fatalf("built CSS asset not found: assets=%v err=%v", assets, err)
	}
	assetPath := strings.TrimPrefix(assets[0], "static/dist")
	rec = httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, assetPath, nil)
	http.DefaultServeMux.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET %s status = %d, want %d", assetPath, rec.Code, http.StatusOK)
	}
}
