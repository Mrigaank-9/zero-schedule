package main

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Mrigaank-9/job-scheduler/internal/api"
	"github.com/Mrigaank-9/job-scheduler/internal/config"
	"github.com/Mrigaank-9/job-scheduler/internal/database/sqlite"
	"github.com/Mrigaank-9/job-scheduler/internal/queue/memory"
	"github.com/Mrigaank-9/job-scheduler/internal/repositiory"
	"github.com/Mrigaank-9/job-scheduler/internal/worker"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("error loading .env file")
	}

	cfg := config.MustLoadConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	jobQueue := memory.NewQueue()
	db := sqlite.SetupSQL(os.Getenv("DB_URI"))

	if err := repositiory.InitJobsTable(ctx, db); err != nil {
		log.Fatal(err)
	}

	go worker.RunWorkers(ctx, jobQueue, cfg.Workers, db)

	router := http.NewServeMux()

	router.HandleFunc(
		"POST /api/v1/job",
		api.New(ctx, db, jobQueue),
	)

	server := http.Server{
		Addr:    cfg.Address,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			slog.Error(
				"server failed",
				slog.String("error", err.Error()),
			)
		}
	}()

	slog.Info(
		"server started",
		slog.String("address", cfg.Address),
	)

	done := make(chan os.Signal, 1)
	signal.Notify(
		done,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)

	<-done

	slog.Info("shutting down server")

	// Stop workers.
	cancel()

	// Give HTTP server time to finish active requests.
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error(
			"failed to shutdown server",
			slog.String("error", err.Error()),
		)
	}

	slog.Info("server shutdown successfully")
}
