// Package main is the MsgGuard API Gateway entry point.
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	aiquotaMemory "github.com/EthanShen10086/voxera-kit/aiquota/memory"
	ffMemory "github.com/EthanShen10086/voxera-kit/featureflag/memory"
	"github.com/EthanShen10086/voxera-kit/llm"
	"github.com/EthanShen10086/voxera-kit/llm/deepseek"
	"github.com/EthanShen10086/voxera-kit/llm/qwen"
	"github.com/EthanShen10086/voxera-kit/loadshed"
	"github.com/EthanShen10086/voxera-kit/loadshed/adaptive"
	kitmw "github.com/EthanShen10086/voxera-kit/middleware"
	"github.com/EthanShen10086/voxera-kit/observability/logger"
	"github.com/EthanShen10086/voxera-kit/observability/metrics"
	"github.com/EthanShen10086/voxera-kit/observability/tracing"
	"github.com/EthanShen10086/voxera-kit/pii"
	piiregex "github.com/EthanShen10086/voxera-kit/pii/regex"
	"github.com/EthanShen10086/voxera-kit/security/headers"

	memadapters "github.com/EthanShen10086/msgguard/pkg/adapters/memory"
	pgadapters "github.com/EthanShen10086/msgguard/pkg/adapters/postgres"
	redisadapters "github.com/EthanShen10086/msgguard/pkg/adapters/redis"
	"github.com/EthanShen10086/msgguard/pkg/config"
	"github.com/EthanShen10086/msgguard/pkg/ports"
	"github.com/EthanShen10086/msgguard/services/gateway/internal/handler"
)

func main() {
	ctx := context.Background()
	cfgPath := envOr("CONFIG_PATH", "../../deploy/config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic(err)
	}

	log, err := logger.NewZapLogger()
	if err != nil {
		panic(err)
	}

	tracer, _ := tracing.NewOTelTracer(ctx, tracing.OTelConfig{
		ServiceName: "msgguard-gateway",
		Endpoint:    cfg.Observability.OTelEndpoint,
		SampleRate:  cfg.Observability.TraceSampleRate,
		Insecure:    true,
	})
	if tracer != nil {
		defer tracer.Shutdown(ctx)
	}

	recorder := metrics.NewPrometheusRecorder()
	shedder := adaptive.New(loadshed.Config{MaxLoad: 0.9, Window: 10 * time.Second})
	redactor := piiregex.NewRedactor(pii.Config{Rules: piiregex.DefaultRules()})
	secCfg := headers.DefaultStrict()

	flagStore := ffMemory.NewAdapter()
	quotaStore := aiquotaMemory.NewStore()
	llmRouter := llm.NewRouter()
	if os.Getenv("QWEN_API_KEY") != "" {
		llmRouter.Register(qwen.New(llm.Config{APIKey: os.Getenv("QWEN_API_KEY"), Model: "qwen-turbo"}), 1)
	}
	if os.Getenv("DEEPSEEK_API_KEY") != "" {
		llmRouter.Register(deepseek.New(llm.Config{APIKey: os.Getenv("DEEPSEEK_API_KEY"), Model: "deepseek-chat"}), 2)
	}

	feedbackStore := wireFeedbackStore(cfg, log)
	var cache ports.Cache = memadapters.NewCache()
	if cfg.Cache.RedisEnabled {
		if redisURL := envOr("REDIS_URL", ""); redisURL != "" {
			if rc, err := redisadapters.NewCache(redisURL); err == nil {
				cache = rc
				log.Info("redis cache connected")
			} else {
				log.Info("redis fallback to memory", logger.Field{Key: "err", Value: err.Error()})
			}
		}
	}

	classifyHandler := handler.NewClassifyHandler(llmRouter, quotaStore, flagStore, cfg, log, cache)
	feedbackHandler := handler.NewFeedbackHandler(log, feedbackStore)
	analyticsHandler := handler.NewAnalyticsHandler(log)
	shadowHandler := handler.NewShadowHandler(log, classifyHandler)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.Handle("/metrics", metrics.HTTPHandler())
	mux.HandleFunc("/api/v1/classify/defer", classifyHandler.Defer)
	mux.HandleFunc("/api/v1/classify", classifyHandler.Classify)
	mux.HandleFunc("/api/v1/classify/shadow", shadowHandler.Compare)
	mux.HandleFunc("/api/v1/feedback", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			feedbackHandler.List(w, r)
			return
		}
		feedbackHandler.Create(w, r)
	})
	mux.HandleFunc("/api/v1/analytics", analyticsHandler.Ingest)
	mux.Handle("/api/v1/rules/", proxyTo(cfg.Gateway.RulesAddr))
	mux.Handle("/api/v1/models/", proxyTo(cfg.Gateway.ModelAddr))
	handler.RegisterHealthRoutes(mux)

	mws := []kitmw.Func{
		kitmw.Recovery(log),
		kitmw.RequestID(),
		kitmw.Logging(log),
		kitmw.Metrics(recorder),
		kitmw.SecurityHeaders(secCfg),
		kitmw.LoadShed(shedder),
		kitmw.Timeout(30 * time.Second),
		kitmw.PIIRedact(redactor),
		kitmw.HealthCheck(nil),
	}
	if tracer != nil {
		mws = append([]kitmw.Func{kitmw.Tracing(tracer)}, mws...)
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Gateway.Port),
		Handler:      kitmw.Chain(mux, mws...),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		log.Info("gateway listening", logger.Field{Key: "port", Value: cfg.Gateway.Port})
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func wireFeedbackStore(cfg *config.Config, log logger.Logger) ports.FeedbackStore {
	driver := cfg.Database.Driver
	if dsn := envOr("DATABASE_DSN", cfg.Database.DSN); dsn != "" && (driver == "postgres" || driver == "") {
		store, err := pgadapters.NewFeedbackStore(dsn)
		if err == nil {
			log.Info("postgres feedback store connected")
			return store
		}
		log.Info("postgres fallback to memory", logger.Field{Key: "err", Value: err.Error()})
	}
	return memadapters.NewFeedbackStore()
}

func proxyTo(addr string) http.Handler {
	target, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}
	return httputil.NewSingleHostReverseProxy(target)
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
