package app

import (
	"context"
	"database/sql"
	"errors"
	"expense-tracker/internal/config"
	"expense-tracker/internal/db"
	"expense-tracker/internal/handler"
	"expense-tracker/internal/middlewares"
	"expense-tracker/internal/repository"
	"expense-tracker/internal/service"
	"log/slog"
	"net/http"
	"time"
)

// App represents the running application
type App struct {
	logger  *slog.Logger
	server  *http.Server
	db      *sql.DB
	service *service.ExpenseService
}

// New initializes the application and wires dependencies
func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	// 1. Run Migrations
	if err := db.RunMigrations(db.MigrateConfig{
		DBURL:         cfg.DB.DSN(),
		MigrationPath: "file://migrations",
		Logger:        logger,
	}); err != nil {
		return nil, err
	}
	logger.Info("✓ migrations complete")

	// 2. Database Connection
	database, err := db.New(cfg.DB.DSN())
	if err != nil {
		return nil, err
	}

	// 3. Wire Dependencies
	repo := repository.NewExpenseRepo(database)
	svc := service.NewExpenseService(repo, nil)
	h := handler.NewHandler(svc)
	mux := http.NewServeMux()

	authRepo := repository.NewAuthRepo(database)
	authService := service.NewAuthService(cfg.JWTSecret, authRepo)
	authHandler := handler.NewAuthHandler(authService)

	mux.HandleFunc("POST /auth/register", authHandler.Register)
	mux.HandleFunc("POST /auth/login", authHandler.Login)
	mux.HandleFunc("POST /auth/logout", authHandler.Logout)
	mux.HandleFunc("POST /auth/refresh", authHandler.Refresh)

	// 4. Register Routes
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("GET /track", h.Get)
	protectedMux.HandleFunc("GET /track/{id}", h.GetById)
	protectedMux.HandleFunc("POST /track", h.Post)
	protectedMux.HandleFunc("PUT /track/{id}", h.Put)
	protectedMux.HandleFunc("PATCH /track/{id}", h.Patch)
	protectedMux.HandleFunc("DELETE /track/{id}", h.Delete)

	protectedHandler := middlewares.AuthMiddleware(authService)(protectedMux)
	mux.Handle("/track", protectedHandler)
	mux.Handle("/track/", protectedHandler)

	handler := middlewares.JSONMiddleware(mux)
	handler = middlewares.Cors(handler)
	handler = middlewares.LoggingMiddleware(logger)(handler)
	handler = middlewares.TimeoutMiddleware(30 * time.Second)(handler)
	handler = middlewares.Recover(logger)(handler)

	// 5. Setup Server
	server := &http.Server{
		Addr:         ":" + cfg.HTTP.Port,
		Handler:      handler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &App{
		logger:  logger,
		server:  server,
		db:      database,
		service: svc,
	}, nil
}

// Run starts the application and blocks until context is cancelled
func (a *App) Run(ctx context.Context) error {
	// Start server in a goroutine
	errChan := make(chan error, 1)
	go func() {
		a.logger.Info("server starting", "addr", a.server.Addr)
		errChan <- a.server.ListenAndServe()
	}()

	// Wait for either error or context cancellation
	select {
	case err := <-errChan:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		a.logger.Info("shutting down server...")
		// Graceful shutdown with timeout
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return err
		}
	}

	// Cleanup resources
	if a.db != nil {
		a.db.Close()
	}

	a.logger.Info("application stopped")
	return nil
}
