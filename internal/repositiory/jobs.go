package repositiory

import (
	"context"
	"errors"
	"time"

	"github.com/Mrigaank-9/job-scheduler/internal/database"
	"github.com/Mrigaank-9/job-scheduler/internal/types"
	"github.com/google/uuid"
)

func InitJobsTable(ctx context.Context, db database.Database) error {
	query := `
		CREATE TABLE IF NOT EXISTS jobs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			command TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at DATETIME NOT NULL,
			updated_at DATETIME NOT NULL,
			started_at DATETIME,
			ended_at DATETIME
		)
	`

	_, err := db.Exec(ctx, query)
	return err
}

func CreateJob(ctx context.Context, db database.Database, job *types.Job) error {
	query := `
		INSERT INTO jobs (id, name, command, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`

	_, err := db.Exec(
		ctx,
		query,
		job.JobID.String(),
		job.Name,
		job.Command,
		job.Status,
		job.CreatedAt,
		job.UpdatedAt,
	)

	return err
}

func GetJobById(ctx context.Context, db database.Database, jobId string) (types.Job, error) {
	query := `
		SELECT id, name, command, status, created_at, updated_at , started_at, ended_at
		FROM jobs
		WHERE id = ?
	`

	var job types.Job
	var id string

	err := db.QueryRow(ctx, query, jobId).Scan(
		&id,
		&job.Name,
		&job.Command,
		&job.Status,
		&job.CreatedAt,
		&job.UpdatedAt,
		&job.StartedAt,
		&job.EndedAt,
	)

	if err != nil {
		return types.Job{}, err
	}

	job.JobID, err = uuid.Parse(id)
	if err != nil {
		return types.Job{}, err
	}

	return job, nil
}

func UpdateStatus(ctx context.Context, db database.Database, job *types.Job, status types.JobStatus) error {
	now := time.Now()
	job.UpdatedAt = now
	query := `
		UPDATE jobs SET status = ?,  updated_at = ? 
		WHERE id = ? 
	`
	result, err := db.Exec(ctx, query, status, job.UpdatedAt, job.JobID.String())
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("job not found")
	}
	return nil
}

func UpdateStartedAt(ctx context.Context, db database.Database, job *types.Job) error {
	now := time.Now()
	job.StartedAt = &now
	job.UpdatedAt = now
	query := `
		UPDATE jobs SET started_at = ?,  updated_at = ? 
		WHERE id = ? 
	`
	result, err := db.Exec(ctx, query, job.StartedAt, job.UpdatedAt, job.JobID.String())
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("job not found")
	}
	return nil
}

func UpdateEndedAt(ctx context.Context, db database.Database, job *types.Job) error {
	now := time.Now()
	job.EndedAt = &now
	job.UpdatedAt = now
	query := `
		UPDATE jobs SET ended_at = ?,  updated_at = ? 
		WHERE id = ? 
	`
	result, err := db.Exec(ctx, query, job.EndedAt, job.UpdatedAt, job.JobID.String())
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("job not found")
	}
	return nil
}
