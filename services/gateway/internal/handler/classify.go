package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/EthanShen10086/voxera-kit/circuitbreaker"
	"github.com/EthanShen10086/voxera-kit/observability/logger"

	"github.com/EthanShen10086/msgguard/pkg/config"
	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type ClassifyHandler struct {
	classifier ports.LLMClassifier
	breaker    circuitbreaker.CircuitBreaker
	cache      ports.Cache
	queue      ports.Queue
	cfg        *config.Config
	log        logger.Logger
}

func NewClassifyHandler(
	classifier ports.LLMClassifier,
	breaker circuitbreaker.CircuitBreaker,
	cache ports.Cache,
	cfg *config.Config,
	log logger.Logger,
	queue ports.Queue,
) *ClassifyHandler {
	return &ClassifyHandler{classifier: classifier, breaker: breaker, cache: cache, queue: queue, cfg: cfg, log: log}
}

type classifyRequest struct {
	Sender string `json:"sender"`
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

type classifyResponse struct {
	Action     string  `json:"action"`
	Category   string  `json:"category"`
	Confidence float64 `json:"confidence"`
	Cached     bool    `json:"cached,omitempty"`
}

func (h *ClassifyHandler) Classify(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	resp, err := h.run(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, resp)
}

func (h *ClassifyHandler) Defer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	resp, err := h.run(r)
	if err != nil {
		writeJSON(w, classifyResponse{Action: "none", Category: "ham", Confidence: 0})
		return
	}
	writeJSON(w, resp)
}

func (h *ClassifyHandler) run(r *http.Request) (classifyResponse, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return classifyResponse{}, err
	}
	var req classifyRequest
	if err := json.Unmarshal(body, &req); err != nil {
		return classifyResponse{}, err
	}
	return h.runInternal(r, req)
}

func (h *ClassifyHandler) runInternal(r *http.Request, req classifyRequest) (classifyResponse, error) {
	if h.cache != nil {
		key := cacheKey(req.Body)
		if cached, err := h.cache.Get(r.Context(), key); err == nil && cached != nil {
			var resp classifyResponse
			if json.Unmarshal(cached, &resp) == nil {
				resp.Cached = true
				return resp, nil
			}
		}
	}

	if !h.cfg.Features.CloudLLM || h.classifier == nil {
		return heuristicClassify(req.Body), nil
	}

	var category string
	err := h.breaker.Execute(r.Context(), func() error {
		var e error
		category, e = h.classifier.Classify(r.Context(), req.Body)
		return e
	})
	if err != nil {
		return heuristicClassify(req.Body), nil
	}

	action := categoryToAction(category)
	result := classifyResponse{Action: action, Category: category, Confidence: 0.85}

	if h.cache != nil {
		key := cacheKey(req.Body)
		if data, err := json.Marshal(result); err == nil {
			_ = h.cache.Set(context.Background(), key, data, 24*time.Hour)
		}
	}
	if h.queue != nil {
		payload, _ := json.Marshal(map[string]string{"body": req.Body, "category": category})
		_ = h.queue.Publish(r.Context(), "msgguard.classify.done", payload)
	}
	return result, nil
}

func categoryToAction(category string) string {
	switch category {
	case "spam", "phishing":
		return "junk"
	case "promotion":
		return "promotion"
	default:
		return "allow"
	}
}

func cacheKey(body string) string {
	sum := sha256.Sum256([]byte(body))
	return "classify:" + hex.EncodeToString(sum[:8])
}

func heuristicClassify(text string) classifyResponse {
	lower := strings.ToLower(text)
	spamWords := []string{"免费", "中奖", "贷款", "free", "winner", "click here"}
	hits := 0
	for _, w := range spamWords {
		if strings.Contains(lower, w) {
			hits++
		}
	}
	if hits >= 2 {
		return classifyResponse{Action: "junk", Category: "spam", Confidence: 0.9}
	}
	if hits == 1 {
		return classifyResponse{Action: "promotion", Category: "promotion", Confidence: 0.75}
	}
	return classifyResponse{Action: "allow", Category: "ham", Confidence: 0.6}
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}