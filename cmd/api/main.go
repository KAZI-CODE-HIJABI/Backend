package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/KAZI-CODE-HIJABI/Backend/internal/config"
	"github.com/KAZI-CODE-HIJABI/Backend/internal/database"
	"github.com/KAZI-CODE-HIJABI/Backend/internal/server"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	level := zapcore.InfoLevel
	_ = level.UnmarshalText([]byte(cfg.LogLevel))
	lc := zap.NewProductionConfig()
	lc.Level = zap.NewAtomicLevelAt(level)
	log, err := lc.Build()
	if err != nil {
		return errors.New("cannot initialize logger")
	}
	defer func() { _ = log.Sync() }()
	gin.SetMode(gin.ReleaseMode)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	connectCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	_, pool, err := database.Open(connectCtx, cfg.DatabaseURL)
	cancel()
	if err != nil {
		return errors.New("database connection failed; check DATABASE_URL and database availability")
	}
	defer pool.Close()
	srv := &http.Server{Addr: cfg.HTTPAddress, Handler: server.Router(pool, log), ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 15 * time.Second, IdleTimeout: 60 * time.Second, MaxHeaderBytes: 1 << 20}
	errs := make(chan error, 1)
	go func() { log.Info("api_started", zap.String("address", cfg.HTTPAddress)); errs <- srv.ListenAndServe() }()
	select {
	case err := <-errs:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			_ = srv.Close()
			return err
		}
		log.Info("api_stopped")
		return nil
	}
}
