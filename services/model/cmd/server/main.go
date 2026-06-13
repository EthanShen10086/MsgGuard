// Package main serves Core ML model metadata and bundles.
package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	fsadapter "github.com/EthanShen10086/msgguard/pkg/adapters/filesystem"
	memadapters "github.com/EthanShen10086/msgguard/pkg/adapters/memory"
	"github.com/EthanShen10086/msgguard/pkg/httpauth"
	"github.com/EthanShen10086/msgguard/pkg/ports"
)

func main() {
	dir := envOr("MODEL_STORAGE_PATH", "../../deploy/models")
	reg, err := fsadapter.NewModelRegistry(dir)
	if err != nil {
		panic(err)
	}
	auth := memadapters.NewAuth(envOr("AUTH_SECRET", "msgguard-dev-secret"))
	srv := &server{registry: reg, dir: dir}
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	mux.HandleFunc("/api/v1/models/latest", srv.latest)
	mux.HandleFunc("/api/v1/models/register", httpauth.RequirePermission(auth, auth, "models", "write", srv.register))
	mux.HandleFunc("/api/v1/models/", srv.download)
	port := envOr("PORT", "8083")
	http.ListenAndServe(":"+port, mux)
}

type server struct {
	registry ports.ModelRegistry
	dir      string
	mu       sync.Mutex
}

func (s *server) latest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	locale := r.URL.Query().Get("locale")
	if locale == "" {
		locale = "zh-Hans"
	}
	meta, err := s.registry.GetLatest(r.Context(), locale)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("ETag", meta.Checksum)
	json.NewEncoder(w).Encode(map[string]any{
		"version": meta.Version, "locale": meta.Locale,
		"checksum": meta.Checksum,
		"url":      "/api/v1/models/" + meta.Version + "/download/spam_classifier.mlmodel",
		"bayes_url": "/api/v1/models/" + meta.Version + "/download/bayes_model.json",
		"minIOS":   "17.0",
	})
}

func (s *server) register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, _ := io.ReadAll(r.Body)
	var meta ports.ModelMeta
	if err := json.Unmarshal(body, &meta); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := s.registry.Register(r.Context(), meta); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "registered"})
}

func (s *server) download(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/models/"), "/")
	if len(parts) < 3 || parts[1] != "download" {
		http.NotFound(w, r)
		return
	}
	version, name := parts[0], parts[2]
	data, err := s.registry.GetArtifact(r.Context(), version, name)
	if err != nil || data == nil {
		for _, loc := range []string{"zh-Hans", "en-US"} {
			fallback := filepath.Join(s.dir, loc, version, name)
			data, err = os.ReadFile(fallback)
			if err == nil {
				break
			}
		}
	}
	if err != nil || data == nil {
		fallback := filepath.Join(s.dir, "..", "..", "ml", "output", name)
		data, err = os.ReadFile(fallback)
	}
	if err != nil {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(data)
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
