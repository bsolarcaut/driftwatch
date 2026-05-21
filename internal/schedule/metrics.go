package schedule

import (
	"sync/atomic"
	"time"
)

// Metrics tracks lightweight runtime statistics for the scheduler.
type Metrics struct {
	TotalRuns   atomic.Int64
	FailedRuns  atomic.Int64
	LastRunTime atomic.Pointer[time.Time]
}

// NewMetrics returns an initialised Metrics instance.
func NewMetrics() *Metrics {
	return &Metrics{}
}

// record updates counters after each job execution.
func (m *Metrics) record(err error) {
	m.TotalRuns.Add(1)
	if err != nil {
		m.FailedRuns.Add(1)
	}
	now := time.Now()
	m.LastRunTime.Store(&now)
}

// Summary returns a snapshot of current metric values.
func (m *Metrics) Summary() map[string]any {
	summary := map[string]any{
		"total_runs":  m.TotalRuns.Load(),
		"failed_runs": m.FailedRuns.Load(),
	}
	if t := m.LastRunTime.Load(); t != nil {
		summary["last_run_time"] = t.Format(time.RFC3339)
	}
	return summary
}
