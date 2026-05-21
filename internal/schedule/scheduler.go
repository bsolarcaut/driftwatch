package schedule

import (
	"context"
	"log/slog"
	"time"
)

// Job is a function that performs a drift check cycle.
type Job func(ctx context.Context) error

// Scheduler runs a Job on a fixed interval.
type Scheduler struct {
	interval time.Duration
	job      Job
	logger   *slog.Logger
}

// New creates a new Scheduler with the given interval and job.
func New(interval time.Duration, job Job, logger *slog.Logger) *Scheduler {
	return &Scheduler{
		interval: interval,
		job:      job,
		logger:   logger,
	}
}

// Run executes the job immediately, then repeats on each tick until ctx is
// cancelled. It returns only when the context is done.
func (s *Scheduler) Run(ctx context.Context) {
	s.logger.Info("scheduler started", "interval", s.interval)

	if err := s.runJob(ctx); err != nil {
		s.logger.Error("job failed", "error", err)
	}

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.runJob(ctx); err != nil {
				s.logger.Error("job failed", "error", err)
			}
		case <-ctx.Done():
			s.logger.Info("scheduler stopped")
			return
		}
	}
}

func (s *Scheduler) runJob(ctx context.Context) error {
	s.logger.Info("running drift check")
	start := time.Now()
	err := s.job(ctx)
	s.logger.Info("drift check complete", "duration", time.Since(start))
	return err
}
