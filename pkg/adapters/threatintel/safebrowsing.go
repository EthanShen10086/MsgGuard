package threatintel

import (
	"context"
	"os"

	"github.com/EthanShen10086/msgguard/pkg/ports"
)

// SafeBrowsing is a stub adapter for Google Safe Browsing URL checks.
// Set GOOGLE_SAFE_BROWSING_API_KEY to enable; without a key all URLs are treated as safe.
type SafeBrowsing struct {
	apiKey string
}

func NewSafeBrowsing() *SafeBrowsing {
	return &SafeBrowsing{apiKey: os.Getenv("GOOGLE_SAFE_BROWSING_API_KEY")}
}

func (s *SafeBrowsing) Enabled() bool {
	return s.apiKey != ""
}

func (s *SafeBrowsing) CheckURL(ctx context.Context, rawURL string) (*ports.URLVerdict, error) {
	verdict := &ports.URLVerdict{URL: rawURL, Malicious: false, Source: "safebrowsing-stub"}
	if !s.Enabled() {
		verdict.Reason = "api_key_not_configured"
		return verdict, nil
	}
	// Production: POST https://safebrowsing.googleapis.com/v4/threatMatches:find
	_ = ctx
	verdict.Reason = "lookup_not_implemented"
	return verdict, nil
}
