package httpauth

import (
	"os"
	"testing"
)

func TestParseAdminEmails(t *testing.T) {
	m := parseAdminEmails("A@x.com, b@y.com ,,")
	if len(m) != 2 {
		t.Fatalf("got %d", len(m))
	}
	if _, ok := m["a@x.com"]; !ok {
		t.Fatal("missing a@x.com")
	}
}

func TestOIDCProviderIsAdminByEmailList(t *testing.T) {
	p := &OIDCProvider{
		admins: map[string]struct{}{"ops@corp.com": {}},
	}
	if !p.isAdmin("ops@corp.com") {
		t.Fatal("expected allow")
	}
	if p.isAdmin("other@corp.com") {
		t.Fatal("expected deny")
	}
}

func TestOIDCProviderIsAdminByDomain(t *testing.T) {
	p := &OIDCProvider{adminDom: "msgguard.app"}
	if !p.isAdmin("user@msgguard.app") {
		t.Fatal("domain match")
	}
	if p.isAdmin("user@evil.com") {
		t.Fatal("should deny")
	}
}

func TestOIDCProviderProductionDenyWithoutAllowlist(t *testing.T) {
	os.Setenv("MSGGUARD_ENV", "production")
	defer os.Unsetenv("MSGGUARD_ENV")
	p := &OIDCProvider{}
	if p.isAdmin("anyone@gmail.com") {
		t.Fatal("production without allowlist must deny")
	}
}

func TestOIDCConfiguredFalseWhenMissing(t *testing.T) {
	os.Unsetenv("OIDC_ISSUER")
	os.Unsetenv("OIDC_CLIENT_ID")
	os.Unsetenv("OIDC_CLIENT_SECRET")
	if OIDCConfigured() {
		t.Fatal("should be false")
	}
}
