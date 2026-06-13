// Flywheel worker consumes NATS triggers and runs ML retrain pipeline.
package main

import (
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

func main() {
	url := os.Getenv("NATS_URL")
	if url == "" {
		url = "nats://localhost:4222"
	}
	mlRoot := os.Getenv("ML_ROOT")
	if mlRoot == "" {
		mlRoot = "/app/ml"
	}
	nc, err := nats.Connect(url)
	if err != nil {
		log.Fatal(err)
	}
	_, err = nc.QueueSubscribe("msgguard.flywheel.trigger", "flywheel-workers", func(m *nats.Msg) {
		log.Printf("flywheel trigger: %s", string(m.Data))
		run := exec.Command("bash", "-lc", "make flywheel && make train && make benchmark")
		run.Dir = mlRoot
		run.Stdout = os.Stdout
		run.Stderr = os.Stderr
		if err := run.Run(); err != nil {
			log.Printf("flywheel failed: %v", err)
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("flywheel worker listening on msgguard.flywheel.trigger")
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
