package app

import (
	"context"
	"os"
	"time"

	aiquotaMemory "github.com/EthanShen10086/voxera-kit/aiquota/memory"
	"github.com/EthanShen10086/voxera-kit/audit"
	auditMemory "github.com/EthanShen10086/voxera-kit/audit/memory"
	"github.com/EthanShen10086/voxera-kit/auth"
	"github.com/EthanShen10086/voxera-kit/circuitbreaker"
	cbAdapter "github.com/EthanShen10086/voxera-kit/circuitbreaker/memory"
	ffMemory "github.com/EthanShen10086/voxera-kit/featureflag/memory"
	"github.com/EthanShen10086/voxera-kit/llm"
	"github.com/EthanShen10086/voxera-kit/llm/deepseek"
	"github.com/EthanShen10086/voxera-kit/llm/qwen"
	"github.com/EthanShen10086/voxera-kit/observability/logger"
	"github.com/EthanShen10086/voxera-kit/observability/tracing"
	"github.com/EthanShen10086/voxera-kit/ratelimiter"
	rlMemory "github.com/EthanShen10086/voxera-kit/ratelimiter/memory"

	fsadapter "github.com/EthanShen10086/msgguard/pkg/adapters/filesystem"
	llmadapter "github.com/EthanShen10086/msgguard/pkg/adapters/llm"
	memadapters "github.com/EthanShen10086/msgguard/pkg/adapters/memory"
	mongoadapters "github.com/EthanShen10086/msgguard/pkg/adapters/mongodb"
	natsadapter "github.com/EthanShen10086/msgguard/pkg/adapters/nats"
	pgadapters "github.com/EthanShen10086/msgguard/pkg/adapters/postgres"
	redisadapters "github.com/EthanShen10086/msgguard/pkg/adapters/redis"
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
	AnalyticsStore  ports.AnalyticsStore
	LLMClassifier   ports.LLMClassifier
	Authenticator   auth.Authenticator
	Authorizer      auth.Authorizer
	RateLimiter     ratelimiter.RateLimiter
	CircuitBreaker  circuitbreaker.CircuitBreaker
	AuditWriter     audit.Writer
	QuotaStore      interface{}
	FlagStore       interface{}
	LLMRouter       *llm.Router
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
	c.FeedbackStore = wireFeedbackStore(cfg, driver, dsn, log)
	c.RuleStore = wireRuleStore(cfg, driver, dsn, log)
	c.AnalyticsStore = wireAnalyticsStore(cfg, driver, dsn, log)
	c.Cache = wireCache(cfg, log)
	c.Queue = wireQueue(cfg, log)
	c.ModelRegistry = wireModelRegistry(cfg, log)

	c.LLMRouter = llm.NewRouter()
	if os.Getenv("QWEN_API_KEY") != "" {
		c.LLMRouter.Register(qwen.New(llm.Config{APIKey: os.Getenv("QWEN_API_KEY"), Model: "qwen-turbo"}), 1)
	}
	if os.Getenv("DEEPSEEK_API_KEY") != "" {
		c.LLMRouter.Register(deepseek.New(llm.Config{APIKey: os.Getenv("DEEPSEEK_API_KEY"), Model: "deepseek-chat"}), 2)
	}
	c.LLMClassifier = llmadapter.NewClassifier(c.LLMRouter)

	secret := envOr("AUTH_SECRET", "msgguard-dev-secret")
	memAuth := memadapters.NewAuth(secret)
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
	c.QuotaStore = aiquotaMemory.NewStore()
	c.FlagStore = ffMemory.NewAdapter()

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

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
