package worker

import (
	"context"
	"log/slog"
	"os/exec"
	"sync"

	"github.com/Mrigaank-9/job-scheduler/internal/database"
	"github.com/Mrigaank-9/job-scheduler/internal/queue"
	"github.com/Mrigaank-9/job-scheduler/internal/repositiory"
	"github.com/Mrigaank-9/job-scheduler/internal/types"
)

func work(ctx context.Context, workerId int64, wg *sync.WaitGroup, jobQueue queue.Queue, db database.Database) {
	defer wg.Done()
	slog.Debug("worker running", slog.Int64("worker-id", workerId))
	for {
		jobId, err := jobQueue.Consume(ctx)
		if err != nil {
			slog.Error("unable to get job", slog.String("job-id", jobId), slog.String("error", err.Error()))
			continue
		}

		// Get job from DB
		job, err := repositiory.GetJobById(ctx, db, jobId)
		if err != nil {
			slog.Error("unable to get job", slog.String("job-id", jobId), slog.String("error", err.Error()))
			continue
		}

		// Mark job as RUNNING
		if err := repositiory.UpdateStatus(ctx, db, &job, types.RUNNING); err != nil {
			slog.Error("unable to update job status to running", slog.String("job-id", jobId), slog.String("error", err.Error()))
			continue
		}

		// Set started_at
		if err := repositiory.UpdateStartedAt(ctx, db, &job); err != nil {
			slog.Error("unable to update job started_at", slog.String("job-id", jobId), slog.String("error", err.Error()))
			continue
		}

		// Spawn subprocess
		cmd := exec.Command(job.Command)

		if err := cmd.Start(); err != nil {
			slog.Error("unable to start job process", slog.String("job-id", jobId), slog.String("error", err.Error()))

			// Mark job as failed
			if updateErr := repositiory.UpdateStatus(ctx, db, &job, types.FAILED); updateErr != nil {
				slog.Error("unable to update job status to failed", slog.String("job-id", jobId), slog.String("error", updateErr.Error()))
			}
			continue
		}

		// Wait for subprocess to finish
		if err := cmd.Wait(); err != nil {
			slog.Error("job process failed", slog.String("job-id", jobId), slog.String("error", err.Error()))
			if updateErr := repositiory.UpdateStatus(ctx, db, &job, types.FAILED); updateErr != nil {
				slog.Error("unable to update job status to failed", slog.String("job-id", jobId), slog.String("error", updateErr.Error()))
			}

			if updateErr := repositiory.UpdateEndedAt(ctx, db, &job); updateErr != nil {
				slog.Error("unable to update job ended_at", slog.String("job-id", jobId), slog.String("error", updateErr.Error()))
			}

			continue
		}

		// Process completed successfully
		if err := repositiory.UpdateStatus(ctx, db, &job, types.DONE); err != nil {
			slog.Error("unable to update job status to done", slog.String("job-id", jobId), slog.String("error", err.Error()))
			continue
		}

		if err := repositiory.UpdateEndedAt(ctx, db, &job); err != nil {
			slog.Error("unable to update job ended_at", slog.String("job-id", jobId), slog.String("error", err.Error()))
			continue
		}
	}
}

func RunWorkers(ctx context.Context, jobQueue queue.Queue, workerCount int64, db database.Database) {
	var wg sync.WaitGroup
	for i := range workerCount {
		wg.Add(1)
		go work(ctx, i, &wg, jobQueue, db)
	}

	wg.Wait()
}
