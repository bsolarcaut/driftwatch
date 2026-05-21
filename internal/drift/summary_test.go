package drift

import (
	"errors"
	"testing"
)

func TestBuildSummary_NoDrift(t *testing.T) {
	results := []DetectResult{
		{ResourceID: "res-1", Diffs: nil, Err: nil},
		{ResourceID: "res-2", Diffs: []Diff{}, Err: nil},
	}
	s := BuildSummary(results)
	if s.TotalResources != 2 {
		t.Errorf("expected 2 total, got %d", s.TotalResources)
	}
	if s.DriftedCount != 0 {
		t.Errorf("expected 0 drifted, got %d", s.DriftedCount)
	}
	if s.ErrorCount != 0 {
		t.Errorf("expected 0 errors, got %d", s.ErrorCount)
	}
	if s.HasDrift() {
		t.Error("expected HasDrift to be false")
	}
}

func TestBuildSummary_WithDrift(t *testing.T) {
	results := []DetectResult{
		{ResourceID: "res-1", Diffs: []Diff{{Attribute: "ami", StateValue: "old", LiveValue: "new"}}, Err: nil},
		{ResourceID: "res-2", Diffs: nil, Err: nil},
	}
	s := BuildSummary(results)
	if s.DriftedCount != 1 {
		t.Errorf("expected 1 drifted, got %d", s.DriftedCount)
	}
	if !s.HasDrift() {
		t.Error("expected HasDrift to be true")
	}
}

func TestBuildSummary_WithErrors(t *testing.T) {
	results := []DetectResult{
		{ResourceID: "res-1", Err: errors.New("fetch failed")},
		{ResourceID: "res-2", Diffs: []Diff{{Attribute: "type", StateValue: "a", LiveValue: "b"}}},
	}
	s := BuildSummary(results)
	if s.ErrorCount != 1 {
		t.Errorf("expected 1 error, got %d", s.ErrorCount)
	}
	if s.DriftedCount != 1 {
		t.Errorf("expected 1 drifted, got %d", s.DriftedCount)
	}
	if !s.HasErrors() {
		t.Error("expected HasErrors to be true")
	}
}

func TestBuildSummary_Empty(t *testing.T) {
	s := BuildSummary(nil)
	if s.TotalResources != 0 {
		t.Errorf("expected 0 total, got %d", s.TotalResources)
	}
	if s.HasDrift() || s.HasErrors() {
		t.Error("expected no drift and no errors for empty results")
	}
}

func TestDriftSummary_String(t *testing.T) {
	s := DriftSummary{TotalResources: 5, DriftedCount: 2, ErrorCount: 1}
	got := s.String()
	want := "drift check complete: 2/5 resources drifted, 1 errors"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}
