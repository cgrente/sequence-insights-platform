package main

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/cgrente/sequence-insights-platform/server/internal/config"
	"github.com/cgrente/sequence-insights-platform/server/internal/db"
	"github.com/cgrente/sequence-insights-platform/server/internal/handlers"
	"github.com/cgrente/sequence-insights-platform/server/internal/jobs"
	"github.com/cgrente/sequence-insights-platform/server/internal/logging"
	"github.com/cgrente/sequence-insights-platform/server/internal/storage"
)

/*
cmd/api is the HTTP entrypoint.

Best practices demonstrated:
- Config via env
- Structured logs
- Dependency injection (Store/Queue)
- Graceful shutdown
*/

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	log := logging.NewJSONLogger(slog.LevelInfo)
	log.Info(
		"starting",
		"env", cfg.Env,
		"addr", net.JoinHostPort(cfg.HTTPHost, strconv.Itoa(cfg.HTTPPort)),
	)
	// Database
	dbConn, err := db.Open(cfg.DBUrl)
	if err != nil {
		log.Error("failed to open database", "err", err)
		os.Exit(1)
	}
	defer dbConn.Close()

	dbConn.SetMaxOpenConns(cfg.DBMaxOpenConns)
	dbConn.SetMaxIdleConns(cfg.DBMaxIdleConns)
	dbConn.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)
	dbConn.SetConnMaxLifetime(cfg.DBConnMaxLifetime)

	ctx := context.Background()
	if err := db.Ping(ctx, dbConn); err != nil {
		log.Error("database ping failed", "err", err)
		os.Exit(1)
	}

	store := storage.NewPostgresStore(dbConn)

	// Jobs
	queue := jobs.NewQueue(log, store, cfg.JobQueueSize)
	workerCtx, workerCancel := context.WithCancel(context.Background())
	queue.RunWorkers(workerCtx, cfg.WorkerCount)
	defer workerCancel()

	// HTTP
	h := handlers.NewRouter(handlers.RouterParams{
		Log:    log,
		Store:  store,
		Queue:  queue,
		APIKey: cfg.APIKey,
	})

	srv := &http.Server{
		Addr:              net.JoinHostPort(cfg.HTTPHost, strconv.Itoa(cfg.HTTPPort)),
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("http listening", "addr", net.JoinHostPort(cfg.HTTPHost, strconv.Itoa(cfg.HTTPPort)))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server error", "err", err)
			os.Exit(1)
		}
	}()

	// Graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	log.Info("shutting down")
	workerCancel()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Info("bye")
}
