package middleware

import (
	"net/http"
	"os"
	"strings"
)

// AuthProductionConfig gates bootstrap endpoints in production.
type AuthProductionConfig struct {
	BootstrapTokenEnabled bool
	DeviceTokenEnabled    bool
	ModelDownloadAuth     bool
}

// LoadAuthProductionConfig reads auth policy from environment.
func LoadAuthProductionConfig() AuthProductionConfig {
	env := strings.ToLower(os.Getenv("MSGGUARD_ENV"))
	isProd := env == "production" || env == "prod"
	bootstrap := os.Getenv("AUTH_BOOTSTRAP_ENABLED")
	device := os.Getenv("AUTH_DEVICE_TOKEN_ENABLED")
	modelAuth := os.Getenv("MODEL_DOWNLOAD_AUTH_REQUIRED")
	return AuthProductionConfig{
		BootstrapTokenEnabled: envBool(bootstrap, !isProd),
		DeviceTokenEnabled:    envBool(device, true),
		ModelDownloadAuth:     envBool(modelAuth, isProd),
	}
}

func envBool(v string, defaultVal bool) bool {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "1", "true", "yes", "on":
		return true
	case "0", "false", "no", "off":
		return false
	default:
		return defaultVal
	}
}

// GateBootstrap wraps admin token issuance — disabled in production unless AUTH_BOOTSTRAP_ENABLED=true.
func GateBootstrap(cfg AuthProductionConfig, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !cfg.BootstrapTokenEnabled {
			http.Error(w, "bootstrap token issuance disabled", http.StatusForbidden)
			return
		}
		next(w, r)
	}
}
