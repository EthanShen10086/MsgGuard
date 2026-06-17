// Package appstore implements App Store Server API integration.
package appstore

import (
	"context"
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Client calls Apple App Store Server API when credentials are configured.
type Client struct {
	issuerID   string
	keyID      string
	bundleID   string
	privateKey *ecdsa.PrivateKey
	http       *http.Client
	baseURL    string
}

// NewClientFromEnv reads APPLE_ISSUER_ID, APPLE_KEY_ID, APPLE_BUNDLE_ID, APPLE_PRIVATE_KEY.
func NewClientFromEnv() (*Client, error) {
	issuer := os.Getenv("APPLE_ISSUER_ID")
	keyID := os.Getenv("APPLE_KEY_ID")
	bundle := os.Getenv("APPLE_BUNDLE_ID")
	pemKey := os.Getenv("APPLE_PRIVATE_KEY")
	if issuer == "" || keyID == "" || pemKey == "" {
		return nil, errors.New("apple credentials not configured")
	}
	if bundle == "" {
		bundle = "com.ethanshen.msgguard"
	}
	block, _ := pem.Decode([]byte(strings.ReplaceAll(pemKey, `\n`, "\n")))
	if block == nil {
		return nil, errors.New("invalid APPLE_PRIVATE_KEY PEM")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	ec, ok := key.(*ecdsa.PrivateKey)
	if !ok {
		return nil, errors.New("APPLE_PRIVATE_KEY must be ECDSA")
	}
	return &Client{
		issuerID: issuer, keyID: keyID, bundleID: bundle, privateKey: ec,
		http: &http.Client{Timeout: 15 * time.Second},
		baseURL: "https://api.storekit.itunes.apple.com",
	}, nil
}

func (c *Client) Enabled() bool { return c != nil }

func (c *Client) bearerToken() (string, error) {
	claims := jwt.MapClaims{
		"iss": c.issuerID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(5 * time.Minute).Unix(),
		"aud": "appstoreconnect-v1",
		"bid": c.bundleID,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	t.Header["kid"] = c.keyID
	return t.SignedString(c.privateKey)
}

// GetTransactionInfo fetches signed transaction info from Apple (production API).
func (c *Client) GetTransactionInfo(ctx context.Context, transactionID string) ([]byte, error) {
	if c == nil {
		return nil, errors.New("apple client not configured")
	}
	tok, err := c.bearerToken()
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/inApps/v1/transactions/%s", c.baseURL, transactionID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+tok)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("apple API status %d", resp.StatusCode)
	}
	var out struct {
		SignedTransactionInfo string `json:"signedTransactionInfo"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return []byte(out.SignedTransactionInfo), nil
}
