package jobs

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	JobID     uuid.UUID
	Command   string
	CreatedAt time.Time
	UpdatedAt time.Time
	StartedAt time.Time
	EndedAt   time.Time
}

func NewJob(command string) *Job {
	job := Job{
		JobID:     uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &job
}
