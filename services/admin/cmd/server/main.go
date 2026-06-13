// Package main is the MsgGuard admin API for quota whitelist and feature flags.
package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/EthanShen10086/voxera-kit/aiquota"
	aiquotaMemory "github.com/EthanShen10086/voxera-kit/aiquota/memory"
	"github.com/EthanShen10086/voxera-kit/featureflag"
	ffMemory "github.com/EthanShen10086/voxera-kit/featureflag/memory"
	memadapters "github.com/EthanShen10086/msgguard/pkg/adapters/memory"
)

func main() {
	auth := memadapters.NewAuth(os.Getenv("AUTH_SECRET"))
	quota := aiquotaMemory.NewStore()
	flags := ffMemory.NewAdapter()
	_ = flags.SetFlag(nil, featureflag.Flag{Key: "cloud_llm", Enabled: true, Percentage: 100})

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	mux.HandleFunc("/api/v1/admin/quota/whitelist", func(w http.ResponseWriter, r *http.Request) {
		claims, err := auth.Authenticate(r.Context(), r.Header.Get("Authorization"))
		if err != nil || !auth.HasRole(r.Context(), claims, "admin") {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		switch r.Method {
		case http.MethodGet:
			entries, _ := quota.ListWhitelist(r.Context())
			json.NewEncoder(w).Encode(entries)
		case http.MethodPost:
			var req struct {
				UserID string `json:"user_id"`
				Reason string `json:"reason"`
			}
			_ = json.NewDecoder(r.Body).Decode(&req)
			_ = quota.AddToWhitelist(r.Context(), aiquota.WhitelistEntry{UserID: req.UserID, Reason: req.Reason})
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]string{"status": "whitelisted"})
		}
	})
	mux.HandleFunc("/api/v1/admin/flags", func(w http.ResponseWriter, r *http.Request) {
		claims, err := auth.Authenticate(r.Context(), r.Header.Get("Authorization"))
		if err != nil || !auth.HasRole(r.Context(), claims, "admin") {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}
		switch r.Method {
		case http.MethodGet:
			list, _ := flags.GetFlags(r.Context())
			json.NewEncoder(w).Encode(list)
		case http.MethodPost, http.MethodPut:
			var flag featureflag.Flag
			_ = json.NewDecoder(r.Body).Decode(&flag)
			_ = flags.SetFlag(r.Context(), flag)
			json.NewEncoder(w).Encode(map[string]string{"status": "updated"})
		}
	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	http.ListenAndServe(":"+port, mux)
}
