package drift_test

import (
	"errors"
	"testing"

	"github.com/your-org/driftwatch/internal/drift"
	"github.com/your-org/driftwatch/internal/terraform"
)

type mockFetcher struct {
	attrs map[string]interface{}
	err   error
}

func (m *mockFetcher) Fetch(_ terraform.Resource) (map[string]interface{}, error) {
	return m.attrs, m.err
}

func TestDetect_NoDrift(t *testing.T) {
	attrs := map[string]interface{}{"id": "i-123", "instance_type": "t3.micro"}
	state := &terraform.State{
		Resources: []terraform.Resource{
			{Type: "aws_instance", Name: "web", Attributes: attrs},
		},
	}
	d := drift.NewDetector(&mockFetcher{attrs: attrs})
	diffs, err := d.Detect(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 0 {
		t.Errorf("expected no diffs, got %d: %v", len(diffs), diffs)
	}
}

func TestDetect_AttributeChanged(t *testing.T) {
	stateAttrs := map[string]interface{}{"id": "i-123", "instance_type": "t3.micro"}
	liveAttrs := map[string]interface{}{"id": "i-123", "instance_type": "t3.large"}
	state := &terraform.State{
		Resources: []terraform.Resource{
			{Type: "aws_instance", Name: "web", Attributes: stateAttrs},
		},
	}
	d := drift.NewDetector(&mockFetcher{attrs: liveAttrs})
	diffs, err := d.Detect(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].Attribute != "instance_type" {
		t.Errorf("expected diff on instance_type, got %s", diffs[0].Attribute)
	}
}

func TestDetect_MissingLiveAttribute(t *testing.T) {
	stateAttrs := map[string]interface{}{"id": "i-123", "tags": "env=prod"}
	liveAttrs := map[string]interface{}{"id": "i-123"}
	state := &terraform.State{
		Resources: []terraform.Resource{
			{Type: "aws_instance", Name: "web", Attributes: stateAttrs},
		},
	}
	d := drift.NewDetector(&mockFetcher{attrs: liveAttrs})
	diffs, err := d.Detect(state)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
}

func TestDetect_FetcherError(t *testing.T) {
	state := &terraform.State{
		Resources: []terraform.Resource{
			{Type: "aws_instance", Name: "web", Attributes: map[string]interface{}{"id": "i-123"}},
		},
	}
	d := drift.NewDetector(&mockFetcher{err: errors.New("API error")})
	_, err := d.Detect(state)
	if err == nil {
		t.Fatal("expected error from fetcher, got nil")
	}
}
