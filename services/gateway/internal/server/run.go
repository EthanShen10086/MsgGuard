package server

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/EthanShen10086/voxera-kit/auth"
	"github.com/EthanShen10086/voxera-kit/loadshed"
	"github.com/EthanShen10086/voxera-kit/loadshed/adaptive"
	kitmw "github.com/EthanShen10086/voxera-kit/middleware"
	"github.com/EthanShen10086/voxera-kit/observability/logger"
	"github.com/EthanShen10086/voxera-kit/observability/metrics"
	"github.com/EthanShen10086/voxera-kit/pii"
	piiregex "github.com/EthanShen10086/voxera-kit/pii/regex"
	"github.com/EthanShen10086/voxera-kit/security/headers"

	mgapp "github.com/EthanShen10086/msgguard/pkg/app"
	"github.com/EthanShen10086/msgguard/pkg/httpauth"
	gwmw "github.com/EthanShen10086/msgguard/services/gateway/internal/middleware"
	"github.com/EthanShen10086/msgguard/services/gateway/internal/handler"
)

// Run starts the gateway HTTP server with all middleware wired from Container.
func Run(c *mgapp.Container) error {
	ctx := context.Background()

	classifyHandler := handler.NewClassifyHandler(
		c.LLMClassifier, c.ThreatIntel, c.CircuitBreaker, c.Cache, c.Config, c.Log, c.Queue,
		c.QuotaStore, c.FlagStore,
	)
	feedbackHandler := handler.NewFeedbackHandler(c.Log, c.FeedbackStore, c.AuditWriter, c.Queue)
	analyticsHandler := handler.NewAnalyticsHandler(c.Log, c.AnalyticsStore)
	shadowHandler := handler.NewShadowHandler(c.Log, classifyHandler, c.FlagStore)
	rulesHandler := handler.NewRulesHandler(c.RuleStore, c.Authenticator, c.Authorizer)
	adminHandler := handler.NewAdminHandler(
		c.Log, c.AnalyticsStore, c.FeedbackStore, c.Authenticator, c.Authorizer,
		c.QuotaStore, c.FlagStore, shadowHandler,
	)
	privacyHandler := handler.NewPrivacyHandler(c.Log, c.AnalyticsStore)
	modelAdminHandler := handler.NewModelAdminHandler(
		c.Log, c.ModelRegistry, c.Authenticator, c.Authorizer,
	)
	entitlementsHandler := handler.NewEntitlementsHandler(c.Log, c.SubscriptionStore, c.AppStoreVerifier)
	authCfg := gwmw.LoadAuthProductionConfig()
	if c.Config.Security.AuthBootstrapEnabled {
		authCfg.BootstrapTokenEnabled = true
	}
	if !c.Config.Security.DeviceTokenEnabled {
		authCfg.DeviceTokenEnabled = false
	}
	if c.Config.Security.ModelDownloadAuth {
		authCfg.ModelDownloadAuth = true
	}
	deviceIssuer := &httpauth.MemoryDeviceIssuer{Auth: c.Authenticator}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/metrics", metrics.HTTPHandler())
	mux.HandleFunc("/api/v1/auth/token", gwmw.GateBootstrap(authCfg, adminHandler.IssueToken))
	mux.HandleFunc("/api/v1/auth/device", httpauth.IssueDeviceTokenHandler(deviceIssuer, authCfg.DeviceTokenEnabled))
	mux.HandleFunc("/api/v1/entitlements/verify", entitlementsHandler.Verify)
	mux.HandleFunc("/api/v1/entitlements/status", entitlementsHandler.Status)
	mux.HandleFunc("/api/v1/classify/defer", classifyHandler.Defer)
	mux.HandleFunc("/api/v1/classify", classifyHandler.Classify)
	mux.HandleFunc("/api/v1/classify/shadow", shadowHandler.Compare)
	mux.HandleFunc("/api/v1/admin/shadow/stats", adminHandler.ShadowStats)
	mux.HandleFunc("/metrics/shadow", shadowHandler.Prometheus)
	mux.HandleFunc("/api/v1/feedback", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			claims, err := authenticate(c, r)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ok, _ := c.Authorizer.Authorize(r.Context(), claims, "feedback", "read")
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			feedbackHandler.List(w, r)
			return
		}
		httpauth.RequireDeviceOrBearer(c.Authenticator, feedbackHandler.Create)(w, r)
	})
	mux.HandleFunc("/api/v1/analytics", httpauth.RequireDeviceOrBearer(c.Authenticator, analyticsHandler.Ingest))
	mux.HandleFunc("/api/v1/privacy/me", privacyHandler.DeleteMe)
	mux.HandleFunc("/api/v1/rules/latest", rulesHandler.Latest)
	mux.HandleFunc("/api/v1/rules/register", rulesHandler.Register)
	mux.HandleFunc("/api/v1/rules/", rulesHandler.ByVersion)
	mux.HandleFunc("/api/v1/admin/metrics/summary", adminHandler.MetricsSummary)
	mux.HandleFunc("/api/v1/admin/quota/whitelist", adminHandler.QuotaWhitelist)
	mux.HandleFunc("/api/v1/admin/flags", adminHandler.FeatureFlags)
	mux.HandleFunc("/api/v1/admin/models/promote", modelAdminHandler.Promote)
	mux.HandleFunc("/api/v1/admin/models/rollback", modelAdminHandler.Rollback)
	modelProxy := proxyTo(c.Config.Gateway.ModelAddr)
	mux.HandleFunc("/api/v1/models/register", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			claims, err := authenticate(c, r)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}
			ok, _ := c.Authorizer.Authorize(r.Context(), claims, "models", "write")
			if !ok {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
		}
		modelProxy.ServeHTTP(w, r)
	})
	mux.HandleFunc("/api/v1/models/", func(w http.ResponseWriter, r *http.Request) {
		if authCfg.ModelDownloadAuth && strings.Contains(r.URL.Path, "/download/") {
			httpauth.RequireModelRead(c.Authenticator, c.Authorizer, modelProxy.ServeHTTP)(w, r)
			return
		}
		modelProxy.ServeHTTP(w, r)
	})
	handler.RegisterHealthRoutes(mux)

	recorder := metrics.NewPrometheusRecorder()
	shedder := adaptive.New(loadshed.Config{MaxLoad: 0.9, Window: 10 * time.Second})
	redactor := piiregex.NewRedactor(pii.Config{Rules: piiregex.DefaultRules()})
	secCfg := headers.DefaultStrict()

	mws := []kitmw.Func{
		kitmw.Recovery(c.Log),
		kitmw.RequestID(),
		kitmw.Logging(c.Log),
		kitmw.Metrics(recorder),
		kitmw.SecurityHeaders(secCfg),
		kitmw.LoadShed(shedder),
		rateLimit(c),
		kitmw.Timeout(30 * time.Second),
		kitmw.PIIRedact(redactor),
		kitmw.HealthCheck(nil),
	}
	if c.Tracer != nil {
		mws = append([]kitmw.Func{kitmw.Tracing(c.Tracer)}, mws...)
	}
	mws = append(mws, httpauth.OIDCMiddleware("/api/v1/admin/"))
	if mtlsAdminRequired(c) {
		prefixes := []string{"/api/v1/admin/"}
		if h := mtlsClientHeader(c); h != "" {
			mws = append(mws, httpauth.RequireClientCertHeader(h, prefixes))
			c.Log.Info("mtls admin header enforcement enabled", logger.Field{Key: "header", Value: h})
		} else {
			mws = append(mws, httpauth.RequireClientCert(prefixes))
			c.Log.Info("mtls admin TLS cert enforcement enabled")
		}
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", c.Config.Gateway.Port),
		Handler:      kitmw.Chain(mux, mws...),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		c.Log.Info("gateway listening", logger.Field{Key: "port", Value: c.Config.Gateway.Port})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	c.Shutdown(shutdownCtx)
	return srv.Shutdown(shutdownCtx)
}

func authenticate(c *mgapp.Container, r *http.Request) (*auth.Claims, error) {
	h := r.Header.Get("Authorization")
	if h == "" {
		return nil, fmt.Errorf("missing auth")
	}
	return c.Authenticator.Authenticate(r.Context(), h)
}

func proxyTo(addr string) http.Handler {
	target, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(target)
}

func mtlsAdminRequired(c *mgapp.Container) bool {
	if os.Getenv("MTLS_ADMIN_REQUIRED") == "true" {
		return true
	}
	return c.Config.Security.MTLSAdminRequired
}

func mtlsClientHeader(c *mgapp.Container) string {
	if h := os.Getenv("MTLS_CLIENT_HEADER"); h != "" {
		return h
	}
	return c.Config.Security.MTLSClientHeader
}

func rateLimit(c *mgapp.Container) kitmw.Func {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			key := r.RemoteAddr
			if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
				key = strings.Split(xff, ",")[0]
			}
			ok, err := c.RateLimiter.Allow(r.Context(), key)
			if err != nil || !ok {
				http.Error(w, "rate limit exceeded", http.StatusTooManyRequests)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
