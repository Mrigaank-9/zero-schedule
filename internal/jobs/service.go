package jobs

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Mrigaank-9/job-scheduler/internal/database"
	"github.com/Mrigaank-9/job-scheduler/internal/queue"
	"github.com/Mrigaank-9/job-scheduler/internal/repositiory"
	"github.com/Mrigaank-9/job-scheduler/internal/types"
)

func CreateJob(ctx context.Context, db database.Database, name string, command string, jobQueue queue.Queue) error {
	if command == "" {
		return errors.New("command not provided")
	}
	if name == "" {
		return errors.New("name not provided")
	}
	job := types.NewJob(name, command)
	jobId := strconv.FormatUint(uint64(job.JobID.ID()), 10)
	repositiory.CreateJob(ctx, db, job)
	err := jobQueue.Publish(ctx, jobId)
	if err != nil {
		return err
	}
	fmt.Println("Created Job as Job Id: ", jobId)
	return nil
}
