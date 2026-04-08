package graceful_shutdown

import (
	"context"
	"sync"
	"time"
)

// UpdateServer simulates a package update server that processes update jobs.
// When Shutdown is called (or the context is cancelled), it should:
//   1. Stop accepting new jobs.
//   2. Wait for all in-flight jobs to complete (with a deadline).
//   3. Return the results of completed jobs.

type UpdateServer struct {
	jobCh   chan Job
	results []Result
	wg      sync.WaitGroup
}

type Job struct {
	PackageName string
	Version     string
}

type Result struct {
	PackageName string
	Version     string
	Success     bool
}

func NewUpdateServer(workers int) *UpdateServer {
	s := &UpdateServer{
		jobCh: make(chan Job, 100),
	}

	for i := 0; i < workers; i++ {
		go s.worker()
	}

	return s
}

func (s *UpdateServer) worker() {
	for job := range s.jobCh {
		// Simulate package installation
		time.Sleep(50 * time.Millisecond)
		s.results = append(s.results, Result{
			PackageName: job.PackageName,
			Version:     job.Version,
			Success:     true,
		})
	}
}

// Submit adds a job to the server. Should return false if the server
// is shutting down.
func (s *UpdateServer) Submit(job Job) bool {
	s.jobCh <- job
	return true
}

// Shutdown stops the server and returns all results.
// It should wait for in-flight jobs up to the given timeout.
func (s *UpdateServer) Shutdown(ctx context.Context) []Result {
	close(s.jobCh)
	return s.results
}
