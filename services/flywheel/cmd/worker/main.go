// Flywheel worker consumes NATS triggers and runs ML retrain pipeline with debounce.
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	url := envOr("NATS_URL", "nats://localhost:4222")
	mlRoot := envOr("ML_ROOT", "/app/ml")
	minInterval := durationOr("FLYWHEEL_MIN_INTERVAL", time.Hour)
	minBatch := intOr("FLYWHEEL_MIN_BATCH", 10)

	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}

	var (
		mu      sync.Mutex
		lastRun time.Time
		pending int
	)

	runPipeline := func() {
		log.Println("flywheel: starting pipeline")
		run := exec.Command("bash", "-lc", "make flywheel && make train-all-locales && make export-all-locales && make benchmark")
		run.Dir = mlRoot
		run.Env = append(os.Environ(), "CI=true")
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		if err := run.Run(); err != nil {
			log.Printf("flywheel failed: %v", err)
			return
		}
		log.Println("flywheel: pipeline complete")
	}

	_, err = nc.QueueSubscribe("msgguard.flywheel.trigger", "flywheel-workers", func(m *nats.Msg) {
		log.Printf("flywheel trigger: %s", string(m.Data))
		mu.Lock()
		pending++
		shouldRun := pending >= minBatch || time.Since(lastRun) >= minInterval
		if shouldRun {
			pending = 0
			lastRun = time.Now()
		}
		mu.Unlock()
		if shouldRun {
			go runPipeline()
		} else {
			log.Printf("flywheel: debounced (%d pending, next after %s)", pending, minInterval)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("flywheel worker listening (debounce: %s or %d events)", minInterval, minBatch)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func envOr(k, d string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return d
}

func durationOr(k string, d time.Duration) time.Duration {
	if v := os.Getenv(k); v != "" {
		if parsed, err := time.ParseDuration(v); err == nil {
			return parsed
		}
	}
	return d
}

func intOr(k string, d int) int {
	if v := os.Getenv(k); v != "" {
		var n int
		if _, err := fmt.Sscanf(v, "%d", &n); err == nil {
			return n
		}
	}
	return d
}
