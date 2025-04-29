package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"task-dashboard/config"
	"time"

	"task-dashboard/internal/api"
	"task-dashboard/internal/service"
	"task-dashboard/internal/web"
)

func fatalLogAndExit(msg string, err error) {
	slog.Error(msg, "error", err)
	os.Exit(1)
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	cfg, err := config.LoadConfig()
	if err != nil {
		fatalLogAndExit("Failed to load configuration", err)
	}

	// Create a root context that cancels on OS interrupt signals
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Create API client
	apiClient := api.NewMotionClient(cfg.APIKey)

	// Validate API key
	slog.Info("Validating API key...")
	if err := apiClient.ValidateAPIKey(); err != nil {
		fatalLogAndExit("API key validation failed", err)
	}
	slog.Info("API key is valid!")

	// Create services
	taskService := service.NewTaskService(apiClient)
	go taskService.StartPeriodicRefresh(ctx, cfg.RefreshInterval)

	slog.Info("Fetching initial task data", "refresh_interval", cfg.RefreshInterval)
	if _, err := taskService.RefreshTasks(); err != nil {
		fatalLogAndExit("Error fetching initial tasks", err)
	}

	server, err := web.NewServer(taskService, cfg.TemplatesDir, cfg.Port)
	if err != nil {
		fatalLogAndExit("Error creating server", err)
	}

	server.StartMetricsCollector(ctx)

	// Start server in background
	go func() {
		slog.Info("Starting server", "port", cfg.Port)
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			fatalLogAndExit("Server error", err)
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	slog.Info("Shutdown signal received.")

	// Create a new context for server shutdown (e.g., 5 seconds timeout)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		fatalLogAndExit("Server shutdown failed", err)
	}

	slog.Info("Shutdown complete. Goodbye.")
}
