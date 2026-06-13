package handler

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/EthanShen10086/voxera-kit/aiquota"
	ff "github.com/EthanShen10086/voxera-kit/featureflag"
	"github.com/EthanShen10086/voxera-kit/llm"
	"github.com/EthanShen10086/voxera-kit/llm/prompt"
	"github.com/EthanShen10086/voxera-kit/observability/logger"

	"github.com/EthanShen10086/msgguard/pkg/config"
	"github.com/EthanShen10086/msgguard/pkg/ports"
)

type ClassifyHandler struct {
	router    *llm.Router
	quota     aiquota.Manager
	flags     ff.Store
	cfg       *config.Config
	log       logger.Logger
	cache     ports.Cache
	mu        sync.Mutex
	failCount int
	openUntil time.Time
}

func NewClassifyHandler(router *llm.Router, quota aiquota.Manager, flags ff.Store, cfg *config.Config, log logger.Logger, cache ports.Cache) *ClassifyHandler {
	return &ClassifyHandler{router: router, quota: quota, flags: flags, cfg: cfg, log: log, cache: cache}
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

	if !h.cfg.Features.CloudLLM {
		return heuristicClassify(req.Body), nil
	}

	if h.isCircuitOpen() {
		return heuristicClassify(req.Body), nil
	}

	if h.router == nil {
		return heuristicClassify(req.Body), nil
	}

	tmpl := prompt.Classify
	_, userPrompt := tmpl.Render(map[string]any{
		"Categories": "spam, promotion, ham, phishing, transaction",
		"Text":       req.Body,
	})
	resp, err := h.router.Route(r.Context(), llm.Request{
		Messages: []llm.Message{
			{Role: llm.RoleSystem, Content: tmpl.System},
			{Role: llm.RoleUser, Content: userPrompt},
		},
	})
	if err != nil {
		h.recordFailure()
		return heuristicClassify(req.Body), nil
	}
	h.recordSuccess()
	category := strings.TrimSpace(strings.ToLower(resp.Content))
	action := "allow"
	switch category {
	case "spam", "phishing":
		action = "junk"
	case "promotion":
		action = "promotion"
	}
	result := classifyResponse{Action: action, Category: category, Confidence: 0.85}
	if h.cache != nil {
		key := cacheKey(req.Body)
		if data, err := json.Marshal(result); err == nil {
			_ = h.cache.Set(context.Background(), key, data, 24*time.Hour)
		}
	}
	return result, nil
}

func (h *ClassifyHandler) isCircuitOpen() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return time.Now().Before(h.openUntil)
}

func (h *ClassifyHandler) recordFailure() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.failCount++
	if h.failCount >= 3 {
		h.openUntil = time.Now().Add(30 * time.Second)
		h.failCount = 0
	}
}

func (h *ClassifyHandler) recordSuccess() {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.failCount = 0
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
