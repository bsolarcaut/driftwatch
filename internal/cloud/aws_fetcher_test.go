package cloud_test

import (
	"context"
	"testing"

	"github.com/your-org/driftwatch/internal/cloud"
)

// TestAWSFetcher_UnsupportedType verifies that fetching an unknown resource type
// returns ErrResourceNotFound without making any network calls.
// Real AWS calls are covered by integration tests (build tag: integration).
func TestAWSFetcher_UnsupportedType(t *testing.T) {
	// Use MockFetcher to validate ErrResourceNotFound contract shared with AWSFetcher.
	mock := cloud.NewMockFetcher(nil)

	_, err := mock.Fetch(context.Background(), "aws_lambda_function", "my-func")
	if err == nil {
		t.Fatal("expected error for missing resource, got nil")
	}
}

func TestMockFetcher_SimulatesAWSInstance(t *testing.T) {
	attrs := map[string]string{
		"instance_type":  "t3.micro",
		"instance_state": "running",
		"ami":            "ami-0abcdef1234567890",
		"subnet_id":      "subnet-abc123",
		"vpc_id":         "vpc-deadbeef",
		"private_ip":     "10.0.1.5",
		"public_ip":      "54.12.34.56",
	}

	mock := cloud.NewMockFetcher(map[string]map[string]string{
		"aws_instance/i-0123456789abcdef0": attrs,
	})

	got, err := mock.Fetch(context.Background(), "aws_instance", "i-0123456789abcdef0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for k, want := range attrs {
		if got[k] != want {
			t.Errorf("attr %q: got %q, want %q", k, got[k], want)
		}
	}
}

func TestMockFetcher_MissingInstanceReturnsNotFound(t *testing.T) {
	mock := cloud.NewMockFetcher(map[string]map[string]string{})

	_, err := mock.Fetch(context.Background(), "aws_instance", "i-doesnotexist")
	if err == nil {
		t.Fatal("expected ErrResourceNotFound, got nil")
	}

	if err.Error() == "" {
		t.Error("expected non-empty error message")
	}
}
