package schedule_test

import (
	"context"
	"errors"
	"log/slog"
	"sync/atomic"
	"testing"
	"time"

	"github.com/yourorg/driftwatch/internal/schedule"
)

func nopLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(nil, &slog.HandlerOptions{Level: slog.LevelError}))
}

func TestScheduler_RunsJobImmediately(t *testing.T) {
	var count atomic.Int32
	job := func(_ context.Context) error {
		count.Add(1)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	s := schedule.New(10*time.Second, job, nopLogger())
	s.Run(ctx)

	if count.Load() < 1 {
		t.Fatal("expected job to run at least once immediately")
	}
}

func TestScheduler_RunsJobOnInterval(t *testing.T) {
	var count atomic.Int32
	job := func(_ context.Context) error {
		count.Add(1)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 350*time.Millisecond)
	defer cancel()

	s := schedule.New(100*time.Millisecond, job, nopLogger())
	s.Run(ctx)

	// immediate run + ~3 ticks within 350 ms
	if count.Load() < 3 {
		t.Fatalf("expected at least 3 runs, got %d", count.Load())
	}
}

func TestScheduler_ContinuesOnJobError(t *testing.T) {
	var count atomic.Int32
	job := func(_ context.Context) error {
		count.Add(1)
		return errors.New("simulated error")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	s := schedule.New(80*time.Millisecond, job, nopLogger())
	s.Run(ctx)

	if count.Load() < 2 {
		t.Fatalf("expected scheduler to continue after error, got %d runs", count.Load())
	}
}

func TestScheduler_StopsOnContextCancel(t *testing.T) {
	var count atomic.Int32
	job := func(_ context.Context) error {
		count.Add(1)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	s := schedule.New(50*time.Millisecond, job, nopLogger())

	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()

	time.Sleep(120 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// ok
	case <-time.After(500 * time.Millisecond):
		t.Fatal("scheduler did not stop after context cancellation")
	}
}
