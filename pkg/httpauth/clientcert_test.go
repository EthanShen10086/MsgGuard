package httpauth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequireClientCert_skipsPublicPaths(t *testing.T) {
	next := RequireClientCert([]string{"/api/v1/admin/"})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()
	next.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}

func TestRequireClientCert_blocksAdminWithoutCert(t *testing.T) {
	next := RequireClientCert([]string{"/api/v1/admin/"})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/metrics/summary", nil)
	rr := httptest.NewRecorder()
	next.ServeHTTP(rr, req)
	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", rr.Code)
	}
}

func TestRequireClientCertHeader(t *testing.T) {
	next := RequireClientCertHeader("X-Client-Cert-Subject", []string{"/api/v1/admin/"})(
		http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { w.WriteHeader(http.StatusOK) }),
	)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/admin/metrics/summary", nil)
	req.Header.Set("X-Client-Cert-Subject", "CN=admin")
	rr := httptest.NewRecorder()
	next.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
}
