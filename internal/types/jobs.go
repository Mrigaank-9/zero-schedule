package types

import (
	"time"

	"github.com/google/uuid"
)

type JobStatus string

const (
	QUEUED  JobStatus = "QUEUED"
	RUNNING JobStatus = "RUNNING"
	DONE    JobStatus = "DONE"
	FAILED  JobStatus = "FAILED"
)

type Job struct {
	JobID     uuid.UUID
	Name      string
	Command   string
	Status    JobStatus
	CreatedAt time.Time
	UpdatedAt time.Time
	StartedAt *time.Time
	EndedAt   *time.Time
}

func NewJob(name string, command string) *Job {
	job := Job{
		JobID:     uuid.New(),
		Command:   command,
		Name:      name,
		Status:    QUEUED,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return &job
}

type RequestJobcreation struct {
	Name    string
	Command string
}
