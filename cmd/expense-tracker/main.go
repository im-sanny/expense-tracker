package main

import (
	"expense-tracker/internal/config"
	"expense-tracker/internal/db"
	"expense-tracker/internal/handler"
	"expense-tracker/internal/repository"
	"log"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	cfg := config.Load()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if err := db.RunMigrations(db.MigrateConfig{
		DBURL:         cfg.DB.DSN(),
		MigrationPath: "file://migrations",
		Logger:        logger,
	}); err != nil {
		logger.Error("migration failed", "error", err)
		os.Exit(1)
	}
	logger.Info("✓ migrations complete, starting application")

	database, err := db.New(cfg.DB.DSN())
	if err != nil {
		log.Fatalf("cannot connect to database: %v", err)
	}
	defer database.Close()

	repo := repository.NewExpenseRepo(database)
	h := handler.NewHandler(repo)
	mux := http.NewServeMux()

	mux.HandleFunc("GET /track", h.Get)
	mux.HandleFunc("GET /track/{id}", h.GetById)
	mux.HandleFunc("POST /track", h.Post)
	mux.HandleFunc("PUT /track/{id}", h.Put)
	mux.HandleFunc("PATCH /track/{id}", h.Patch)
	mux.HandleFunc("DELETE /track/{id}", h.Delete)

	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "3000"
	}
	
	addr := ":" + port

	log.Printf("Server running on port http://localhost%s", addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
