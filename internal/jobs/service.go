package jobs

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/Mrigaank-9/job-scheduler/internal/queue"
)

func CreateJob(ctx context.Context, command string, jobQueue queue.Queue) error {
	if command == "" {
		return errors.New("command is not provided")
	}
	job := NewJob(command)
	jobId := strconv.FormatUint(uint64(job.JobID.ID()), 10)
	fmt.Println("Created Job as Job Id: ", jobId)
	err := jobQueue.Publish(ctx, jobId)
	if err != nil {
		return err
	}
	return nil
}
