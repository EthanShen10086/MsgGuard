package app

import (
	"context"
	"os"
	"time"

	"github.com/EthanShen10086/voxera-kit/aiquota"
	aiquotaMemory "github.com/EthanShen10086/voxera-kit/aiquota/memory"
	"github.com/EthanShen10086/voxera-kit/audit"
	auditMemory "github.com/EthanShen10086/voxera-kit/audit/memory"
	"github.com/EthanShen10086/voxera-kit/auth"
	"github.com/EthanShen10086/voxera-kit/circuitbreaker"
	cbAdapter "github.com/EthanShen10086/voxera-kit/circuitbreaker/memory"
	"github.com/EthanShen10086/voxera-kit/featureflag"
	ffMemory "github.com/EthanShen10086/voxera-kit/featureflag/memory"
	"github.com/EthanShen10086/voxera-kit/llm"
	"github.com/EthanShen10086/voxera-kit/llm/deepseek"
	"github.com/EthanShen10086/voxera-kit/llm/qwen"
	"github.com/EthanShen10086/voxera-kit/observability/logger"
	"github.com/EthanShen10086/voxera-kit/observability/tracing"
	"github.com/EthanShen10086/voxera-kit/ratelimiter"
	rlMemory "github.com/EthanShen10086/voxera-kit/ratelimiter/memory"

	appstoreadapter "github.com/EthanShen10086/msgguard/pkg/adapters/appstore"
	fsadapter "github.com/EthanShen10086/msgguard/pkg/adapters/filesystem"
	llmadapter "github.com/EthanShen10086/msgguard/pkg/adapters/llm"
	memadapters "github.com/EthanShen10086/msgguard/pkg/adapters/memory"
	mongoadapters "github.com/EthanShen10086/msgguard/pkg/adapters/mongodb"
	natsadapter "github.com/EthanShen10086/msgguard/pkg/adapters/nats"
	pgadapters "github.com/EthanShen10086/msgguard/pkg/adapters/postgres"
	redisadapters "github.com/EthanShen10086/msgguard/pkg/adapters/redis"
	tiadapter "github.com/EthanShen10086/msgguard/pkg/adapters/threatintel"
	"github.com/EthanShen10086/msgguard/pkg/config"
	"github.com/EthanShen10086/msgguard/pkg/ports"
)

// Container holds all wired dependencies (IoC root).
type Container struct {
	Config          *config.Config
	Log             logger.Logger
	Tracer          tracing.Tracer
	FeedbackStore   ports.FeedbackStore
	RuleStore       ports.RuleStore
	Cache           ports.Cache
	Queue           ports.Queue
	ModelRegistry   ports.ModelRegistry
	ThreatIntel     ports.ThreatIntel
	AnalyticsStore  ports.AnalyticsStore
	LLMClassifier   ports.LLMClassifier
	Authenticator   auth.Authenticator
	Authorizer      auth.Authorizer
	RateLimiter     ratelimiter.RateLimiter
	CircuitBreaker  circuitbreaker.CircuitBreaker
	AuditWriter     audit.Writer
	QuotaStore        aiquota.Manager
	FlagStore         featureflag.Store
	SubscriptionStore ports.SubscriptionStore
	AppStoreVerifier  *appstoreadapter.Verifier
	LLMRouter         *llm.Router
}

func NewContainer(cfg *config.Config) (*Container, error) {
	ctx := context.Background()
	log, err := logger.NewZapLogger()
	if err != nil {
		return nil, err
	}

	c := &Container{Config: cfg, Log: log}

	tracer, _ := tracing.NewOTelTracer(ctx, tracing.OTelConfig{
		ServiceName: "msgguard-gateway",
		Endpoint:    cfg.Observability.OTelEndpoint,
		SampleRate:  cfg.Observability.TraceSampleRate,
		Insecure:    true,
	})
	c.Tracer = tracer

	dsn := envOr("DATABASE_DSN", cfg.Database.DSN)
	driver := envOr("DATABASE_DRIVER", cfg.Database.Driver)
	if dsn != "" && (driver == "postgres" || driver == "") {
		if err := pgadapters.Migrate(dsn); err != nil {
			log.Info("postgres migration warning", logger.Field{Key: "error", Value: err.Error()})
		} else {
			log.Info("postgres migrations applied")
		}
	}
	c.FeedbackStore = wireFeedbackStore(cfg, driver, dsn, log)
	c.RuleStore = wireRuleStore(cfg, driver, dsn, log)
	c.AnalyticsStore = wireAnalyticsStore(cfg, driver, dsn, log)
	c.Cache = wireCache(cfg, log)
	c.Queue = wireQueue(cfg, log)
	c.ModelRegistry = wireModelRegistry(cfg, log)
	c.ThreatIntel = wireThreatIntel(log)
	c.SubscriptionStore = wireSubscriptionStore(driver, dsn, log)
	c.AppStoreVerifier = appstoreadapter.NewVerifier()

	c.LLMRouter = llm.NewRouter()
	if os.Getenv("QWEN_API_KEY") != "" {
		c.LLMRouter.Register(qwen.New(llm.Config{APIKey: os.Getenv("QWEN_API_KEY"), Model: "qwen-turbo"}), 1)
	}
	if os.Getenv("DEEPSEEK_API_KEY") != "" {
		c.LLMRouter.Register(deepseek.New(llm.Config{APIKey: os.Getenv("DEEPSEEK_API_KEY"), Model: "deepseek-chat"}), 2)
	}
	c.LLMClassifier = llmadapter.NewClassifier(c.LLMRouter)

	secret := envOr("AUTH_SECRET", "")
	allowDev := os.Getenv("MSGGUARD_ENV") != "production" && os.Getenv("MSGGUARD_ENV") != "prod"
	if err := memadapters.ValidateSecret(secret, os.Getenv("MSGGUARD_ENV"), allowDev); err != nil && os.Getenv("CI") != "true" {
		log.Info("auth secret warning", logger.Field{Key: "error", Value: err.Error()})
	}
	if secret == "" && allowDev {
		secret = "msgguard-dev-secret"
	}
	memAuth := memadapters.NewAuthWithOptions(secret, allowDev)
	c.Authenticator = memAuth
	c.Authorizer = memAuth

	c.RateLimiter = rlMemory.New(ratelimiter.Config{
		Enabled: true,
		Rate:    cfg.Gateway.RateLimitRPS,
		Burst:   int(cfg.Gateway.RateLimitRPS * 2),
	})
	if c.RateLimiter == nil || cfg.Gateway.RateLimitRPS <= 0 {
		c.RateLimiter = rlMemory.New(ratelimiter.Config{Enabled: true, Rate: 100, Burst: 200})
	}

	c.CircuitBreaker = cbAdapter.New(circuitbreaker.Config{
		MaxFailures: 3, Timeout: 30 * time.Second, HalfOpenMaxCalls: 2,
	})

	c.AuditWriter = auditMemory.NewAdapter()
	c.QuotaStore, c.FlagStore = wireQuotaAndFlags(driver, dsn, log)
	_ = c.FlagStore.SetFlag(ctx, featureflag.Flag{Key: "cloud_llm", Enabled: cfg.Features.CloudLLM, Percentage: 100})
	_ = c.FlagStore.SetFlag(ctx, featureflag.Flag{Key: "shadow_mode", Enabled: false, Percentage: 0})
	_ = c.FlagStore.SetFlag(ctx, featureflag.Flag{Key: "model_canary", Enabled: false, Percentage: 0})

	return c, nil
}

func (c *Container) Shutdown(ctx context.Context) {
	if c.Tracer != nil {
		_ = c.Tracer.Shutdown(ctx)
	}
}

func wireFeedbackStore(cfg *config.Config, driver, dsn string, log logger.Logger) ports.FeedbackStore {
	if dsn != "" && driver == "mongodb" {
		if s, err := mongoadapters.NewFeedbackStore(dsn); err == nil {
			log.Info("mongodb feedback store connected")
			return s
		}
	}
	if dsn != "" && (driver == "postgres" || driver == "") {
		if s, err := pgadapters.NewFeedbackStore(dsn); err == nil {
			log.Info("postgres feedback store connected")
			return s
		}
	}
	return memadapters.NewFeedbackStore()
}

func wireRuleStore(cfg *config.Config, driver, dsn string, log logger.Logger) ports.RuleStore {
	if dsn != "" && driver == "mongodb" {
		if s, err := mongoadapters.NewRuleStore(dsn); err == nil {
			log.Info("mongodb rule store connected")
			return s
		}
	}
	if dsn != "" && (driver == "postgres" || driver == "") {
		if s, err := pgadapters.NewRuleStore(dsn); err == nil {
			log.Info("postgres rule store connected")
			return s
		}
	}
	return memadapters.NewRuleStore()
}

func wireAnalyticsStore(cfg *config.Config, driver, dsn string, log logger.Logger) ports.AnalyticsStore {
	if dsn != "" && driver == "mongodb" {
		if s, err := mongoadapters.NewAnalyticsStore(dsn); err == nil {
			log.Info("mongodb analytics store connected")
			return s
		}
	}
	if dsn != "" && (driver == "postgres" || driver == "") {
		if s, err := pgadapters.NewAnalyticsStore(dsn); err == nil {
			log.Info("postgres analytics store connected")
			return s
		}
	}
	return memadapters.NewAnalyticsStore()
}

func wireCache(cfg *config.Config, log logger.Logger) ports.Cache {
	if cfg.Cache.RedisEnabled {
		if url := envOr("REDIS_URL", ""); url != "" {
			if rc, err := redisadapters.NewCache(url); err == nil {
				log.Info("redis cache connected")
				return rc
			}
		}
	}
	return memadapters.NewCache()
}

func wireQueue(cfg *config.Config, log logger.Logger) ports.Queue {
	if cfg.Queue.Driver == "nats" {
		url := envOr("NATS_URL", cfg.Queue.URL)
		if url == "" {
			url = "nats://localhost:4222"
		}
		if q, err := natsadapter.NewQueue(url); err == nil {
			log.Info("nats queue connected")
			return q
		}
		log.Info("nats fallback to noop")
	}
	return memadapters.NewQueue()
}

func wireSubscriptionStore(driver, dsn string, log logger.Logger) ports.SubscriptionStore {
	if dsn != "" && driver == "mongodb" {
		if s, err := mongoadapters.NewSubscriptionStore(dsn); err == nil {
			log.Info("mongodb subscription store connected")
			return s
		}
	}
	if dsn != "" && (driver == "postgres" || driver == "") {
		if s, err := pgadapters.NewSubscriptionStore(dsn); err == nil {
			log.Info("postgres subscription store connected")
			return s
		}
	}
	return memadapters.NewSubscriptionStore()
}

func wireQuotaAndFlags(driver, dsn string, log logger.Logger) (aiquota.Manager, featureflag.Store) {
	if dsn != "" && (driver == "postgres" || driver == "") {
		if q, err := pgadapters.NewQuotaManager(dsn); err == nil {
			log.Info("postgres quota whitelist connected")
			if f, err := pgadapters.NewFlagStore(dsn); err == nil {
				log.Info("postgres feature flags connected")
				return q, f
			}
			return q, ffMemory.NewAdapter()
		}
	}
	return aiquotaMemory.NewStore(), ffMemory.NewAdapter()
}

func wireModelRegistry(cfg *config.Config, log logger.Logger) ports.ModelRegistry {
	path := cfg.ModelStorage.Path
	if path == "" {
		path = "./deploy/models"
	}
	if reg, err := fsadapter.NewModelRegistry(path); err == nil {
		return reg
	}
	return memadapters.NewModelRegistry()
}

func wireThreatIntel(log logger.Logger) ports.ThreatIntel {
	sb := tiadapter.NewSafeBrowsing()
	if sb.Enabled() {
		log.Info("safe browsing threat intel enabled (stub)")
	}
	return sb
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
