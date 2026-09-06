# Job Scheduler

A lightweight job scheduler built in Go that accepts jobs through an HTTP API, persists job metadata in SQLite, queues jobs in memory, and executes them concurrently using a configurable worker pool.

## Overview

The scheduler follows a simple execution pipeline:

```text
HTTP API
   │
   ▼
Create Job
   │
   ├──► SQLite
   │
   └──► In-Memory Queue
            │
            ▼
      Worker Pool
            │
            ▼
      OS Subprocess
            │
            ▼
    DONE / FAILED
```

The project is currently at **V1**, focused on building the core scheduling and execution system.

## Features

* REST API for creating jobs
* SQLite persistence for job metadata
* In-memory buffered job queue
* Configurable worker pool
* Concurrent job execution using goroutines
* OS subprocess execution using Go's `os/exec`
* Job lifecycle tracking
* Context-aware queue operations
* Graceful HTTP server shutdown
* Graceful worker shutdown through context cancellation
* Basic error handling and structured logging using `log/slog`

## Job Lifecycle

Each job moves through the following states:

```text
PENDING
   │
   ▼
RUNNING
   │
   ├──────────► DONE
   │
   └──────────► FAILED
```

The scheduler stores important execution timestamps including:

* `created_at`
* `updated_at`
* `started_at`
* `ended_at`

## Architecture

The project is structured into separate layers:

```text
cmd/
└── jobSchedulor/
    └── main.go

internal/
├── api/
├── config/
├── database/
│   └── sqlite/
├── jobs/
├── queue/
│   ├── memory/
│   └── queue.go
├── repositiory/
├── types/
└── worker/
```

### API Layer

Responsible for receiving HTTP requests and exposing job-related endpoints.

Current endpoint:

```http
POST /api/v1/job
```

### Job Service

Handles job creation and coordinates persistence with queue publishing.

```text
Create Job
    ↓
Persist job in SQLite
    ↓
Publish job ID to queue
```

### Queue

The queue is defined through an interface so the underlying implementation can be replaced later.

Current implementation:

```text
MemoryQueue
```

The queue uses a buffered Go channel with capacity `1024`.

Workers block while waiting for jobs rather than continuously polling.

### Worker Pool

The scheduler starts a configurable number of workers:

```text
Worker 1 ──┐
Worker 2 ──┤
Worker 3 ──┤──► Queue
Worker N ──┘
```

Each worker:

1. Consumes a job ID from the queue
2. Loads the job from SQLite
3. Marks the job as `RUNNING`
4. Records `started_at`
5. Spawns the job as an OS subprocess
6. Waits for the subprocess to finish
7. Marks the job as `DONE` or `FAILED`
8. Records `ended_at`

## Subprocess Execution

Jobs are executed using Go's `os/exec` package.

For example:

```bash
/opt/homebrew/bin/python3 /path/to/run.py
```

The worker starts the process and waits for completion:

```go
cmd := exec.Command("sh", "-c", job.Command)

if err := cmd.Start(); err != nil {
    // process failed to start
}

if err := cmd.Wait(); err != nil {
    // process failed
}
```

The worker goroutine waits for its subprocess, while other workers remain available to execute other jobs.

## Configuration

Configuration is loaded from environment variables.

Example:

```env
ADDRESS=:8080
WORKERS=3
DATABASE_PATH=storage/sqlite.db
```

See `.env.example` for the available configuration values.

## Running Locally

### 1. Clone the repository

```bash
git clone <repository-url>
cd job-scheduler
```

### 2. Configure environment variables

Create a `.env` file:

```env
ADDRESS=:8080
WORKERS=3
DATABASE_PATH=storage/sqlite.db
```

### 3. Run the application

```bash
go run cmd/jobSchedulor/main.go
```

The server will start on the configured address.

## Example

Create a job:

```http
POST /api/v1/job
Content-Type: application/json
```

Example request:

```json
{
  "name": "test-job",
  "command": "/opt/homebrew/bin/python3 /path/to/run.py"
}
```

The job is persisted and published to the queue.

A worker then picks it up and executes the command.

## Graceful Shutdown

The application listens for:

```text
SIGINT
SIGTERM
```

When shutdown is requested:

1. The worker context is cancelled.
2. Workers stop waiting for new jobs.
3. The HTTP server stops accepting new requests.
4. Existing HTTP requests are given time to finish.
5. The application exits cleanly.

## Current Scope — V1

V1 intentionally focuses on the core scheduling functionality.

### Implemented

* Job creation
* Job persistence
* In-memory queue
* Worker pool
* Concurrent execution
* Subprocess management
* Job status tracking
* Execution timestamps
* Graceful shutdown

### Planned for V2

* Persistent execution logs
* stdout/stderr capture
* Job log retrieval API
* Job execution duration/progress information
* Retry failed jobs
* Retry/attempt tracking
* Job filtering and querying
* More robust process management
* Persistent/distributed queue

Future versions may replace the in-memory queue with infrastructure such as Redis and introduce additional distributed-system components.

## Tech Stack

* **Go**
* **net/http**
* **SQLite**
* **database/sql**
* **os/exec**
* **Goroutines**
* **Channels**
* **Context**
* **slog**
* **godotenv**

## Project Goal

The goal of this project is to understand and implement the core components of a job execution system while exploring Go's concurrency primitives, worker pools, queues, process management, persistence, and graceful shutdown.

The project will be incrementally evolved toward a more production-oriented distributed job scheduling system.
