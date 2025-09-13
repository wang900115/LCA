package bootstrap

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func Run(cancelTime time.Duration, srv *http.Server, scheduler gocron.Scheduler) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[SYSTEM] HTTP server error: %s", err.Error())
		}
		log.Println("[SYSTEM] HTTP server init successfully")
	}()

	go func() {
		scheduler.Start()
		log.Println("[SYSTEM] Scheduler init successfully")
	}()

	<-ctx.Done()
	log.Println("[SYSTEM] Shutdown signal received ...")

	shutdownCTX, cancel := context.WithTimeout(context.Background(), cancelTime)
	defer cancel()
	if err := srv.Shutdown(shutdownCTX); err != nil {
		log.Fatalf("[SYSTEM] HTTP server forced to shutdown: %s", err.Error())
	}
	log.Println("[SYSTEM] HTTP server exited gracefully")

	if err := scheduler.StopJobs(); err != nil {
		log.Fatalf("[SYSTEM] Scheduler forced exited: %s", err.Error())
	}
	log.Println("[SYSTEM] Scheduler exited gracefully")
}
