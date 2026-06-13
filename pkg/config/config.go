package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Features      Features      `yaml:"features"`
	Security      Security      `yaml:"security"`
	Gateway       Gateway       `yaml:"gateway"`
	Database      Database      `yaml:"database"`
	Cache         Cache         `yaml:"cache"`
	Queue         Queue         `yaml:"queue"`
	ModelStorage  ModelStorage  `yaml:"model"`
	Observability Observability `yaml:"observability"`
	Scheduler     Scheduler     `yaml:"scheduler"`
	LLM           LLM           `yaml:"llm"`
}

type Database struct {
	Driver string `yaml:"driver"`
	DSN    string `yaml:"dsn"`
}

type Queue struct {
	Driver string `yaml:"driver"`
	URL    string `yaml:"url"`
}

type ModelStorage struct {
	Driver string `yaml:"driver"`
	Path   string `yaml:"path"`
}

type Features struct {
	CloudLLM           bool `yaml:"cloud_llm"`
	CrashReportUpload  bool `yaml:"crash_report_upload"`
	IPWhitelistEnabled bool `yaml:"ip_whitelist_enabled"`
}

type Security struct {
	MTLSAdminRequired bool   `yaml:"mtls_admin_required"`
	MTLSClientHeader  string `yaml:"mtls_client_header"`
}

type Gateway struct {
	Port           int      `yaml:"port"`
	RateLimitRPS   float64  `yaml:"rate_limit_rps"`
	IPWhitelist    []string `yaml:"ip_whitelist"`
	RulesAddr      string   `yaml:"rules_addr"`
	ClassifyAddr   string   `yaml:"classify_addr"`
	ModelAddr      string   `yaml:"model_addr"`
	FeedbackAddr   string   `yaml:"feedback_addr"`
}

type Cache struct {
	RedisEnabled bool `yaml:"redis_enabled"`
	RulesTTL     int  `yaml:"rules_ttl"`
}

type Observability struct {
	TraceSampleRate float64 `yaml:"trace_sample_rate"`
	MetricsEnabled  bool    `yaml:"metrics_enabled"`
	OTelEndpoint    string  `yaml:"otel_endpoint"`
}

type Scheduler struct {
	RuleSyncCron string `yaml:"rule_sync_cron"`
}

type LLM struct {
	PrimaryProvider  string `yaml:"primary_provider"`
	FallbackProvider string `yaml:"fallback_provider"`
	DailyQuotaFree   int    `yaml:"daily_quota_free"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Default(), nil
	}
	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

func Default() *Config {
	return &Config{
		Features: Features{IPWhitelistEnabled: false},
		Gateway: Gateway{
			Port:         8080,
			RateLimitRPS: 100,
			RulesAddr:    "http://localhost:8081",
			ClassifyAddr: "http://localhost:8082",
			ModelAddr:    "http://localhost:8083",
			FeedbackAddr: "http://localhost:8084",
		},
		Cache:         Cache{RedisEnabled: true, RulesTTL: 3600},
		Database:      Database{Driver: "memory"},
		Queue:         Queue{Driver: "noop"},
		ModelStorage:  ModelStorage{Driver: "filesystem", Path: "./deploy/models"},
		Observability: Observability{TraceSampleRate: 1.0, MetricsEnabled: true, OTelEndpoint: "localhost:4318"},
		LLM:           LLM{PrimaryProvider: "qwen", FallbackProvider: "deepseek", DailyQuotaFree: 10},
	}
}
