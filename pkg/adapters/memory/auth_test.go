package memory

import (
	"context"
	"testing"

	"github.com/EthanShen10086/voxera-kit/auth"
)

func TestValidateSecretProduction(t *testing.T) {
	if err := ValidateSecret("", "production", false); err == nil {
		t.Fatal("expected error for empty secret in production")
	}
	if err := ValidateSecret("short", "production", false); err == nil {
		t.Fatal("expected error for short secret in production")
	}
	if err := ValidateSecret("this-is-a-secure-production-secret-key-32", "production", false); err != nil {
		t.Fatalf("unexpected: %v", err)
	}
}

func TestGenerateAndAuthenticateAdmin(t *testing.T) {
	a := NewAuthWithOptions("test-secret-key-for-unit-tests-only", true)
	ctx := context.Background()
	pair, err := a.GenerateToken(ctx, &auth.Claims{UserID: "u1", Roles: []string{"admin"}})
	if err != nil {
		t.Fatal(err)
	}
	claims, err := a.Authenticate(ctx, pair.AccessToken)
	if err != nil {
		t.Fatal(err)
	}
	if claims.UserID != "u1" {
		t.Fatalf("uid %q", claims.UserID)
	}
	if !a.HasRole(ctx, claims, "admin") {
		t.Fatal("expected admin role")
	}
}

func TestRBACAdminCanWriteModels(t *testing.T) {
	a := NewAuthWithOptions("test-secret-key-for-unit-tests-only", true)
	ctx := context.Background()
	pair, _ := a.GenerateToken(ctx, &auth.Claims{UserID: "admin1", Roles: []string{"admin"}})
	claims, _ := a.Authenticate(ctx, pair.AccessToken)
	ok, err := a.Authorize(ctx, claims, "models", "write")
	if err != nil || !ok {
		t.Fatalf("admin should write models: ok=%v err=%v", ok, err)
	}
}

func TestRBACDeviceCannotWriteModels(t *testing.T) {
	a := NewAuthWithOptions("test-secret-key-for-unit-tests-only", true)
	ctx := context.Background()
	pair, _ := a.GenerateToken(ctx, &auth.Claims{UserID: "dev1", Roles: []string{"device"}})
	claims, _ := a.Authenticate(ctx, pair.AccessToken)
	ok, _ := a.Authorize(ctx, claims, "models", "write")
	if ok {
		t.Fatal("device must not write models")
	}
}

func TestRBACMLEngineerFeedbackRead(t *testing.T) {
	a := NewAuthWithOptions("test-secret-key-for-unit-tests-only", true)
	ctx := context.Background()
	pair, _ := a.GenerateToken(ctx, &auth.Claims{UserID: "ml1", Roles: []string{"ml_engineer"}})
	claims, _ := a.Authenticate(ctx, pair.AccessToken)
	ok, _ := a.Authorize(ctx, claims, "feedback", "read")
	if !ok {
		t.Fatal("ml_engineer should read feedback")
	}
}

func TestRevokeToken(t *testing.T) {
	a := NewAuthWithOptions("test-secret-key-for-unit-tests-only", true)
	ctx := context.Background()
	pair, _ := a.GenerateToken(ctx, &auth.Claims{UserID: "u1", Roles: []string{"admin"}})
	if err := a.RevokeToken(ctx, pair.AccessToken); err != nil {
		t.Fatal(err)
	}
	if _, err := a.Authenticate(ctx, pair.AccessToken); err == nil {
		t.Fatal("revoked token should fail")
	}
}

func TestInvalidTokenRejected(t *testing.T) {
	a := NewAuthWithOptions("test-secret-key-for-unit-tests-only", true)
	if _, err := a.Authenticate(context.Background(), "not.a.valid.token"); err == nil {
		t.Fatal("expected invalid token error")
	}
}
