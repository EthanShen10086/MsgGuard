package redis

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/EthanShen10086/voxera-kit/auth"
)

// RevocationStore persists revoked token hashes in Redis.
type RevocationStore struct {
	client *goredis.Client
	prefix string
	ttl    time.Duration
}

func NewRevocationStore(url string) (*RevocationStore, error) {
	opt, err := goredis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	c := goredis.NewClient(opt)
	if err := c.Ping(context.Background()).Err(); err != nil {
		return nil, err
	}
	return &RevocationStore{client: c, prefix: "msgguard:revoked:", ttl: 30 * 24 * time.Hour}, nil
}

func (r *RevocationStore) Revoke(ctx context.Context, token string) error {
	token = strings.TrimPrefix(token, "Bearer ")
	h := hashToken(token)
	return r.client.Set(ctx, r.prefix+h, "1", r.ttl).Err()
}

func (r *RevocationStore) IsRevoked(ctx context.Context, token string) (bool, error) {
	token = strings.TrimPrefix(token, "Bearer ")
	n, err := r.client.Exists(ctx, r.prefix+hashToken(token)).Result()
	return n > 0, err
}

func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// RevokingAuth wraps Authenticator and checks Redis revocation before/after delegate.
type RevokingAuth struct {
	inner auth.Authenticator
	rev   *RevocationStore
}

func NewRevokingAuth(inner auth.Authenticator, rev *RevocationStore) *RevokingAuth {
	return &RevokingAuth{inner: inner, rev: rev}
}

func (r *RevokingAuth) Authenticate(ctx context.Context, token string) (*auth.Claims, error) {
	if revoked, err := r.rev.IsRevoked(ctx, token); err == nil && revoked {
		return nil, errors.New("token revoked")
	}
	return r.inner.Authenticate(ctx, token)
}

func (r *RevokingAuth) GenerateToken(ctx context.Context, claims *auth.Claims) (*auth.TokenPair, error) {
	return r.inner.GenerateToken(ctx, claims)
}

func (r *RevokingAuth) RefreshToken(ctx context.Context, refreshToken string) (*auth.TokenPair, error) {
	return r.inner.RefreshToken(ctx, refreshToken)
}

func (r *RevokingAuth) RevokeToken(ctx context.Context, token string) error {
	if err := r.rev.Revoke(ctx, token); err != nil {
		return err
	}
	return r.inner.RevokeToken(ctx, token)
}

// ErrInvalidToken alias when voxera-kit doesn't export - use string error from inner
var _ auth.Authenticator = (*RevokingAuth)(nil)
