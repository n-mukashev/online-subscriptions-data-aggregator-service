package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/config"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/logger"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/routes"
	"github.com/mukashev-n/online-subscriptions-data-aggregator-service/internal/storage"
)

// @title Users Online Subscriptions Data Aggregator API
// @version 1.0
// @description API documentation for Users Online Subscriptions Data Aggregator
// @BasePath /
func main() {
	cfg := config.MustLoad()
	logger.InitLogger(cfg.Env)
	storage.InitDB(cfg)
	server := gin.Default()
	routes.RegisterRoutes(server)
	srv := &http.Server{
		Addr:    cfg.ServerConfig.Url,
		Handler: server,
	}
	go func() {
		logger.Log.Info("server started", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Error("server failed", "err", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Log.Info("shutdown signal received")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Log.Error("server forced to shutdown", "err", err)
	} else {
		logger.Log.Info("server stopped gracefully")
	}
	if storage.DB != nil {
		if err := storage.DB.Close(); err != nil {
			logger.Log.Error("error closing db", "err", err)
		} else {
			logger.Log.Info("db connection closed")
		}
	}
}
