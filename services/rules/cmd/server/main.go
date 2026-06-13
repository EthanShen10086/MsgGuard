// Package main serves spam filter rules.
package main

import (
	"encoding/json"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) { w.Write([]byte("ok")) })
	mux.HandleFunc("/api/v1/rules/latest", func(w http.ResponseWriter, r *http.Request) {
		rules := map[string]any{
			"version": "1.0.0",
			"locale":  "zh-Hans",
			"tags":    []string{"advertising", "charity", "shortCode", "finance", "order"},
			"keywords_block": []string{"免费领取", "中奖", "贷款", "free gift", "winner"},
			"keywords_allow": []string{"验证码", "verification code", "取件码"},
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("ETag", `"rules-v1"`)
		json.NewEncoder(w).Encode(rules)
	})
	port := envOr("PORT", "8081")
	http.ListenAndServe(":"+port, mux)
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}
