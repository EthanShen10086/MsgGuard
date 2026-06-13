// Package main is the MsgGuard API Gateway entry point.
package main

import (
	"os"

	"github.com/EthanShen10086/msgguard/pkg/app"
	"github.com/EthanShen10086/msgguard/pkg/config"
	"github.com/EthanShen10086/msgguard/services/gateway/internal/server"
)

func main() {
	cfgPath := envOr("CONFIG_PATH", "../../deploy/config.yaml")
	cfg, err := config.Load(cfgPath)
	if err != nil {
		panic(err)
	}
	container, err := app.NewContainer(cfg)
	if err != nil {
		panic(err)
	}
	if err := server.Run(container); err != nil {
		panic(err)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
