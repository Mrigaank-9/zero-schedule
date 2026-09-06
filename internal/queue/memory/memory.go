package memory

import (
	"context"

	"github.com/Mrigaank-9/job-scheduler/internal/queue"
)

type MemoryQueue struct {
	Jobs chan string
}

func NewQueue() queue.Queue {
	jobChan := make(chan string, 1024)
	return &MemoryQueue{Jobs: jobChan}
}

func (m *MemoryQueue) Publish(ctx context.Context, jobId string) error {
	select {
	case m.Jobs <- jobId:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (m *MemoryQueue) Consume(ctx context.Context) (string, error) {
	select {
	case jobId := <-m.Jobs:
		return jobId, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
func (m *MemoryQueue) Ack(ctx context.Context, jobId string) error {
	return nil

}
func (m *MemoryQueue) Retry(ctx context.Context, jobId string) error {
	return nil

}
