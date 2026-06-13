package handler

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

// MLModelHealthMonitor tracks model performance metrics reported by clients.
type MLModelHealthMonitor struct {
	mu              sync.Mutex
	totalInferences int64
	totalLatencyMs  float64
	lastReportAt    time.Time
}

var globalHealth = &MLModelHealthMonitor{}

type healthReport struct {
	MeanLatencyMs float64   `json:"mean_latency_ms"`
	TotalCalls    int64     `json:"total_calls"`
	LastReportAt  time.Time `json:"last_report_at,omitempty"`
}

type healthPayload struct {
	LatencyMs float64 `json:"latency_ms"`
	Layer     string  `json:"layer"`
}

func RegisterHealthRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/api/v1/model/health", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			globalHealth.mu.Lock()
			report := healthReport{
				TotalCalls:   globalHealth.totalInferences,
				LastReportAt: globalHealth.lastReportAt,
			}
			if globalHealth.totalInferences > 0 {
				report.MeanLatencyMs = globalHealth.totalLatencyMs / float64(globalHealth.totalInferences)
			}
			globalHealth.mu.Unlock()
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(report)
		case http.MethodPost:
			var p healthPayload
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			globalHealth.mu.Lock()
			globalHealth.totalInferences++
			globalHealth.totalLatencyMs += p.LatencyMs
			globalHealth.lastReportAt = time.Now()
			globalHealth.mu.Unlock()
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
