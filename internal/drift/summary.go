package drift

import "fmt"

// DriftSummary holds aggregated results from a drift detection run.
type DriftSummary struct {
	TotalResources int
	DriftedCount   int
	ErrorCount     int
	Results        []DetectResult
}

// DetectResult holds the outcome of checking a single resource.
type DetectResult struct {
	ResourceID string
	Diffs      []Diff
	Err        error
}

// HasDrift returns true if any resources have drifted.
func (s *DriftSummary) HasDrift() bool {
	return s.DriftedCount > 0
}

// HasErrors returns true if any resources encountered fetch errors.
func (s *DriftSummary) HasErrors() bool {
	return s.ErrorCount > 0
}

// String returns a human-readable one-line summary.
func (s *DriftSummary) String() string {
	return fmt.Sprintf(
		"drift check complete: %d/%d resources drifted, %d errors",
		s.DriftedCount, s.TotalResources, s.ErrorCount,
	)
}

// BuildSummary aggregates a slice of DetectResults into a DriftSummary.
func BuildSummary(results []DetectResult) DriftSummary {
	s := DriftSummary{
		TotalResources: len(results),
		Results:        results,
	}
	for _, r := range results {
		if r.Err != nil {
			s.ErrorCount++
			continue
		}
		if len(r.Diffs) > 0 {
			s.DriftedCount++
		}
	}
	return s
}
