package main

import (
	"context"
	"expense-tracker/internal/app"
	"expense-tracker/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 1. Load Configuration
	cfg := config.Load()

	// 2. Setup Logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))

	// 3. Initialize Application
	application, err := app.New(cfg, logger)
	if err != nil {
		logger.Error("failed to initialize app", "error", err)
		os.Exit(1)
	}

	// 4. Setup Context with Signal Handling
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Listen for interrupt signals (Ctrl+C, kill)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger.Info("received shutdown signal")
		cancel() // Triggers graceful shutdown in app.Run
	}()

	// 5. Run Application
	if err := application.Run(ctx); err != nil {
		logger.Error("application error", "error", err)
		os.Exit(1)
	}
}
