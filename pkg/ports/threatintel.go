package ports

import "context"

// URLVerdict is the result of an L0 threat-intel URL lookup.
type URLVerdict struct {
	URL       string `json:"url"`
	Malicious bool   `json:"malicious"`
	Source    string `json:"source"`
	Reason    string `json:"reason,omitempty"`
}

// ThreatIntel checks URLs against external reputation feeds (Safe Browsing, etc.).
type ThreatIntel interface {
	CheckURL(ctx context.Context, rawURL string) (*URLVerdict, error)
}
