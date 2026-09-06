package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/Mrigaank-9/job-scheduler/internal/database"
	"github.com/Mrigaank-9/job-scheduler/internal/jobs"
	"github.com/Mrigaank-9/job-scheduler/internal/queue"
	"github.com/Mrigaank-9/job-scheduler/internal/types"
	"github.com/Mrigaank-9/job-scheduler/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(ctx context.Context, db database.Database, queue queue.Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var jobRequest types.RequestJobcreation

		// decode the json into struct
		err := json.NewDecoder(r.Body).Decode(&jobRequest)
		if errors.Is(err, io.EOF) {
			// empty bodyjob
			response.WriteJson(w, http.StatusBadRequest, response.GenernalError(fmt.Errorf("empty body")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GenernalError(err))
			return
		}

		// request validation to use
		slog.Info("creating a student")
		if err := validator.New().Struct(jobRequest); err != nil {
			validationErr := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationErrors(validationErr))
			return
		}
		job, err := jobs.CreateJob(ctx, db, jobRequest.Name, jobRequest.Command, queue)
		if err != nil {
			response.WriteJson(w, http.StatusInternalServerError, err)
			return
		}

		slog.Info(
			"user created successfully",
			"id", job.JobID.String(),
			"name", job.Name,
			"command", job.Command,
		)
		response.WriteJson(w, http.StatusCreated, map[string]string{
			"id": job.JobID.String(),
		})
	}
}
