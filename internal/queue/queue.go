package queue

import "context"

type Queue interface {
	Publish(ctx context.Context, jobId string) error
	Consume(ctx context.Context) (string, error)
	Ack(ctx context.Context, jobId string) error
	Retry(ctx context.Context, jobId string) error
}
