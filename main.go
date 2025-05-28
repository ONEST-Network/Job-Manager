package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ONEST-Network/scheme-manager-adapter/api/routes"
	"github.com/ONEST-Network/scheme-manager-adapter/pkg/config"
	"github.com/ONEST-Network/scheme-manager-adapter/pkg/database/postgres"
	"github.com/sirupsen/logrus"
)

func init() {
	// Configure logger
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetOutput(os.Stdout)

	// Set log level based on environment variable
	logLevel := os.Getenv("LOG_LEVEL")
	if logLevel == "" {
		logLevel = "info"
	}
	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		logrus.WithError(err).WithField("log_level", logLevel).Warn("Invalid log level provided, falling back to INFO level")
		logrus.SetLevel(logrus.InfoLevel)
	} else {
		logrus.SetLevel(level)
	}
}

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		logrus.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	ctx := context.Background()
	db, err := postgres.InitDB(ctx, cfg)
	if err != nil {
		logrus.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Set up router
	router := routes.SetupRouter(db, cfg)

	// Create server
	serverAddr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	// Run server in a goroutine
	go func() {
		logrus.Infof("Starting server on %s", serverAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logrus.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// kill (no param) default sends syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logrus.Info("Shutting down server...")

	// Create context with timeout for shutdown
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Server.GracefulTimeout)*time.Second)
	defer cancel()

	// Shutdown server
	if err := server.Shutdown(shutdownCtx); err != nil {
		logrus.Fatalf("Server forced to shutdown: %v", err)
	}
	logrus.Info("Server exiting")
}
