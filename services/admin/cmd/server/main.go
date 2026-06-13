// Package main is the MsgGuard admin API for quota whitelist and feature flags.
package main

import (
	"encoding/json"
	"net/http"
	"os"

	aiquotaMemory "github.com/EthanShen10086/voxera-kit/aiquota/memory"
	"github.com/EthanShen10086/voxera-kit/aiquota"
	ffMemory "github.com/EthanShen10086/voxera-kit/featureflag/memory"
	memadapters "github.com/EthanShen10086/msgguard/pkg/adapters/memory"
)

func main() {
	auth := memadapters.NewAuth(os.Getenv("AUTH_SECRET"))
	quota := aiquotaMemory.NewStore()
	flags := ffMemory.NewAdapter()

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	mux.HandleFunc("/api/v1/admin/quota/whitelist", func(w http.ResponseWriter, r *http.Request) {
		claims, err := auth.Authenticate(r.Context(), r.Header.Get("Authorization"))
		if err != nil || !auth.HasRole(r.Context(), claims, "admin") {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		if r.Method == http.MethodPost {
			var req struct {
				UserID string `json:"user_id"`
			}
			_ = json.NewDecoder(r.Body).Decode(&req)
			_ = quota.AddToWhitelist(r.Context(), aiquota.WhitelistEntry{UserID: req.UserID, Reason: "admin grant"})
			w.WriteHeader(http.StatusCreated)
			return
		}
		w.Write([]byte(`{"status":"ok"}`))
	})
	mux.HandleFunc("/api/v1/admin/flags", func(w http.ResponseWriter, r *http.Request) {
		claims, err := auth.Authenticate(r.Context(), r.Header.Get("Authorization"))
		if err != nil || !auth.HasRole(r.Context(), claims, "admin") {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		_ = flags
		json.NewEncoder(w).Encode(map[string]string{"cloud_llm": "config-driven"})
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	http.ListenAndServe(":"+port, mux)
}
