// Package appstore provides a stub App Store Server API transaction verifier.
package appstore

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

// VerificationResult holds parsed JWS metadata from a signed transaction.
type VerificationResult struct {
	ProductID    string
	TransactionID string
	Valid        bool
	ExpiresAt    *time.Time
}

// Verifier validates App Store signed transactions (JWS parse + optional Apple API).
type Verifier struct {
	client *Client
}

func NewVerifier() *Verifier {
	c, _ := NewClientFromEnv()
	return &Verifier{client: c}
}

// VerifySignedTransaction checks JWS structure and decodes the payload header fields.
func (v *Verifier) VerifySignedTransaction(signedTransaction, productID string) (*VerificationResult, error) {
	if strings.TrimSpace(signedTransaction) == "" {
		return nil, errors.New("signed_transaction required")
	}
	if strings.TrimSpace(productID) == "" {
		return nil, errors.New("product_id required")
	}

	parts := strings.Split(signedTransaction, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWS: expected 3 segments")
	}
	for _, part := range parts {
		if _, err := base64.RawURLEncoding.DecodeString(part); err != nil {
			if _, err2 := base64.URLEncoding.DecodeString(part); err2 != nil {
				return nil, errors.New("invalid JWS segment encoding")
			}
		}
	}

	var header struct {
		Alg string `json:"alg"`
	}
	headerRaw, err := decodeJWSSegment(parts[0])
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(headerRaw, &header); err != nil {
		return nil, errors.New("invalid JWS header")
	}
	if header.Alg == "" {
		return nil, errors.New("JWS header missing alg")
	}

	payloadRaw, err := decodeJWSSegment(parts[1])
	if err != nil {
		return nil, err
	}
	var payload struct {
		ProductID     string `json:"productId"`
		TransactionID string `json:"transactionId"`
		ExpiresDate   int64  `json:"expiresDate"`
	}
	_ = json.Unmarshal(payloadRaw, &payload)

	result := &VerificationResult{
		ProductID:     productID,
		TransactionID: payload.TransactionID,
		Valid:         true,
	}
	if payload.ProductID != "" {
		result.ProductID = payload.ProductID
	}
	if payload.ExpiresDate > 0 {
		t := time.UnixMilli(payload.ExpiresDate)
		result.ExpiresAt = &t
	}
	return result, nil
}

func decodeJWSSegment(seg string) ([]byte, error) {
	if b, err := base64.RawURLEncoding.DecodeString(seg); err == nil {
		return b, nil
	}
	return base64.URLEncoding.DecodeString(seg)
}
