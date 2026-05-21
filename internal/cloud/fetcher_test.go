package cloud_test

import (
	"context"
	"errors"
	"testing"

	"github.com/your-org/driftwatch/internal/cloud"
)

func TestMockFetcher_Found(t *testing.T) {
	f := cloud.NewMockFetcher(map[string]cloud.ResourceAttributes{
		"aws_instance/i-abc123": {"instance_type": "t3.micro", "ami": "ami-0abcdef"},
	})

	attrs, err := f.Fetch(context.Background(), "aws_instance", "i-abc123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if attrs["instance_type"] != "t3.micro" {
		t.Errorf("expected instance_type=t3.micro, got %v", attrs["instance_type"])
	}
}

func TestMockFetcher_NotFound(t *testing.T) {
	f := cloud.NewMockFetcher(nil)

	_, err := f.Fetch(context.Background(), "aws_instance", "i-missing")
	if err == nil {
		t.Fatal("expected ErrResourceNotFound, got nil")
	}

	var notFound *cloud.ErrResourceNotFound
	if !errors.As(err, &notFound) {
		t.Fatalf("expected *ErrResourceNotFound, got %T", err)
	}
	if notFound.ResourceType != "aws_instance" || notFound.ResourceID != "i-missing" {
		t.Errorf("unexpected ErrResourceNotFound fields: %+v", notFound)
	}
}

func TestMockFetcher_ReturnsCopy(t *testing.T) {
	f := cloud.NewMockFetcher(map[string]cloud.ResourceAttributes{
		"aws_s3_bucket/my-bucket": {"region": "us-east-1"},
	})

	attrs, err := f.Fetch(context.Background(), "aws_s3_bucket", "my-bucket")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Mutate the returned copy.
	attrs["region"] = "eu-west-1"

	// Fetch again and confirm the original is unchanged.
	attrs2, _ := f.Fetch(context.Background(), "aws_s3_bucket", "my-bucket")
	if attrs2["region"] != "us-east-1" {
		t.Errorf("original resource was mutated; got region=%v", attrs2["region"])
	}
}

func TestMockFetcher_NilMapInitialized(t *testing.T) {
	f := cloud.NewMockFetcher(nil)
	if f.Resources == nil {
		t.Error("expected Resources map to be initialized, got nil")
	}
}

func TestErrResourceNotFound_Error(t *testing.T) {
	err := &cloud.ErrResourceNotFound{ResourceType: "aws_vpc", ResourceID: "vpc-999"}
	want := "resource not found: type=aws_vpc id=vpc-999"
	if err.Error() != want {
		t.Errorf("expected %q, got %q", want, err.Error())
	}
}
