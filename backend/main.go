// Command ConMonitr backend serves the container-monitoring REST + WebSocket API.
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"conmonitr/backend/internal/api"
	"conmonitr/backend/internal/config"
	"conmonitr/backend/internal/docker"
)

func main() {
	config.Load(".env")

	addr := os.Getenv("CONMONITR_ADDR")
	if addr == "" {
		addr = ":8081"
	}

	svc, err := docker.NewService()
	if err != nil {
		log.Fatalf("docker client: %v", err)
	}
	defer svc.Close()

	if version, err := svc.Ping(context.Background()); err != nil {
		log.Printf("warning: docker not reachable yet: %v", err)
	} else {
		log.Printf("connected to docker engine, API version %s", version)
	}

	handler := api.NewHandler(svc)
	router := api.NewRouter(handler)

	srv := &http.Server{
		Addr:    addr,
		Handler: router,
	}

	go func() {
		log.Printf("ConMonitr backend listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Println("shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}
}
