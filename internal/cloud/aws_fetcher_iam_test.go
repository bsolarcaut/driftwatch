package cloud

import (
	"context"
	"errors"
	"testing"
)

// TestAWSFetcher_IAMRole_UnsupportedType verifies that the AWSFetcher returns
// an error for resource types it does not recognise, even after adding IAM.
func TestAWSFetcher_IAMRole_UnsupportedType(t *testing.T) {
	f := &AWSFetcher{}
	_, err := f.Fetch(context.Background(), "aws_lambda_function", "my-fn")
	if err == nil {
		t.Fatal("expected error for unsupported type, got nil")
	}
}

// TestMockFetcher_SimulatesIAMRole confirms the mock fetcher can stand in for
// real IAM calls during detector tests.
func TestMockFetcher_SimulatesIAMRole(t *testing.T) {
	f := NewMockFetcher(map[string]map[string]string{
		"my-role": {
			"name": "my-role",
			"arn":  "arn:aws:iam::123456789012:role/my-role",
			"path": "/",
		},
	})

	attrs, err := f.Fetch(context.Background(), "aws_iam_role", "my-role")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attrs["arn"] != "arn:aws:iam::123456789012:role/my-role" {
		t.Errorf("unexpected arn: %q", attrs["arn"])
	}
}

// TestMockFetcher_MissingIAMRoleReturnsNotFound ensures missing IAM roles
// surface as ErrResourceNotFound via the mock.
func TestMockFetcher_MissingIAMRoleReturnsNotFound(t *testing.T) {
	f := NewMockFetcher(nil)
	_, err := f.Fetch(context.Background(), "aws_iam_role", "ghost-role")
	var nf ErrResourceNotFound
	if !errors.As(err, &nf) {
		t.Fatalf("expected ErrResourceNotFound, got %v", err)
	}
}
