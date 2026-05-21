// Package schedule provides a simple interval-based scheduler for running
// drift-check jobs. The Scheduler executes a Job immediately on start and
// then repeats on a fixed interval until its context is cancelled.
//
// Typical usage:
//
//	s := schedule.New(5*time.Minute, myDriftJob, logger)
//	s.Run(ctx) // blocks until ctx is done
package schedule
