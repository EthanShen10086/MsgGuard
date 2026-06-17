package handler

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/EthanShen10086/msgguard/pkg/config"
	"github.com/EthanShen10086/msgguard/pkg/ports"
)

func TestHeuristicClassifySpam(t *testing.T) {
	resp := heuristicClassify("免费贷款中奖 click here")
	if resp.Category != "spam" || resp.Action != "junk" {
		t.Fatalf("got %+v", resp)
	}
}

func TestHeuristicClassifyOTPAllow(t *testing.T) {
	resp := heuristicClassify("您的验证码是 123456")
	if resp.Action != "allow" || resp.Category != "ham" {
		t.Fatalf("got %+v", resp)
	}
}

func TestCategoryToAction(t *testing.T) {
	if categoryToAction("phishing") != "junk" {
		t.Fatal("phishing -> junk")
	}
	if categoryToAction("promotion") != "promotion" {
		t.Fatal("promotion")
	}
}

func TestCacheKeyDeterministic(t *testing.T) {
	a := cacheKey("hello")
	b := cacheKey("hello")
	if a != b || !strings.HasPrefix(a, "classify:") {
		t.Fatalf("keys %q %q", a, b)
	}
}

func TestExtractURLs(t *testing.T) {
	urls := extractURLs("see https://evil.com/path and http://x.y/z.")
	if len(urls) != 2 {
		t.Fatalf("got %v", urls)
	}
}

func TestClassifyHandlerHeuristicWhenLLMDisabled(t *testing.T) {
	h := NewClassifyHandler(nil, nil, nil, nil, &config.Config{}, nil, nil, nil, nil)
	req := httptest.NewRequest("POST", "/api/v1/classify", strings.NewReader(`{"body":"免费贷款"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Classify(w, req)
	if w.Code != 200 {
		t.Fatalf("status %d body %s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "spam") && !strings.Contains(w.Body.String(), "promotion") {
		t.Fatalf("expected heuristic classification: %s", w.Body.String())
	}
}

func TestClassifyHandlerThreatIntelBlocks(t *testing.T) {
	ti := &stubThreatIntel{malicious: true}
	h := NewClassifyHandler(nil, ti, nil, nil, &config.Config{Features: config.Features{CloudLLM: true}}, nil, nil, nil, nil)
	req := httptest.NewRequest("POST", "/api/v1/classify", strings.NewReader(`{"body":"visit https://bad.example"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	h.Classify(w, req)
	if !strings.Contains(w.Body.String(), "phishing") {
		t.Fatalf("expected phishing: %s", w.Body.String())
	}
}

type stubThreatIntel struct {
	malicious bool
}

func (s *stubThreatIntel) CheckURL(ctx context.Context, rawURL string) (*ports.URLVerdict, error) {
	return &ports.URLVerdict{URL: rawURL, Malicious: s.malicious, Source: "test"}, nil
}
