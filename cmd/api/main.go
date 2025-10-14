package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/szks-repo/cloud-run-blog/internal/blog"
	"github.com/szks-repo/cloud-run-blog/internal/server"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	addr := ":" + envOrDefault("PORT", "8080")

	repo := blog.NewInMemoryRepository()
	srv, err := server.New(repo)
	if err != nil {
		log.Fatalf("initializing server: %v", err)
	}

	if err := srv.Run(ctx, addr); err != nil && !errors.Is(err, context.Canceled) && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

func envOrDefault(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
