package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/tonnomolt/hortons-diary/internal/config"
	"github.com/tonnomolt/hortons-diary/internal/httpserver"
	"github.com/tonnomolt/hortons-diary/internal/attacks"
	webassets "github.com/tonnomolt/hortons-diary/web"
)

const shutdownTimeout = 10 * time.Second

func main() {
	if len(os.Args) == 3 && os.Args[1] == "healthcheck" {
		os.Exit(runHealthcheck(os.Args[2]))
	}

	if err := run(); err != nil {
		slog.Error("server stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	//TEMP TEST BEGINS
	record, err := attacks.NewPainRecord(time.Now(), 8, 45)
	if err != nil {
		return fmt.Errorf("create demo record: %w", err)
	}
	fmt.Printf("Demo PainRecord: %+v
	", record)
	//TEMP TEST ENDS

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load configuration: %w", err)
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: cfg.LogLevel}))
	slog.SetDefault(logger)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		return fmt.Errorf("configure database pool: %w", err)
	}
	defer db.Close()

	if err := db.Ping(ctx); err != nil {
		return fmt.Errorf("connect to database: %w", err)
	}

	handler, err := httpserver.NewHandler(db, webassets.Files)
	if err != nil {
		return fmt.Errorf("create HTTP handler: %w", err)
	}

	server := httpserver.New(cfg.HTTPAddr, handler)
	errCh := make(chan error, 1)
	go func() {
		logger.Info("HTTP server listening", "address", cfg.HTTPAddr, "environment", cfg.Environment)
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-ctx.Done():
		logger.Info("shutdown requested")
	case err := <-errCh:
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("serve HTTP: %w", err)
		}
		return nil
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("graceful shutdown: %w", err)
	}
	return nil
}

func runHealthcheck(url string) int {
	client := http.Client{Timeout: 2 * time.Second}
	response, err := client.Get(url)
	if err != nil {
		return 1
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return 1
	}
	return 0
}
